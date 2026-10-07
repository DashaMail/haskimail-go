package haskimail

import (
	"bytes"
	"fmt"
	"strconv"
	"strings"
	"time"
)

// TimeLocation — часовой пояс для дат API, пришедших без указания зоны
// (например "2024-01-02 15:04:05"). По умолчанию — локальный.
var TimeLocation = time.Local

var timeLayouts = []string{
	time.RFC3339Nano,
	"2006-01-02T15:04:05.999999999Z0700",
	"2006-01-02T15:04:05Z0700",
	"2006-01-02 15:04:05",
	"2006-01-02T15:04:05",
	"2006-01-02T15:04:05.999999999",
	time.RFC1123Z,
	"Mon, 2 Jan 2006 15:04:05 -0700",
	"2006-01-02", // дневная статистика (Days[].Date)
}

// Time — дата/время API. Разбирает все форматы, встречающиеся в ответах Haskimail.
type Time struct {
	time.Time
}

// ParseTime разбирает строку в одном из поддерживаемых форматов.
func ParseTime(s string) (time.Time, error) {
	for _, layout := range timeLayouts {
		if t, err := time.ParseInLocation(layout, s, TimeLocation); err == nil {
			return t, nil
		}
	}
	return time.Time{}, fmt.Errorf("haskimail: не удалось разобрать дату: %q", s)
}

// UnmarshalJSON реализует json.Unmarshaler.
func (t *Time) UnmarshalJSON(b []byte) error {
	if bytes.Equal(b, []byte("null")) {
		return nil
	}
	s, err := strconv.Unquote(string(b))
	if err != nil {
		return fmt.Errorf("haskimail: дата должна быть строкой: %s", b)
	}
	// Пустая дата MySQL — «даты нет», а не ошибка разбора.
	if s == "" || strings.HasPrefix(s, "0000-00-00") {
		return nil
	}
	parsed, err := ParseTime(s)
	if err != nil {
		return err
	}
	t.Time = parsed
	return nil
}

// MarshalJSON реализует json.Marshaler (RFC 3339).
func (t Time) MarshalJSON() ([]byte, error) {
	if t.IsZero() {
		return []byte("null"), nil
	}
	return []byte(strconv.Quote(t.Format(time.RFC3339Nano))), nil
}
