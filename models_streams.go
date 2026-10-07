package haskimail

// Типы каналов.
const (
	MessageStreamTransactional = "Transactional"
	MessageStreamBroadcasts    = "Broadcasts"
	MessageStreamInbound       = "Inbound"
)

// SubscriptionManagementConfiguration — настройки обработки отписок.
type SubscriptionManagementConfiguration struct {
	// UnsubscribeHandlingType: "Haskimail" (по умолчанию) или "Custom".
	UnsubscribeHandlingType string `json:"UnsubscribeHandlingType,omitempty"`
}

// MessageStream — канал отправки. ID — числовой ID канала строкой («2081»),
// его передают в Message.MessageStream. Имена JSON-полей совпадают с API:
// при создании канала сервер читает ServerID с учётом регистра.
type MessageStream struct {
	ID                                  string                               `json:"ID,omitempty"`
	ServerID                            int64                                `json:"ServerID,omitempty"`
	Name                                string                               `json:"Name,omitempty"`
	Description                         string                               `json:"Description,omitempty"`
	MessageStreamType                   string                               `json:"MessageStreamType,omitempty"`
	CreatedAt                           *Time                                `json:"CreatedAt,omitempty"`
	UpdatedAt                           *Time                                `json:"UpdatedAt,omitempty"`
	ArchivedAt                          *Time                                `json:"ArchivedAt,omitempty"`
	SubscriptionManagementConfiguration *SubscriptionManagementConfiguration `json:"SubscriptionManagementConfiguration,omitempty"`
}

// MessageStreams — список каналов.
type MessageStreams struct {
	TotalCount     int             `json:"TotalCount"`
	MessageStreams []MessageStream `json:"MessageStreams"`
}

// MessageStreamArchiveResponse — результат архивации канала.
type MessageStreamArchiveResponse struct {
	ID                string `json:"ID"`
	ServerID          int64  `json:"ServerID"`
	ArchivedAt        *Time  `json:"ArchivedAt"`
	ExpectedPurgeDate *Time  `json:"ExpectedPurgeDate"`
}

// MessageStreamUnarchiveResponse — результат разархивации канала.
type MessageStreamUnarchiveResponse = MessageStream
