//go:build integration

// Интеграционные тесты против живого API Haskimail.
//
// Запуск:
//
//	export HASKIMAIL_SERVER_TOKEN=...      # обязательно
//	export HASKIMAIL_ACCOUNT_TOKEN=...     # для тестов API аккаунта
//	export HASKIMAIL_SENDER_EMAIL=...      # подтверждённый отправитель — для отправки писем
//	export HASKIMAIL_RECIPIENT_EMAIL=...   # получатель тестовых писем
//	export HASKIMAIL_API_URL=...           # необязательно, хост API
//	export HASKIMAIL_INSECURE=1            # необязательно, HTTP вместо HTTPS
//	go test -tags integration -v -run Integration ./...
//
// Тесты только читают данные, кроме: отправки писем (если заданы SENDER и
// RECIPIENT) и создания временных шаблона и вебхука, которые сразу удаляются.
//
// Дополнительно каждый ответ проверяется на поля, которых нет в Go-моделях:
// такие поля выводятся как предупреждения "НЕИЗВЕСТНЫЕ ПОЛЯ".
package haskimail

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
	"sync"
	"testing"
	"time"
)

// recorder запоминает тело последнего ответа для проверки полноты моделей.
type recorder struct {
	mu   sync.Mutex
	last []byte
}

func (r *recorder) RoundTrip(req *http.Request) (*http.Response, error) {
	resp, err := http.DefaultTransport.RoundTrip(req)
	if err != nil {
		return nil, err
	}
	body, _ := io.ReadAll(resp.Body)
	resp.Body.Close()
	resp.Body = io.NopCloser(bytes.NewReader(body))
	r.mu.Lock()
	r.last = body
	r.mu.Unlock()
	return resp, nil
}

func (r *recorder) body() []byte {
	r.mu.Lock()
	defer r.mu.Unlock()
	return r.last
}

type env struct {
	server, account, sender, recipient string
	rec                                *recorder
	opts                               []Option
}

func loadEnv(t *testing.T) *env {
	t.Helper()
	e := &env{
		server:    os.Getenv("HASKIMAIL_SERVER_TOKEN"),
		account:   os.Getenv("HASKIMAIL_ACCOUNT_TOKEN"),
		sender:    os.Getenv("HASKIMAIL_SENDER_EMAIL"),
		recipient: os.Getenv("HASKIMAIL_RECIPIENT_EMAIL"),
		rec:       &recorder{},
	}
	if e.server == "" {
		t.Skip("HASKIMAIL_SERVER_TOKEN не задан")
	}
	e.opts = []Option{WithHTTPClient(&http.Client{Transport: e.rec, Timeout: 30 * time.Second})}
	if u := os.Getenv("HASKIMAIL_API_URL"); u != "" {
		e.opts = append(e.opts, WithBaseURL(u))
	}
	if os.Getenv("HASKIMAIL_INSECURE") == "1" {
		e.opts = append(e.opts, WithInsecure())
	}
	return e
}

func ctxT(t *testing.T) context.Context {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	t.Cleanup(cancel)
	return ctx
}

// checkModel повторно разбирает последний ответ в строгом режиме и
// сообщает о полях API, отсутствующих в модели.
func checkModel[T any](t *testing.T, rec *recorder) {
	t.Helper()
	raw := bytes.TrimSpace(rec.body())
	if len(raw) == 0 {
		return
	}
	var target any = new(T)
	if raw[0] == '[' {
		var probe T
		if !strings.HasPrefix(fmt.Sprintf("%T", probe), "[]") {
			target = new([]T)
		}
	}
	dec := json.NewDecoder(bytes.NewReader(raw))
	dec.DisallowUnknownFields()
	if err := dec.Decode(target); err != nil {
		t.Logf("НЕИЗВЕСТНЫЕ ПОЛЯ в %T: %v\nответ: %s", *new(T), err, truncate(raw, 600))
	}
}

func truncate(b []byte, n int) string {
	if len(b) > n {
		return string(b[:n]) + "…"
	}
	return string(b)
}

