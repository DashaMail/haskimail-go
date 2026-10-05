package haskimail

import (
	"context"
	"net/http"
	"net/url"
)

// Client — клиент серверного API Haskimail. Безопасен для конкурентного использования.
type Client struct {
	*base
}

// ========== Отправка сообщений ==========

// DeliverMessage отправляет одно письмо.
func (c *Client) DeliverMessage(ctx context.Context, m *Message) (*MessageResponse, error) {
	return do[MessageResponse](ctx, c.base, http.MethodPost, "/email", nil, m)
}

// DeliverMessages отправляет пакет писем. Результат по каждому письму —
// в соответствующем элементе ответа (проверяйте ErrorCode).
func (c *Client) DeliverMessages(ctx context.Context, ms []Message) ([]MessageResponse, error) {
	r, err := do[[]MessageResponse](ctx, c.base, http.MethodPost, "/email/batch", nil, ms)
	if err != nil {
		return nil, err
	}
	return *r, nil
}

// DeliverMessageWithTemplate отправляет письмо по шаблону.
func (c *Client) DeliverMessageWithTemplate(ctx context.Context, m *TemplatedMessage) (*MessageResponse, error) {
	return do[MessageResponse](ctx, c.base, http.MethodPost, "/email/withTemplate", nil, m)
}

// DeliverMessagesWithTemplate отправляет пакет писем по шаблонам.
func (c *Client) DeliverMessagesWithTemplate(ctx context.Context, ms []TemplatedMessage) ([]MessageResponse, error) {
	body := struct {
		Messages []TemplatedMessage `json:"Messages"`
	}{ms}
	r, err := do[[]MessageResponse](ctx, c.base, http.MethodPost, "/email/batchWithTemplates", nil, body)
	if err != nil {
		return nil, err
	}
	return *r, nil
}

// ========== Отказы доставки ==========

// GetBounces возвращает список отказов. Обычно требуются параметры count и offset.
func (c *Client) GetBounces(ctx context.Context, p *Params) (*Bounces, error) {
	return do[Bounces](ctx, c.base, http.MethodGet, "/bounces", p, nil)
}

// GetBounce возвращает отказ по ID.
func (c *Client) GetBounce(ctx context.Context, id int64) (*Bounce, error) {
	return do[Bounce](ctx, c.base, http.MethodGet, "/bounces/"+itoa(id), nil, nil)
}

// ActivateBounce повторно активирует адрес после отказа.
func (c *Client) ActivateBounce(ctx context.Context, id int64) (*Bounce, error) {
	return do[Bounce](ctx, c.base, http.MethodPut, "/bounces/"+itoa(id)+"/activate", nil, nil)
}

// GetDeliveryStats возвращает сводку по доставке.
func (c *Client) GetDeliveryStats(ctx context.Context) (*DeliveryStats, error) {
	return do[DeliveryStats](ctx, c.base, http.MethodGet, "/deliverystats", nil, nil)
}

// ========== Шаблоны ==========

// GetTemplate возвращает шаблон по ID.
func (c *Client) GetTemplate(ctx context.Context, id int64) (*Template, error) {
	return do[Template](ctx, c.base, http.MethodGet, "/templates/"+itoa(id), nil, nil)
}

// GetTemplateByAlias возвращает шаблон по алиасу.
func (c *Client) GetTemplateByAlias(ctx context.Context, alias string) (*Template, error) {
	return do[Template](ctx, c.base, http.MethodGet, "/templates/"+url.PathEscape(alias), nil, nil)
}

// CreateTemplate создаёт шаблон.
func (c *Client) CreateTemplate(ctx context.Context, t *Template) (*Template, error) {
	return do[Template](ctx, c.base, http.MethodPost, "/templates", nil, t)
}

// SetTemplate изменяет шаблон по ID.
func (c *Client) SetTemplate(ctx context.Context, id int64, t *Template) (*Template, error) {
	return do[Template](ctx, c.base, http.MethodPut, "/templates/"+itoa(id), nil, t)
}

// SetTemplateByAlias изменяет шаблон по алиасу.
func (c *Client) SetTemplateByAlias(ctx context.Context, alias string, t *Template) (*Template, error) {
	return do[Template](ctx, c.base, http.MethodPut, "/templates/"+url.PathEscape(alias), nil, t)
}

// GetTemplates возвращает список шаблонов.
func (c *Client) GetTemplates(ctx context.Context, p *Params) (*Templates, error) {
	return do[Templates](ctx, c.base, http.MethodGet, "/templates", p, nil)
}

// DeleteTemplate удаляет шаблон по ID.
func (c *Client) DeleteTemplate(ctx context.Context, id int64) (*RequestResponse, error) {
	return do[RequestResponse](ctx, c.base, http.MethodDelete, "/templates/"+itoa(id), nil, nil)
}

