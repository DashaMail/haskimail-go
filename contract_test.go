package haskimail

import (
	"encoding/json"
	"strings"
	"testing"
)

// Ответы и запросы в том виде, в каком их отдаёт и читает API Haskimail
// (haski/api/functions на сервере): сверка типов и имён JSON-полей.

func TestContractDomainDetails(t *testing.T) {
	body := `{"ID":3127,"Name":"example.com","DKIMVerified":true,"WeakDKIM":false,"SPFVerified":true,
		"ReturnPathDomainVerified":true,
		"SPFTextValue":{"valid":1,"name":"example.com.","record_type":"TXT","value":"v=spf1 include:_spf.haskisender.ru ~all"},
		"DKIMTextValue":{"valid":0,"name":"dm._domainkey.example.com.","record_type":"TXT","value":"v=DKIM1;p=AAA;t=s"},
		"DKIMHost":"dm._domainkey.example.com"}`
	var d DomainDetails
	if err := json.Unmarshal([]byte(body), &d); err != nil {
		t.Fatal(err)
	}
	if d.ID != 3127 || !d.SPFVerified || !d.DKIMVerified || d.DKIMHost != "dm._domainkey.example.com" {
		t.Fatalf("домен разобран неверно: %+v", d)
	}
	if d.SPFTextValue == nil || d.SPFTextValue.Valid != 1 || !strings.Contains(d.SPFTextValue.Value, "_spf.haskisender.ru") {
		t.Fatalf("SPFTextValue: %+v", d.SPFTextValue)
	}
	if d.DKIMTextValue == nil || d.DKIMTextValue.Name != "dm._domainkey.example.com." || d.DKIMTextValue.RecordType != "TXT" {
		t.Fatalf("DKIMTextValue: %+v", d.DKIMTextValue)
	}
}

func TestContractBounces(t *testing.T) {
	body := `{"TotalCount":1,"Bounces":[{"Type":"HardBounce","Name":"Hard bounce","Email":"a@example.com",
		"RecordType":"Bounce","MessageID":"6f1c","ServerID":11834,"MessageStream":"2081","Description":"550",
		"BouncedAt":"2026-10-07 14:30:15","TypeCode":"550","Subject":"s","Tag":"orders"}]}`
	var b Bounces
	if err := json.Unmarshal([]byte(body), &b); err != nil {
		t.Fatal(err)
	}
	if got := b.Bounces[0]; got.TypeCode != "550" || got.ServerID != 11834 || got.MessageStream != "2081" || got.BouncedAt == nil {
		t.Fatalf("возврат разобран неверно: %+v", got)
	}
}

func TestContractStreams(t *testing.T) {
	body := `{"TotalCount":1,"MessageStreams":[{"ID":"2081","ServerID":11834,"Name":"Заказы",
		"MessageStreamType":"transactional","ArchivedAt":null}]}`
	var s MessageStreams
	if err := json.Unmarshal([]byte(body), &s); err != nil {
		t.Fatal(err)
	}
	if got := s.MessageStreams[0]; got.ID != "2081" || got.ServerID != 11834 || got.ArchivedAt != nil {
		t.Fatalf("канал разобран неверно: %+v", got)
	}

	var a MessageStreamArchiveResponse
	if err := json.Unmarshal([]byte(`{"ID":"2081","ServerID":11834,"ArchivedAt":"2026-10-07 15:00:00",
		"ExpectedPurgeDate":"2026-11-21 15:00:00"}`), &a); err != nil {
		t.Fatal(err)
	}
	if a.ArchivedAt == nil || a.ExpectedPurgeDate == nil || a.ServerID != 11834 {
		t.Fatalf("архивация разобрана неверно: %+v", a)
	}
}

func TestContractRequestNames(t *testing.T) {
	// Сервер читает эти поля с учётом регистра.
	cases := map[string]any{
		`"SourceServerID":1`: TemplatesPushRequest{SourceServerID: 1, DestinationServerID: 2},
		`"ServerID":11834`:   MessageStream{ServerID: 11834, Name: "Заказы", MessageStreamType: MessageStreamTransactional},
	}
	for want, v := range cases {
		b, err := json.Marshal(v)
		if err != nil {
			t.Fatal(err)
		}
		if !strings.Contains(string(b), want) {
			t.Errorf("в %s нет %s", b, want)
		}
	}
}

func TestContractStatsAndZeroDate(t *testing.T) {
	var st OutboundStats
	if err := json.Unmarshal([]byte(`{"Sent":615,"Bounced":64,"Opened":166,"Clicked":72,"Unsubsribed":3,
		"SpamComplaints":10,"BounceRate":10.41,"SpamComplaintsRate":1.63,"OpenRate":26.99,"ClickRate":11.71}`), &st); err != nil {
		t.Fatal(err)
	}
	if st.Opened != 166 || st.Clicked != 72 || st.Unsubscribed != 3 || st.OpenRate != 26.99 {
		t.Fatalf("сводка разобрана неверно: %+v", st)
	}

	var s Suppressions
	if err := json.Unmarshal([]byte(`{"Suppressions":[{"EmailAddress":"a@example.com","SuppressionReason":"ManualSuppression",
		"Origin":"Recipient","CreatedAt":"0000-00-00 00:00:00"}]}`), &s); err != nil {
		t.Fatal(err)
	}
	if s.Suppressions[0].CreatedAt != nil && !s.Suppressions[0].CreatedAt.IsZero() {
		t.Fatalf("пустая дата MySQL должна стать нулевой: %v", s.Suppressions[0].CreatedAt)
	}
}