func must[T any](t *testing.T, e *env, v *T, err error) *T {
	t.Helper()
	if err != nil {
		t.Fatalf("ошибка API: %v", err)
	}
	if v == nil {
		t.Fatal("пустой результат")
	}
	checkModel[T](t, e.rec)
	return v
}

func page() *Params { return NewParams().Int("count", 10).Int("offset", 0) }

// ===== Серверный API: только чтение =====

func TestIntegrationServerRead(t *testing.T) {
	e := loadEnv(t)
	c := NewClient(e.server, e.opts...)

	t.Run("DeliveryStats", func(t *testing.T) {
		s, err := c.GetDeliveryStats(ctxT(t))
		s = must(t, e, s, err)
		t.Logf("неактивных адресов: %d, типов отказов: %d", s.InactiveMails, len(s.Bounces))
	})
	t.Run("Bounces", func(t *testing.T) {
		b, err := c.GetBounces(ctxT(t), page())
		b = must(t, e, b, err)
		t.Logf("отказов всего: %d", b.TotalCount)
		if len(b.Bounces) > 0 && b.Bounces[0].BouncedAt == nil {
			t.Error("BouncedAt не разобран")
		}
	})
	t.Run("MessageStreams", func(t *testing.T) {
		s, err := c.GetMessageStreams(ctxT(t))
		s = must(t, e, s, err)
		for _, ms := range s.MessageStreams {
			t.Logf("канал %q (%s)", ms.ID, ms.MessageStreamType)
		}
		if len(s.MessageStreams) > 0 {
			one, err := c.GetMessageStream(ctxT(t), s.MessageStreams[0].ID)
			must(t, e, one, err)
		}
	})
	t.Run("Messages", func(t *testing.T) {
		m, err := c.GetMessages(ctxT(t), page())
		m = must(t, e, m, err)
		t.Logf("исходящих сообщений: %d", m.TotalCount)
		if len(m.Messages) > 0 {
			id := m.Messages[0].MessageID
			d, err := c.GetMessageDetails(ctxT(t), id)
			d = must(t, e, d, err)
			if d.MessageID != id {
				t.Errorf("MessageID: %q != %q", d.MessageID, id)
			}
		}
	})
	t.Run("Opens", func(t *testing.T) {
		o, err := c.GetMessageOpens(ctxT(t), page())
		must(t, e, o, err)
	})
	t.Run("Clicks", func(t *testing.T) {
		o, err := c.GetMessageClicks(ctxT(t), page())
		must(t, e, o, err)
	})
	t.Run("Templates", func(t *testing.T) {
		o, err := c.GetTemplates(ctxT(t), page())
		must(t, e, o, err)
	})
	t.Run("Webhooks", func(t *testing.T) {
		o, err := c.GetWebhooks(ctxT(t), nil)
		must(t, e, o, err)
	})
	t.Run("Suppressions", func(t *testing.T) {
		o, err := c.GetSuppressions(ctxT(t), "outbound", nil)
		must(t, e, o, err)
	})

	// Статистика
	stats := map[string]func(*testing.T, context.Context){
		"Outbound": func(t *testing.T, ctx context.Context) { v, err := c.GetOutboundStats(ctx, nil); must(t, e, v, err) },
		"Sends":    func(t *testing.T, ctx context.Context) { v, err := c.GetSendStats(ctx, nil); must(t, e, v, err) },
		"Bounces":  func(t *testing.T, ctx context.Context) { v, err := c.GetBounceStats(ctx, nil); must(t, e, v, err) },
		"Spam":     func(t *testing.T, ctx context.Context) { v, err := c.GetSpamStats(ctx, nil); must(t, e, v, err) },
		"Opens":    func(t *testing.T, ctx context.Context) { v, err := c.GetOpenStats(ctx, nil); must(t, e, v, err) },
		"Clicks":   func(t *testing.T, ctx context.Context) { v, err := c.GetClickStats(ctx, nil); must(t, e, v, err) },
	}
	for name, call := range stats {
		t.Run("Stats/"+name, func(t *testing.T) { call(t, ctxT(t)) })
	}
}

// ===== Ошибки =====

