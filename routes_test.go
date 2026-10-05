package haskimail

import (
	"context"
	"reflect"
	"testing"
)

// TestAllRoutes проверяет, что каждый метод SDK обращается к нужному
// HTTP-методу и пути, передаёт query, тело и правильный токен.
func TestAllRoutes(t *testing.T) {
	ctx := context.Background()
	p := NewParams().Int("count", 1)
	type route struct {
		name, method, path string
		hasBody, hasQuery  bool
		account            bool
		call               func(*Client, *AccountClient) error
	}
	e := func(_ any, err error) error { return err }
	routes := []route{
		// Отправка
		{"DeliverMessage", "POST", "/email", true, false, false, func(c *Client, _ *AccountClient) error { return e(c.DeliverMessage(ctx, &Message{From: "a"})) }},
		{"DeliverMessages", "POST", "/email/batch", true, false, false, func(c *Client, _ *AccountClient) error { return e(c.DeliverMessages(ctx, []Message{{From: "a"}})) }},
		{"DeliverMessageWithTemplate", "POST", "/email/withTemplate", true, false, false, func(c *Client, _ *AccountClient) error {
			return e(c.DeliverMessageWithTemplate(ctx, &TemplatedMessage{TemplateID: 1}))
		}},
		{"DeliverMessagesWithTemplate", "POST", "/email/batchWithTemplates", true, false, false, func(c *Client, _ *AccountClient) error {
			return e(c.DeliverMessagesWithTemplate(ctx, []TemplatedMessage{{TemplateID: 1}}))
		}},
		// Отказы
		{"GetBounces", "GET", "/bounces", false, true, false, func(c *Client, _ *AccountClient) error { return e(c.GetBounces(ctx, p)) }},
		{"GetBounce", "GET", "/bounces/5", false, false, false, func(c *Client, _ *AccountClient) error { return e(c.GetBounce(ctx, 5)) }},
		{"ActivateBounce", "PUT", "/bounces/5/activate", false, false, false, func(c *Client, _ *AccountClient) error { return e(c.ActivateBounce(ctx, 5)) }},
		{"GetDeliveryStats", "GET", "/deliverystats", false, false, false, func(c *Client, _ *AccountClient) error { return e(c.GetDeliveryStats(ctx)) }},
		// Шаблоны
		{"GetTemplate", "GET", "/templates/7", false, false, false, func(c *Client, _ *AccountClient) error { return e(c.GetTemplate(ctx, 7)) }},
		{"GetTemplateByAlias", "GET", "/templates/welcome", false, false, false, func(c *Client, _ *AccountClient) error { return e(c.GetTemplateByAlias(ctx, "welcome")) }},
		{"CreateTemplate", "POST", "/templates", true, false, false, func(c *Client, _ *AccountClient) error { return e(c.CreateTemplate(ctx, &Template{Subject: "s"})) }},
		{"SetTemplate", "PUT", "/templates/7", true, false, false, func(c *Client, _ *AccountClient) error { return e(c.SetTemplate(ctx, 7, &Template{Subject: "s"})) }},
		{"SetTemplateByAlias", "PUT", "/templates/welcome", true, false, false, func(c *Client, _ *AccountClient) error {
			return e(c.SetTemplateByAlias(ctx, "welcome", &Template{Subject: "s"}))
		}},
		{"GetTemplates", "GET", "/templates", false, true, false, func(c *Client, _ *AccountClient) error { return e(c.GetTemplates(ctx, p)) }},
		{"DeleteTemplate", "DELETE", "/templates/7", false, false, false, func(c *Client, _ *AccountClient) error { return e(c.DeleteTemplate(ctx, 7)) }},
		{"DeleteTemplateByAlias", "DELETE", "/templates/welcome", false, false, false, func(c *Client, _ *AccountClient) error { return e(c.DeleteTemplateByAlias(ctx, "welcome")) }},
		// Исходящие
		{"GetMessages", "GET", "/messages", false, true, false, func(c *Client, _ *AccountClient) error { return e(c.GetMessages(ctx, p)) }},
		{"GetMessageDetails", "GET", "/messages/m1/details", false, false, false, func(c *Client, _ *AccountClient) error { return e(c.GetMessageDetails(ctx, "m1")) }},
		{"GetMessageOpens", "GET", "/messages/opens", false, true, false, func(c *Client, _ *AccountClient) error { return e(c.GetMessageOpens(ctx, p)) }},
		{"GetMessageOpensByID", "GET", "/messages/opens/m1", false, true, false, func(c *Client, _ *AccountClient) error { return e(c.GetMessageOpensByID(ctx, "m1", p)) }},
		{"GetMessageClicks", "GET", "/messages/clicks", false, true, false, func(c *Client, _ *AccountClient) error { return e(c.GetMessageClicks(ctx, p)) }},
		{"GetMessageClicksByID", "GET", "/messages/clicks/m1", false, true, false, func(c *Client, _ *AccountClient) error { return e(c.GetMessageClicksByID(ctx, "m1", p)) }},
		// Статистика
		{"GetOutboundStats", "GET", "/stats", false, true, false, func(c *Client, _ *AccountClient) error { return e(c.GetOutboundStats(ctx, p)) }},
		{"GetSendStats", "GET", "/stats/sends", false, true, false, func(c *Client, _ *AccountClient) error { return e(c.GetSendStats(ctx, p)) }},
		{"GetBounceStats", "GET", "/stats/bounces", false, true, false, func(c *Client, _ *AccountClient) error { return e(c.GetBounceStats(ctx, p)) }},
		{"GetSpamStats", "GET", "/stats/spam", false, true, false, func(c *Client, _ *AccountClient) error { return e(c.GetSpamStats(ctx, p)) }},
		{"GetOpenStats", "GET", "/stats/opens", false, true, false, func(c *Client, _ *AccountClient) error { return e(c.GetOpenStats(ctx, p)) }},
		{"GetOpenPlatformStats", "GET", "/stats/opens/platforms", false, true, false, func(c *Client, _ *AccountClient) error { return e(c.GetOpenPlatformStats(ctx, p)) }},
		{"GetClickStats", "GET", "/stats/clicks", false, true, false, func(c *Client, _ *AccountClient) error { return e(c.GetClickStats(ctx, p)) }},
		{"GetClickLocationStats", "GET", "/stats/clicks/location", false, true, false, func(c *Client, _ *AccountClient) error { return e(c.GetClickLocationStats(ctx, p)) }},
		{"GetClickPlatformStats", "GET", "/stats/clicks/platforms", false, true, false, func(c *Client, _ *AccountClient) error { return e(c.GetClickPlatformStats(ctx, p)) }},
		// Вебхуки
		{"GetWebhooks", "GET", "/webhooks", false, true, false, func(c *Client, _ *AccountClient) error { return e(c.GetWebhooks(ctx, p)) }},
		{"GetWebhook", "GET", "/webhooks/3", false, false, false, func(c *Client, _ *AccountClient) error { return e(c.GetWebhook(ctx, 3)) }},
		{"CreateWebhook", "POST", "/webhooks", true, false, false, func(c *Client, _ *AccountClient) error { return e(c.CreateWebhook(ctx, &Webhook{URL: "u"})) }},
		{"SetWebhook", "PUT", "/webhooks/3", true, false, false, func(c *Client, _ *AccountClient) error { return e(c.SetWebhook(ctx, 3, &Webhook{URL: "u"})) }},
		{"DeleteWebhook", "DELETE", "/webhooks/3", false, false, false, func(c *Client, _ *AccountClient) error { return e(c.DeleteWebhook(ctx, 3)) }},
		// Стоп-листы
		{"CreateSuppressions", "POST", "/message-streams/outbound/suppressions", true, false, false, func(c *Client, _ *AccountClient) error {
			return e(c.CreateSuppressions(ctx, "outbound", NewSuppressionEntries("x@y.ru")))
		}},
		{"GetSuppressions", "GET", "/message-streams/outbound/suppressions", false, true, false, func(c *Client, _ *AccountClient) error {
			return e(c.GetSuppressions(ctx, "outbound", p))
		}},
		{"DeleteSuppressions", "POST", "/message-streams/outbound/suppressions/delete", true, false, false, func(c *Client, _ *AccountClient) error {
			return e(c.DeleteSuppressions(ctx, "outbound", NewSuppressionEntries("x@y.ru")))
		}},
		// Каналы
		{"GetMessageStreams", "GET", "/message-streams", false, false, false, func(c *Client, _ *AccountClient) error { return e(c.GetMessageStreams(ctx)) }},
		{"GetMessageStream", "GET", "/message-streams/s1", false, false, false, func(c *Client, _ *AccountClient) error { return e(c.GetMessageStream(ctx, "s1")) }},
		{"CreateMessageStream", "POST", "/message-streams", true, false, false, func(c *Client, _ *AccountClient) error {
			return e(c.CreateMessageStream(ctx, &MessageStream{ID: "s1"}))
		}},
		{"SetMessageStream", "PATCH", "/message-streams/s1", true, false, false, func(c *Client, _ *AccountClient) error {
			return e(c.SetMessageStream(ctx, "s1", &MessageStream{Name: "n"}))
		}},
		{"ArchiveMessageStream", "POST", "/message-streams/s1/archive", false, false, false, func(c *Client, _ *AccountClient) error { return e(c.ArchiveMessageStream(ctx, "s1")) }},
		{"UnarchiveMessageStream", "POST", "/message-streams/s1/unarchive", false, false, false, func(c *Client, _ *AccountClient) error { return e(c.UnarchiveMessageStream(ctx, "s1")) }},

		// ===== AccountClient =====
		{"GetServer", "GET", "/servers/9", false, false, true, func(_ *Client, a *AccountClient) error { return e(a.GetServer(ctx, 9)) }},
		{"CreateServer", "POST", "/servers", true, false, true, func(_ *Client, a *AccountClient) error { return e(a.CreateServer(ctx, &Server{Name: "n"})) }},
		{"SetServer", "PUT", "/servers/9", true, false, true, func(_ *Client, a *AccountClient) error { return e(a.SetServer(ctx, 9, &Server{Name: "n"})) }},
		{"GetServers", "GET", "/servers", false, true, true, func(_ *Client, a *AccountClient) error { return e(a.GetServers(ctx, p)) }},
		{"DeleteServer", "DELETE", "/servers/9", false, false, true, func(_ *Client, a *AccountClient) error { return e(a.DeleteServer(ctx, 9)) }},
		{"GetDomains", "GET", "/domains", false, false, true, func(_ *Client, a *AccountClient) error { return e(a.GetDomains(ctx)) }},
		{"GetDomainDetails", "GET", "/domains/4", false, false, true, func(_ *Client, a *AccountClient) error { return e(a.GetDomainDetails(ctx, 4)) }},
		{"CreateDomain", "POST", "/domains", true, false, true, func(_ *Client, a *AccountClient) error { return e(a.CreateDomain(ctx, &Domain{Name: "d.ru"})) }},
		{"SetDomain", "PUT", "/domains/4", true, false, true, func(_ *Client, a *AccountClient) error { return e(a.SetDomain(ctx, 4, &Domain{Name: "d.ru"})) }},
		{"DeleteDomain", "DELETE", "/domains/4", false, false, true, func(_ *Client, a *AccountClient) error { return e(a.DeleteDomain(ctx, 4)) }},
		{"VerifyDomainSPF", "PUT", "/domains/4/verifyspf", false, false, true, func(_ *Client, a *AccountClient) error { return e(a.VerifyDomainSPF(ctx, 4)) }},
		{"VerifyDomainDKIM", "PUT", "/domains/4/verifyDkim", false, false, true, func(_ *Client, a *AccountClient) error { return e(a.VerifyDomainDKIM(ctx, 4)) }},
		{"GetSenderSignatures", "GET", "/senders", false, true, true, func(_ *Client, a *AccountClient) error { return e(a.GetSenderSignatures(ctx, p)) }},
		{"GetSenderSignatureDetails", "GET", "/senders/2", false, false, true, func(_ *Client, a *AccountClient) error { return e(a.GetSenderSignatureDetails(ctx, 2)) }},
		{"CreateSenderSignature", "POST", "/senders", true, false, true, func(_ *Client, a *AccountClient) error {
			return e(a.CreateSenderSignature(ctx, &SignatureToCreate{FromEmail: "a@b.ru"}))
		}},
		{"SetSenderSignature", "PUT", "/senders/2", true, false, true, func(_ *Client, a *AccountClient) error {
			return e(a.SetSenderSignature(ctx, 2, &SignatureToCreate{Name: "n"}))
		}},
		{"DeleteSenderSignature", "DELETE", "/senders/2", false, false, true, func(_ *Client, a *AccountClient) error { return e(a.DeleteSenderSignature(ctx, 2)) }},
		{"ResendSenderSignatureConfirmation", "POST", "/senders/2/resend", false, false, true, func(_ *Client, a *AccountClient) error {
			return e(a.ResendSenderSignatureConfirmation(ctx, 2))
		}},
		{"VerifySenderSignatureSPF", "POST", "/senders/2/verifyspf", false, false, true, func(_ *Client, a *AccountClient) error { return e(a.VerifySenderSignatureSPF(ctx, 2)) }},
		{"RequestNewSenderSignatureDKIM", "POST", "/senders/2/requestnewdkim", false, false, true, func(_ *Client, a *AccountClient) error {
			return e(a.RequestNewSenderSignatureDKIM(ctx, 2))
		}},
		{"PushTemplates", "PUT", "/templates/push", true, false, true, func(_ *Client, a *AccountClient) error {
			return e(a.PushTemplates(ctx, &TemplatesPushRequest{SourceServerID: 1, DestinationServerID: 2}))
		}},
	}

	// Полнота: каждый публичный API-метод клиентов должен быть в таблице.
	covered := map[string]bool{}
	for _, r := range routes {
		covered[r.name] = true
	}
	for _, typ := range []reflect.Type{reflect.TypeOf(&Client{}), reflect.TypeOf(&AccountClient{})} {
		for i := 0; i < typ.NumMethod(); i++ {
			name := typ.Method(i).Name
			if name == "BaseURL" {
				continue
			}
			if !covered[name] {
				t.Errorf("метод %s.%s не покрыт тестом маршрутов", typ.Elem().Name(), name)
			}
		}
	}
	if len(routes) != 66 {
		t.Errorf("ожидалось 66 маршрутов (45 серверных + 21 аккаунта), в таблице %d", len(routes))
	}

	for _, r := range routes {
		t.Run(r.name, func(t *testing.T) {
			// Ответ "{}" подходит для объектов; для пакетных методов нужен массив.
			resp := `{}`
			if r.name == "DeliverMessages" || r.name == "DeliverMessagesWithTemplate" {
				resp = `[]`
			}
			srv, c := newTestServer(t, 200, resp)
			cl := NewClient("srv-tok", WithBaseURL(host(srv)), WithInsecure())
			ac := NewAccountClient("acc-tok", WithBaseURL(host(srv)), WithInsecure())
			if err := r.call(cl, ac); err != nil {
				t.Fatalf("ошибка: %v", err)
			}
			if c.method != r.method || c.path != r.path {
				t.Errorf("запрос %s %s, ожидалось %s %s", c.method, c.path, r.method, r.path)
			}
			if (len(c.rawBody) > 0) != r.hasBody {
				t.Errorf("тело: %q, ожидалось наличие=%v", c.rawBody, r.hasBody)
			}
			if (c.query != "") != r.hasQuery {
				t.Errorf("query: %q, ожидалось наличие=%v", c.query, r.hasQuery)
			}
			wantHdr, otherHdr, wantTok := ServerTokenHeader, AccountTokenHeader, "srv-tok"
			if r.account {
				wantHdr, otherHdr, wantTok = AccountTokenHeader, ServerTokenHeader, "acc-tok"
			}
			if c.header.Get(wantHdr) != wantTok || c.header.Get(otherHdr) != "" {
				t.Errorf("токен: %s=%q, %s=%q", wantHdr, c.header.Get(wantHdr), otherHdr, c.header.Get(otherHdr))
			}
		})
	}
}
