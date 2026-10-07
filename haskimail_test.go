package haskimail

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"
)

type captured struct {
	method, path, query string
	header              http.Header
	body                map[string]any
	rawBody             []byte
}

func newTestServer(t *testing.T, status int, resp string) (*httptest.Server, *captured) {
	t.Helper()
	c := &captured{}
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		c.method, c.path, c.query, c.header = r.Method, r.URL.EscapedPath(), r.URL.RawQuery, r.Header.Clone()
		c.rawBody, _ = io.ReadAll(r.Body)
		if len(c.rawBody) > 0 {
			_ = json.Unmarshal(c.rawBody, &c.body)
		}
		w.WriteHeader(status)
		_, _ = io.WriteString(w, resp)
	}))
	t.Cleanup(srv.Close)
	return srv, c
}

func host(srv *httptest.Server) string { return strings.TrimPrefix(srv.URL, "http://") }

func TestDeliverMessage(t *testing.T) {
	srv, c := newTestServer(t, 200, `{"ErrorCode":0,"Message":"OK","To":"to@test.ru","SubmittedAt":"2024-03-01T10:20:30.1234567+03:00","MessageID":"abc-123"}`)
	cl := NewClient("tok", WithBaseURL(host(srv)), WithInsecure())

	m := &Message{From: "from@test.ru", To: "to@test.ru", Subject: "Тема", HTMLBody: "<b>Тело</b>",
		TrackOpens: Bool(false), TrackLinks: TrackLinksHTMLAndText}
	m.AddMetadata("campaign", "welcome")
	m.AddHeader("X-Test", "1")

	r, err := cl.DeliverMessage(context.Background(), m)
	if err != nil {
		t.Fatal(err)
	}
	if r.MessageID != "abc-123" || r.SubmittedAt == nil || r.SubmittedAt.Year() != 2024 {
		t.Fatalf("плохой ответ: %+v", r)
	}
	if c.method != "POST" || c.path != "/email" {
		t.Fatalf("запрос: %s %s", c.method, c.path)
	}
	if c.header.Get(ServerTokenHeader) != "tok" || c.header.Get("Content-Type") != "application/json" ||
		!strings.HasPrefix(c.header.Get("User-Agent"), "Haskimail-Go") {
		t.Fatalf("заголовки: %v", c.header)
	}
	for _, k := range []string{"From", "To", "Subject", "HtmlBody", "TrackLinks", "Metadata", "Headers"} {
		if _, ok := c.body[k]; !ok {
			t.Errorf("нет поля %s в %s", k, c.rawBody)
		}
	}
	if v, ok := c.body["TrackOpens"]; !ok || v != false {
		t.Errorf("TrackOpens=false должен передаваться явно: %s", c.rawBody)
	}
	for _, k := range []string{"Cc", "Bcc", "TextBody", "Attachments"} {
		if _, ok := c.body[k]; ok {
			t.Errorf("пустое поле %s не должно сериализоваться", k)
		}
	}
}

func TestErrors(t *testing.T) {
	cases := []struct {
		status int
		want   error
	}{
		{401, ErrInvalidAPIKey}, {408, ErrTimeout}, {422, ErrInvalidMessage},
		{500, ErrInternalServer}, {503, ErrUnknown},
	}
	for _, tc := range cases {
		srv, _ := newTestServer(t, tc.status, `{"ErrorCode":10,"Message":"Bad"}`)
		cl := NewClient("tok", WithBaseURL(host(srv)), WithInsecure())
		_, err := cl.GetDeliveryStats(context.Background())
		if !errors.Is(err, tc.want) {
			t.Fatalf("%d: ожидалась %v, получено %v", tc.status, tc.want, err)
		}
		var apiErr *Error
		if !errors.As(err, &apiErr) || apiErr.ErrorCode != 10 || apiErr.StatusCode != tc.status || apiErr.Message != "Bad" {
			t.Fatalf("%d: %+v", tc.status, apiErr)
		}
	}
	srv, _ := newTestServer(t, 422, `не json`)
	_, err := NewClient("t", WithBaseURL(host(srv)), WithInsecure()).GetDeliveryStats(context.Background())
	var apiErr *Error
	if !errors.As(err, &apiErr) || apiErr.ErrorCode != 0 || apiErr.Message != "Неизвестная ошибка" {
		t.Fatalf("некорректный JSON: %v", err)
	}
}

