package haskimail

// OutboundStats — сводная статистика исходящих писем.
type OutboundStats struct {
	Sent                  int     `json:"Sent"`
	Bounced               int     `json:"Bounced"`
	SMTPAPIErrors         int     `json:"SmtpApiErrors"`
	BounceRate            float64 `json:"BounceRate"`
	SpamComplaintsRate    float64 `json:"SpamComplaintsRate"`
	Opens                 int     `json:"Opens"`
	UniqueOpens           int     `json:"UniqueOpens"`
	Tracked               int     `json:"Tracked"`
	WithLinkTracking      int     `json:"WithLinkTracking"`
	WithOpenTracking      int     `json:"WithOpenTracking"`
	TotalTrackedLinksSent int     `json:"TotalTrackedLinksSent"`
	UniqueLinksClicked    int     `json:"UniqueLinksClicked"`
	TotalClicks           int     `json:"TotalClicks"`
	WithClientRecorded    int     `json:"WithClientRecorded"`
	WithPlatformRecorded  int     `json:"WithPlatformRecorded"`
	WithReadTimeRecorded  int     `json:"WithReadTimeRecorded"`
	SpamComplaints        int     `json:"SpamComplaints"`
	// Opened и Clicked — открытия и клики за период, включая повторные.
	Opened  int `json:"Opened"`
	Clicked int `json:"Clicked"`
	// Unsubscribed — адресов, отписавшихся за период (в API поле называется Unsubsribed).
	Unsubscribed int     `json:"Unsubsribed"`
	OpenRate     float64 `json:"OpenRate"`
	ClickRate    float64 `json:"ClickRate"`
}

// SentStat — отправлено за день.
type SentStat struct {
	Date *Time `json:"Date"`
	Sent int   `json:"Sent"`
}

// BounceStat — отказы за день по типам.
type BounceStat struct {
	Date         *Time `json:"Date"`
	HardBounce   int   `json:"HardBounce"`
	SoftBounce   int   `json:"SoftBounce"`
	SMTPAPIError int   `json:"SmtpApiError"`
	Transient    int   `json:"Transient"`
	Unknown      int   `json:"Unknown"`
	DMARCPolicy  int   `json:"DMARCPolicy"`
	Subscribe    int   `json:"Subscribe"`
}

// SpamStat — жалобы на спам за день.
type SpamStat struct {
	Date          *Time `json:"Date"`
	SpamComplaint int   `json:"SpamComplaint"`
}

// OpenStat — открытия за день.
type OpenStat struct {
	Date   *Time `json:"Date"`
	Opens  int   `json:"Opens"`
	Unique int   `json:"Unique"`
}

// ClickStat — клики за день.
type ClickStat struct {
	Date   *Time `json:"Date"`
	Clicks int   `json:"Clicks"`
	Unique int   `json:"Unique"`
}

// DayStats — статистика с разбивкой по дням.
type DayStats[T any] struct {
	Days []T `json:"Days"`
}

// Псевдонимы для конкретных видов статистики.
type (
	OutboundSendStats   = DayStats[SentStat]
	OutboundBounceStats = DayStats[BounceStat]
	OutboundSpamStats   = DayStats[SpamStat]
	OutboundOpenStats   = DayStats[OpenStat]
	OutboundClickStats  = DayStats[ClickStat]
)