// DeleteTemplateByAlias удаляет шаблон по алиасу.
func (c *Client) DeleteTemplateByAlias(ctx context.Context, alias string) (*RequestResponse, error) {
	return do[RequestResponse](ctx, c.base, http.MethodDelete, "/templates/"+url.PathEscape(alias), nil, nil)
}

// ========== Исходящие сообщения ==========

// GetMessages ищет исходящие сообщения.
func (c *Client) GetMessages(ctx context.Context, p *Params) (*OutboundMessages, error) {
	return do[OutboundMessages](ctx, c.base, http.MethodGet, "/messages", p, nil)
}

// GetMessageDetails возвращает подробности сообщения.
func (c *Client) GetMessageDetails(ctx context.Context, messageID string) (*OutboundMessageDetails, error) {
	return do[OutboundMessageDetails](ctx, c.base, http.MethodGet, "/messages/"+url.PathEscape(messageID)+"/details", nil, nil)
}

// GetMessageOpens возвращает открытия по всем сообщениям.
func (c *Client) GetMessageOpens(ctx context.Context, p *Params) (*OutboundMessageOpens, error) {
	return do[OutboundMessageOpens](ctx, c.base, http.MethodGet, "/messages/opens", p, nil)
}

// GetMessageOpensByID возвращает открытия одного сообщения.
func (c *Client) GetMessageOpensByID(ctx context.Context, messageID string, p *Params) (*OutboundMessageOpens, error) {
	return do[OutboundMessageOpens](ctx, c.base, http.MethodGet, "/messages/opens/"+url.PathEscape(messageID), p, nil)
}

// GetMessageClicks возвращает клики по всем сообщениям.
func (c *Client) GetMessageClicks(ctx context.Context, p *Params) (*OutboundMessageClicks, error) {
	return do[OutboundMessageClicks](ctx, c.base, http.MethodGet, "/messages/clicks", p, nil)
}

// GetMessageClicksByID возвращает клики одного сообщения.
func (c *Client) GetMessageClicksByID(ctx context.Context, messageID string, p *Params) (*OutboundMessageClicks, error) {
	return do[OutboundMessageClicks](ctx, c.base, http.MethodGet, "/messages/clicks/"+url.PathEscape(messageID), p, nil)
}

// ========== Статистика ==========

// GetOutboundStats возвращает сводную статистику.
func (c *Client) GetOutboundStats(ctx context.Context, p *Params) (*OutboundStats, error) {
	return do[OutboundStats](ctx, c.base, http.MethodGet, "/stats", p, nil)
}

// GetSendStats возвращает число отправленных писем по дням.
func (c *Client) GetSendStats(ctx context.Context, p *Params) (*OutboundSendStats, error) {
	return do[OutboundSendStats](ctx, c.base, http.MethodGet, "/stats/sends", p, nil)
}

// GetBounceStats возвращает отказы по дням.
func (c *Client) GetBounceStats(ctx context.Context, p *Params) (*OutboundBounceStats, error) {
	return do[OutboundBounceStats](ctx, c.base, http.MethodGet, "/stats/bounces", p, nil)
}

// GetSpamStats возвращает жалобы на спам по дням.
func (c *Client) GetSpamStats(ctx context.Context, p *Params) (*OutboundSpamStats, error) {
	return do[OutboundSpamStats](ctx, c.base, http.MethodGet, "/stats/spam", p, nil)
}

// GetOpenStats возвращает открытия по дням.
func (c *Client) GetOpenStats(ctx context.Context, p *Params) (*OutboundOpenStats, error) {
	return do[OutboundOpenStats](ctx, c.base, http.MethodGet, "/stats/opens", p, nil)
}

// GetOpenPlatformStats возвращает открытия по платформам.
func (c *Client) GetOpenPlatformStats(ctx context.Context, p *Params) (*OutboundOpenPlatformStats, error) {
	return do[OutboundOpenPlatformStats](ctx, c.base, http.MethodGet, "/stats/opens/platforms", p, nil)
}

// GetClickStats возвращает клики по дням.
func (c *Client) GetClickStats(ctx context.Context, p *Params) (*OutboundClickStats, error) {
	return do[OutboundClickStats](ctx, c.base, http.MethodGet, "/stats/clicks", p, nil)
}

// GetClickLocationStats возвращает клики по расположению (HTML/текст).
func (c *Client) GetClickLocationStats(ctx context.Context, p *Params) (*OutboundClickLocationStats, error) {
	return do[OutboundClickLocationStats](ctx, c.base, http.MethodGet, "/stats/clicks/location", p, nil)
}