func TestIntegrationInvalidToken(t *testing.T) {
	e := loadEnv(t)
	c := NewClient("заведомо-неверный-токен", e.opts...)
	_, err := c.GetDeliveryStats(ctxT(t))
	if err == nil {
		t.Fatal("ожидалась ошибка для неверного токена")
	}
	var apiErr *Error
	if !asError(err, &apiErr) {
		t.Fatalf("ожидалась *Error, получено %T: %v", err, err)
	}
	t.Logf("HTTP %d, код %d: %s (ErrInvalidAPIKey=%v)", apiErr.StatusCode, apiErr.ErrorCode, apiErr.Message, apiErr.Unwrap() == ErrInvalidAPIKey)
	if apiErr.StatusCode != http.StatusUnauthorized {
		t.Logf("ВНИМАНИЕ: API вернул %d вместо 401 — проверьте маппинг ошибок", apiErr.StatusCode)
	}
}

func TestIntegrationInvalidMessage(t *testing.T) {
	e := loadEnv(t)
	c := NewClient(e.server, e.opts...)
	_, err := c.DeliverMessage(ctxT(t), &Message{From: "не-адрес", To: "тоже-не-адрес", Subject: "x", TextBody: "x"})
	var apiErr *Error
	if !asError(err, &apiErr) {
		t.Fatalf("ожидалась *Error, получено %v", err)
	}
	t.Logf("HTTP %d, код %d: %s", apiErr.StatusCode, apiErr.ErrorCode, apiErr.Message)
}

// ===== Изменяющие операции с откатом =====

func TestIntegrationTemplateLifecycle(t *testing.T) {
	e := loadEnv(t)
	c := NewClient(e.server, e.opts...)
	alias := fmt.Sprintf("go-sdk-test-%d", time.Now().UnixNano())

	created, err := c.CreateTemplate(ctxT(t), &Template{
		BaseTemplate: BaseTemplate{Name: "Go SDK test", Alias: alias},
		Subject:      "Привет, {{name}}",
		HTMLBody:     "<p>Привет, {{name}}!</p>",
		TextBody:     "Привет, {{name}}!",
	})
	created = must(t, e, created, err)
	if created.TemplateID == 0 {
		t.Fatalf("TemplateID не получен: %+v", created)
	}
	t.Cleanup(func() {
		if _, err := c.DeleteTemplate(context.Background(), created.TemplateID); err != nil {
			t.Errorf("не удалось удалить шаблон %d: %v", created.TemplateID, err)
		}
	})

	byID, err := c.GetTemplate(ctxT(t), created.TemplateID)
	byID = must(t, e, byID, err)
	byAlias, err := c.GetTemplateByAlias(ctxT(t), alias)
	byAlias = must(t, e, byAlias, err)
	if byID.TemplateID != byAlias.TemplateID {
		t.Errorf("по ID и по алиасу разные шаблоны: %d / %d", byID.TemplateID, byAlias.TemplateID)
	}

	upd, err := c.SetTemplate(ctxT(t), created.TemplateID, &Template{Subject: "Обновлено"})
	must(t, e, upd, err)

	if e.sender != "" && e.recipient != "" {
		r, err := c.DeliverMessageWithTemplate(ctxT(t), &TemplatedMessage{
			TemplateID: created.TemplateID, From: e.sender, To: e.recipient,
			TemplateModel: map[string]any{"name": "Go SDK"},
		})
		r = must(t, e, r, err)
		t.Logf("письмо по шаблону: %s", r.MessageID)
	}
}

