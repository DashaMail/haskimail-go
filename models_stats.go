package haskimail

import (
	"encoding/json"
	"fmt"
)

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

// PlatformStat — открытия или клики за день по платформам: Desktop — компьютер,
// Mobile — телефон или планшет, WebMail — веб-почта (Gmail, Mail.ru, Яндекс, Рамблер),
// Unknown — не определилось, в том числе боты.
type PlatformStat struct {
	Date    *Time `json:"Date"`
	Desktop int   `json:"Desktop"`
	Mobile  int   `json:"Mobile"`
	WebMail int   `json:"WebMail"`
	Unknown int   `json:"Unknown"`
}

// PlatformStats — открытия или клики по платформам: по дням и итог за период.
type PlatformStats struct {
	Days    []PlatformStat `json:"Days"`
	Desktop int            `json:"Desktop"`
	Mobile  int            `json:"Mobile"`
	WebMail int            `json:"WebMail"`
	Unknown int            `json:"Unknown"`
}

// Псевдонимы для разбивки по платформам.
type (
	OutboundOpenPlatformStats  = PlatformStats
	OutboundClickPlatformStats = PlatformStats
)

// NamedCountStat — счётчики за день по названиям: почтовым клиентам или браузерам.
// В JSON названия — ключи объекта рядом с Date.
type NamedCountStat struct {
	Date   *Time
	Counts map[string]int
}

// UnmarshalJSON разбирает день: Date — дата, остальные ключи — названия со счётчиками.
func (s *NamedCountStat) UnmarshalJSON(b []byte) error {
	var raw map[string]json.RawMessage
	if err := json.Unmarshal(b, &raw); err != nil {
		return err
	}
	date, counts, err := splitNamedCounts(raw, "Date")
	if err != nil {
		return err
	}
	s.Counts = counts
	if date != nil {
		s.Date = new(Time)
		if err := json.Unmarshal(date, s.Date); err != nil {
			return err
		}
	}
	return nil
}

// NamedCountStats — открытия по почтовым клиентам или клики по браузерам: по дням
// и итог за период (Totals). В JSON итоги — ключи объекта рядом с Days.
type NamedCountStats struct {
	Days   []NamedCountStat
	Totals map[string]int
}

// UnmarshalJSON разбирает ответ: Days — дни, остальные ключи — итоги по названиям.
func (s *NamedCountStats) UnmarshalJSON(b []byte) error {
	var raw map[string]json.RawMessage
	if err := json.Unmarshal(b, &raw); err != nil {
		return err
	}
	days, totals, err := splitNamedCounts(raw, "Days")
	if err != nil {
		return err
	}
	s.Totals = totals
	s.Days = nil
	if days != nil {
		return json.Unmarshal(days, &s.Days)
	}
	return nil
}

// splitNamedCounts отделяет служебный ключ (Date или Days) от счётчиков по названиям.
func splitNamedCounts(raw map[string]json.RawMessage, key string) (json.RawMessage, map[string]int, error) {
	counts := make(map[string]int, len(raw))
	for name, value := range raw {
		if name == key {
			continue
		}
		var n int
		if err := json.Unmarshal(value, &n); err != nil {
			return nil, nil, fmt.Errorf("haskimail: счётчик %q: %w", name, err)
		}
		counts[name] = n
	}
	return raw[key], counts, nil
}

// Псевдонимы для разбивки по названиям.
type (
	OutboundOpenEmailClientStats = NamedCountStats
	OutboundClickBrowserStats    = NamedCountStats
)