func TestParamsAndPaths(t *testing.T) {
	srv, c := newTestServer(t, 200, `{"TotalCount":1,"Bounces":[{"ID":42,"Type":"HardBounce","BouncedAt":"2024-01-02 03:04:05"}]}`)
	cl := NewClient("tok", WithBaseURL(host(srv)), WithInsecure())
	p := NewParams().Int("count", 10).Int("offset", 0).String("tag", "a b&c").Bool("inactive", true).
		Date("fromdate", time.Date(2024, 5, 6, 0, 0, 0, 0, time.UTC))
	b, err := cl.GetBounces(context.Background(), p)
	if err != nil {
		t.Fatal(err)
	}
	if b.TotalCount != 1 || b.Bounces[0].ID != 42 || b.Bounces[0].BouncedAt.Hour() != 3 {
		t.Fatalf("%+v", b)
	}
	if c.query != "count=10&fromdate=2024-05-06&inactive=true&offset=0&tag=a+b%26c" {
		t.Fatalf("query: %s", c.query)
	}

	_, _ = cl.GetMessageDetails(context.Background(), "id/with space")
	if c.path != "/messages/id%2Fwith%20space/details" {
		t.Fatalf("path: %s", c.path)
	}
	_, _ = cl.GetSuppressions(context.Background(), "outbound", nil)
	if c.path != "/message-streams/outbound/suppressions" || c.query != "" {
		t.Fatalf("path: %s?%s", c.path, c.query)
	}
}

func TestAccountClientHeaderAndArrayQuirk(t *testing.T) {
	// API иногда возвращает массив вместо объекта — берём первый элемент.
	srv, c := newTestServer(t, 200, `[{"ID":7,"Name":"d.ru","DKIMVerified":true}]`)
	ac := NewAccountClient("acc", WithBaseURL(host(srv)), WithInsecure())
	r, err := ac.GetDomainDetails(context.Background(), 7)
	if err != nil {
		t.Fatal(err)
	}
	if r.ID != 7 || r.Name != "d.ru" || !r.DKIMVerified {
		t.Fatalf("%+v", r)
	}
	if c.header.Get(AccountTokenHeader) != "acc" || c.header.Get(ServerTokenHeader) != "" {
		t.Fatalf("заголовки: %v", c.header)
	}
}

func TestTimeFormats(t *testing.T) {
	for _, s := range []string{
		"2024-01-02 03:04:05",
		"2024-01-02T03:04:05",
		"2024-01-02T03:04:05Z",
		"2024-01-02T03:04:05+0300",
		"2024-01-02T03:04:05.1234567+03:00",
		"2024-01-02T03:04:05.123+0300",
		"Tue, 02 Jan 2024 03:04:05 +0300",
		"Tue, 2 Jan 2024 03:04:05 +0300",
	} {
		var v struct{ T *Time }
		if err := json.Unmarshal([]byte(`{"T":"`+s+`"}`), &v); err != nil || v.T.Day() != 2 || v.T.Minute() != 4 {
			t.Errorf("%q: %v %v", s, err, v.T)
		}
	}
	var d struct{ T *Time }
	if err := json.Unmarshal([]byte(`{"T":"2024-01-02"}`), &d); err != nil || d.T.Day() != 2 {
		t.Errorf("дата без времени: %v", err)
	}
	var v struct{ T *Time }
	if err := json.Unmarshal([]byte(`{"T":null}`), &v); err != nil || v.T != nil {
		t.Errorf("null: %v", err)
	}
	if err := json.Unmarshal([]byte(`{"T":"мусор"}`), &v); err == nil {
		t.Error("ожидалась ошибка")
	}
}

func TestAttachmentsAndRecipients(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "report.pdf")
	if err := os.WriteFile(path, []byte("%PDF-1.4 test"), 0o600); err != nil {
		t.Fatal(err)
	}
	a, err := AttachmentFromFile(path, "")
	if err != nil {
		t.Fatal(err)
	}
	dec, _ := base64.StdEncoding.DecodeString(a.Content)
	if a.Name != "report.pdf" || a.ContentType != "application/pdf" || string(dec) != "%PDF-1.4 test" {
		t.Fatalf("%+v", a)
	}
	raw, _ := json.Marshal(NewAttachment("a.png", []byte{1}, "image/png", "cid:logo"))
	if !strings.Contains(string(raw), `"ContentId":"cid:logo"`) {
		t.Fatalf("%s", raw)
	}
	got := FormatRecipientsMap(map[string]string{"Пётр": "p@x.ru", "Анна": "a@x.ru"})
	if got != `"Анна"<a@x.ru>,"Пётр"<p@x.ru>` {
		t.Fatalf("%s", got)
	}
}

func TestContextCancel(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		<-r.Context().Done()
	}))
	defer srv.Close()
	ctx, cancel := context.WithTimeout(context.Background(), 50*time.Millisecond)
	defer cancel()
	_, err := NewClient("t", WithBaseURL(host(srv)), WithInsecure()).GetDeliveryStats(ctx)
	if !errors.Is(err, context.DeadlineExceeded) {
		t.Fatalf("ожидался DeadlineExceeded, получено %v", err)
	}
}