func TestIntegrationWebhookLifecycle(t *testing.T) {
	e := loadEnv(t)
	c := NewClient(e.server, e.opts...)

	created, err := c.CreateWebhook(ctxT(t), &Webhook{
		URL:           "https://example.com/haskimail-go-sdk-test",
		MessageStream: "outbound",
		HTTPHeaders:   []Header{{Name: "X-Test", Value: "1"}},
		Triggers: &WebhookTriggers{
			Delivery: &WebhookTrigger{Enabled: Bool(true)},
			Open:     &OpenWebhookTrigger{Enabled: Bool(true), PostFirstOpenOnly: Bool(false)},
			Bounce:   &BounceWebhookTrigger{Enabled: Bool(false), IncludeContent: Bool(false)},
		},
	})
	created = must(t, e, created, err)
	if created.ID == 0 {
		t.Fatalf("ID вебхука не получен: %+v", created)
	}
	t.Cleanup(func() {
		if _, err := c.DeleteWebhook(context.Background(), created.ID); err != nil {
			t.Errorf("не удалось удалить вебхук %d: %v", created.ID, err)
		}
	})

	got, err := c.GetWebhook(ctxT(t), created.ID)
	got = must(t, e, got, err)
	if got.Triggers == nil || got.Triggers.Delivery == nil || got.Triggers.Delivery.Enabled == nil || !*got.Triggers.Delivery.Enabled {
		t.Errorf("триггер Delivery не сохранился: %+v", got.Triggers)
	}
	if got.Triggers != nil && got.Triggers.Bounce != nil && got.Triggers.Bounce.Enabled != nil && *got.Triggers.Bounce.Enabled {
		t.Error("Bounce.Enabled=false не дошёл до сервера")
	}
}

// ===== Отправка писем =====

func TestIntegrationSend(t *testing.T) {
	e := loadEnv(t)
	if e.sender == "" || e.recipient == "" {
		t.Skip("HASKIMAIL_SENDER_EMAIL / HASKIMAIL_RECIPIENT_EMAIL не заданы")
	}
	c := NewClient(e.server, e.opts...)

	t.Run("Simple", func(t *testing.T) {
		r, err := c.DeliverMessage(ctxT(t), &Message{
			From: e.sender, To: e.recipient, Subject: "Go SDK: простое письмо",
			HTMLBody: "<p>Тест <b>Go SDK</b></p>", TextBody: "Тест Go SDK",
		})
		r = must(t, e, r, err)
		if r.MessageID == "" || r.ErrorCode != 0 {
			t.Errorf("ответ: %+v", r)
		}
		if r.SubmittedAt == nil {
			t.Error("SubmittedAt не разобран")
		}
		t.Logf("MessageID=%s SubmittedAt=%v", r.MessageID, r.SubmittedAt)
	})

	t.Run("TrackingHeadersMetadataAttachment", func(t *testing.T) {
		m := &Message{
			From: e.sender, To: e.recipient, Subject: "Go SDK: отслеживание и вложение",
			HTMLBody:   `<p>Ссылка: <a href="https://example.com">example</a></p><img src="cid:dot">`,
			TrackOpens: Bool(true), TrackLinks: TrackLinksHTMLAndText, Tag: "go-sdk-test",
		}
		m.AddHeader("X-Go-SDK", "test")
		m.AddMetadata("source", "go-sdk-integration")
		m.AddAttachment(NewAttachment("hello.txt", []byte("Привет из Go"), "text/plain", ""))
		r, err := c.DeliverMessage(ctxT(t), m)
		r = must(t, e, r, err)
		t.Logf("MessageID=%s", r.MessageID)
	})
}

// ===== API аккаунта: только чтение =====

func TestIntegrationAccountRead(t *testing.T) {
	e := loadEnv(t)
	if e.account == "" {
		t.Skip("HASKIMAIL_ACCOUNT_TOKEN не задан")
	}
	a := NewAccountClient(e.account, e.opts...)

	t.Run("Servers", func(t *testing.T) {
		s, err := a.GetServers(ctxT(t), page())
		s = must(t, e, s, err)
		t.Logf("серверов: %d", s.TotalCount)
		if len(s.Servers) > 0 {
			one, err := a.GetServer(ctxT(t), s.Servers[0].ID)
			must(t, e, one, err)
		}
	})
	t.Run("Domains", func(t *testing.T) {
		d, err := a.GetDomains(ctxT(t))
		d = must(t, e, d, err)
		if len(d.Domains) > 0 {
			one, err := a.GetDomainDetails(ctxT(t), d.Domains[0].ID)
			must(t, e, one, err)
		}
	})
}

func asError(err error, target **Error) bool {
	e, ok := err.(*Error)
	if ok {
		*target = e
	}
	return ok
}
