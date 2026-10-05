package haskimail

// OutboundMessage — исходящее сообщение в поиске.
type OutboundMessage struct {
	Recipient     string `json:"Recipient"`
	ReceivedAt    *Time  `json:"ReceivedAt"`
	Tag           string `json:"Tag"`
	From          string `json:"From"`
	Status        string `json:"Status"`
	MessageID     string `json:"MessageID"`
	Subject       string `json:"Subject"`
	MessageStream string `json:"MessageStream"`
}

// OutboundMessages — страница результатов поиска сообщений.
type OutboundMessages struct {
	TotalCount int               `json:"TotalCount"`
	Messages   []OutboundMessage `json:"Messages"`
}

// OutboundMessageDetailsEvent — событие жизненного цикла сообщения.
type OutboundMessageDetailsEvent struct {
	Recipient  string            `json:"Recipient"`
	Type       string            `json:"Type"`
	ReceivedAt *Time             `json:"ReceivedAt"`
	Details    map[string]string `json:"Details"`
}

// OutboundMessageDetails — подробности сообщения.
type OutboundMessageDetails struct {
	OutboundMessage
	TextBody      string                        `json:"TextBody"`
	HTMLBody      string                        `json:"HtmlBody"`
	Body          string                        `json:"Body"`
	MessageEvents []OutboundMessageDetailsEvent `json:"MessageEvents"`
}

// OutboundMessageOpen — событие открытия письма.
type OutboundMessageOpen struct {
	Recipient     string `json:"Recipient"`
	ReceivedAt    *Time  `json:"ReceivedAt"`
	Tag           string `json:"Tag"`
	From          string `json:"From"`
	RecordType    string `json:"RecordType"`
	OS            string `json:"Os"`
	Browser       string `json:"Browser"`
	Webservice    string `json:"Webservice"`
	UserAgent     string `json:"UserAgent"`
	Language      string `json:"Language"`
	Region        string `json:"Region"`
	Country       string `json:"Country"`
	MessageID     string `json:"MessageID"`
	Subject       string `json:"Subject"`
	MessageStream string `json:"MessageStream"`
}

// OutboundMessageOpens — страница открытий.
type OutboundMessageOpens struct {
	TotalCount int                   `json:"TotalCount"`
	Opens      []OutboundMessageOpen `json:"Opens"`
}

// OutboundMessageClick — событие клика по ссылке.
type OutboundMessageClick struct {
	OutboundMessageOpen
	OriginalLink  string `json:"OriginalLink"`
	ClickLocation string `json:"ClickLocation"`
}

// OutboundMessageClicks — страница кликов.
type OutboundMessageClicks struct {
	TotalCount int                    `json:"TotalCount"`
	Clicks     []OutboundMessageClick `json:"Clicks"`
}