// GetClickPlatformStats возвращает клики по платформам.
func (c *Client) GetClickPlatformStats(ctx context.Context, p *Params) (*OutboundClickPlatformStats, error) {
	return do[OutboundClickPlatformStats](ctx, c.base, http.MethodGet, "/stats/clicks/platforms", p, nil)
}

// ========== Вебхуки ==========

// GetWebhooks возвращает список вебхуков (можно фильтровать по MessageStream).
func (c *Client) GetWebhooks(ctx context.Context, p *Params) (*Webhooks, error) {
	return do[Webhooks](ctx, c.base, http.MethodGet, "/webhooks", p, nil)
}

// GetWebhook возвращает вебхук по ID.
func (c *Client) GetWebhook(ctx context.Context, id int64) (*Webhook, error) {
	return do[Webhook](ctx, c.base, http.MethodGet, "/webhooks/"+itoa(id), nil, nil)
}

// CreateWebhook создаёт вебхук.
func (c *Client) CreateWebhook(ctx context.Context, w *Webhook) (*Webhook, error) {
	return do[Webhook](ctx, c.base, http.MethodPost, "/webhooks", nil, w)
}

// SetWebhook изменяет вебхук.
func (c *Client) SetWebhook(ctx context.Context, id int64, w *Webhook) (*Webhook, error) {
	return do[Webhook](ctx, c.base, http.MethodPut, "/webhooks/"+itoa(id), nil, w)
}

// DeleteWebhook удаляет вебхук.
func (c *Client) DeleteWebhook(ctx context.Context, id int64) (*RequestResponse, error) {
	return do[RequestResponse](ctx, c.base, http.MethodDelete, "/webhooks/"+itoa(id), nil, nil)
}

// ========== Стоп-списки ==========

func suppressionsPath(stream string) string {
	return "/message-streams/" + url.PathEscape(stream) + "/suppressions"
}

// CreateSuppressions добавляет адреса в стоп-лист канала.
func (c *Client) CreateSuppressions(ctx context.Context, messageStream string, e SuppressionEntries) (*SuppressionStatuses, error) {
	return do[SuppressionStatuses](ctx, c.base, http.MethodPost, suppressionsPath(messageStream), nil, e)
}

// GetSuppressions возвращает стоп-лист канала.
func (c *Client) GetSuppressions(ctx context.Context, messageStream string, p *Params) (*Suppressions, error) {
	return do[Suppressions](ctx, c.base, http.MethodGet, suppressionsPath(messageStream), p, nil)
}

// DeleteSuppressions удаляет адреса из стоп-листа канала.
func (c *Client) DeleteSuppressions(ctx context.Context, messageStream string, e SuppressionEntries) (*SuppressionStatuses, error) {
	return do[SuppressionStatuses](ctx, c.base, http.MethodPost, suppressionsPath(messageStream)+"/delete", nil, e)
}

// ========== Каналы ==========

// GetMessageStreams возвращает каналы сервера.
func (c *Client) GetMessageStreams(ctx context.Context) (*MessageStreams, error) {
	return do[MessageStreams](ctx, c.base, http.MethodGet, "/message-streams", nil, nil)
}

// GetMessageStream возвращает канал по ID.
func (c *Client) GetMessageStream(ctx context.Context, id string) (*MessageStream, error) {
	return do[MessageStream](ctx, c.base, http.MethodGet, "/message-streams/"+url.PathEscape(id), nil, nil)
}

// CreateMessageStream создаёт канал.
func (c *Client) CreateMessageStream(ctx context.Context, s *MessageStream) (*MessageStream, error) {
	return do[MessageStream](ctx, c.base, http.MethodPost, "/message-streams", nil, s)
}

// SetMessageStream изменяет канал (PATCH).
func (c *Client) SetMessageStream(ctx context.Context, id string, s *MessageStream) (*MessageStream, error) {
	return do[MessageStream](ctx, c.base, http.MethodPatch, "/message-streams/"+url.PathEscape(id), nil, s)
}

// ArchiveMessageStream архивирует канал.
func (c *Client) ArchiveMessageStream(ctx context.Context, id string) (*MessageStreamArchiveResponse, error) {
	return do[MessageStreamArchiveResponse](ctx, c.base, http.MethodPost, "/message-streams/"+url.PathEscape(id)+"/archive", nil, nil)
}

// UnarchiveMessageStream восстанавливает канал из архива.
func (c *Client) UnarchiveMessageStream(ctx context.Context, id string) (*MessageStreamUnarchiveResponse, error) {
	return do[MessageStreamUnarchiveResponse](ctx, c.base, http.MethodPost, "/message-streams/"+url.PathEscape(id)+"/unarchive", nil, nil)
}
