package haskimail

// HTTPAuth — Basic-аутентификация для вызова вебхука.
type HTTPAuth struct {
	Username string `json:"Username,omitempty"`
	Password string `json:"Password,omitempty"`
}

// WebhookTrigger — базовый триггер (клик, доставка).
type WebhookTrigger struct {
	Enabled *bool `json:"Enabled,omitempty"`
}

// OpenWebhookTrigger — триггер открытий.
type OpenWebhookTrigger struct {
	Enabled           *bool `json:"Enabled,omitempty"`
	PostFirstOpenOnly *bool `json:"PostFirstOpenOnly,omitempty"`
}

// BounceWebhookTrigger — триггер отказов и жалоб на спам.
type BounceWebhookTrigger struct {
	Enabled        *bool `json:"Enabled,omitempty"`
	IncludeContent *bool `json:"IncludeContent,omitempty"`
}

// SpamWebhookTrigger — триггер жалоб на спам.
type SpamWebhookTrigger = BounceWebhookTrigger

// SubscriptionChange — триггер изменения подписки.
type SubscriptionChange = WebhookTrigger

// WebhookTriggers — набор триггеров вебхука.
type WebhookTriggers struct {
	Open               *OpenWebhookTrigger   `json:"Open,omitempty"`
	Click              *WebhookTrigger       `json:"Click,omitempty"`
	Delivery           *WebhookTrigger       `json:"Delivery,omitempty"`
	Bounce             *BounceWebhookTrigger `json:"Bounce,omitempty"`
	SpamComplaint      *SpamWebhookTrigger   `json:"SpamComplaint,omitempty"`
	SubscriptionChange *SubscriptionChange   `json:"SubscriptionChange,omitempty"`
}

// Webhook — настройка вебхука.
type Webhook struct {
	ID            int64            `json:"ID,omitempty"`
	URL           string           `json:"Url,omitempty"`
	HTTPAuth      *HTTPAuth        `json:"HttpAuth,omitempty"`
	HTTPHeaders   []Header         `json:"HttpHeaders,omitempty"`
	MessageStream string           `json:"MessageStream,omitempty"`
	Triggers      *WebhookTriggers `json:"Triggers,omitempty"`
}

// Webhooks — список вебхуков.
type Webhooks struct {
	Webhooks []Webhook `json:"Webhooks"`
}

// Типы входящих событий вебхуков (для разбора тела запроса в вашем обработчике).

// DeliveryWebhook — событие доставки.
type DeliveryWebhook struct {
	ServerID      int64             `json:"ServerID"`
	MessageID     string            `json:"MessageID"`
	Recipient     string            `json:"Recipient"`
	Tag           string            `json:"Tag"`
	DeliveredAt   *Time             `json:"DeliveredAt"`
	Details       string            `json:"Details"`
	RecordType    string            `json:"RecordType"`
	Metadata      map[string]string `json:"Metadata"`
	MessageStream string            `json:"MessageStream"`
}

// BounceWebhook — событие отказа.
type BounceWebhook struct {
	Bounce
	Metadata map[string]string `json:"Metadata"`
}

// OpenWebhook — событие открытия.
type OpenWebhook struct {
	OutboundMessageOpen
	FirstOpen bool              `json:"FirstOpen"`
	Metadata  map[string]string `json:"Metadata"`
}

// ClickWebhook — событие клика.
type ClickWebhook struct {
	OutboundMessageClick
	Metadata map[string]string `json:"Metadata"`
}
