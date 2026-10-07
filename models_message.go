package haskimail

import (
	"encoding/base64"
	"mime"
	"net/http"
	"os"
	"path/filepath"
	"sort"
	"strings"
)

// Режимы отслеживания ссылок (поле TrackLinks). В письме None выключает
// отслеживание, любое другое значение включает его. В настройках сервера
// различаются все четыре: None, HtmlAndText, HtmlOnly, TextOnly.
const (
	TrackLinksNone        = "None"
	TrackLinksHTMLAndText = "HtmlAndText"
	TrackLinksHTMLOnly    = "HtmlOnly"
	TrackLinksTextOnly    = "TextOnly"

	// Deprecated: используйте TrackLinksHTMLOnly. API принимает и это значение.
	TrackLinksHTML = "Html"
	// Deprecated: используйте TrackLinksTextOnly. API принимает и это значение.
	TrackLinksText = "Text"
)

// Header — пользовательский заголовок письма.
type Header struct {
	Name  string `json:"Name"`
	Value string `json:"Value"`
}

// Attachment — вложение письма. Content — содержимое в Base64.
type Attachment struct {
	Name          string `json:"Name,omitempty"`
	Content       string `json:"Content,omitempty"`
	ContentType   string `json:"ContentType,omitempty"`
	ContentID     string `json:"ContentId,omitempty"`
	ContentLength int    `json:"ContentLength,omitempty"`
}

// NewAttachment создаёт вложение из байтов (кодирует их в Base64).
// contentID задаётся для встраиваемых (inline) изображений, иначе пустая строка.
func NewAttachment(name string, content []byte, contentType, contentID string) Attachment {
	return Attachment{
		Name:        name,
		Content:     base64.StdEncoding.EncodeToString(content),
		ContentType: contentType,
		ContentID:   contentID,
	}
}

// AttachmentFromFile читает файл с диска и создаёт вложение.
// Тип содержимого определяется по расширению, а если не удалось — по содержимому.
func AttachmentFromFile(path, contentID string) (Attachment, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return Attachment{}, err
	}
	ct := mime.TypeByExtension(strings.ToLower(filepath.Ext(path)))
	if ct == "" {
		ct = http.DetectContentType(data)
	}
	return NewAttachment(filepath.Base(path), data, ct, contentID), nil
}

// Recipient — получатель с отображаемым именем.
type Recipient struct {
	Name  string `json:"Name,omitempty"`
	Email string `json:"Email,omitempty"`
}

// FormatRecipients собирает строку адресатов вида "Имя"<email>,"Имя2"<email2>
// для полей To, Cc и Bcc.
func FormatRecipients(rs ...Recipient) string {
	parts := make([]string, 0, len(rs))
	for _, r := range rs {
		if r.Name == "" {
			parts = append(parts, r.Email)
			continue
		}
		parts = append(parts, `"`+r.Name+`"<`+r.Email+`>`)
	}
	return strings.Join(parts, ",")
}

// FormatRecipientsMap собирает строку адресатов из карты: ключ — имя,
// значение — email. Порядок детерминирован (сортировка по имени).
func FormatRecipientsMap(m map[string]string) string {
	names := make([]string, 0, len(m))
	for n := range m {
		names = append(names, n)
	}
	sort.Strings(names)
	rs := make([]Recipient, 0, len(m))
	for _, n := range names {
		rs = append(rs, Recipient{Name: n, Email: m[n]})
	}
	return FormatRecipients(rs...)
}

// Message — письмо для отправки через /email.
type Message struct {
	MessageStream string            `json:"MessageStream,omitempty"`
	From          string            `json:"From,omitempty"`
	To            string            `json:"To,omitempty"`
	Cc            string            `json:"Cc,omitempty"`
	Bcc           string            `json:"Bcc,omitempty"`
	ReplyTo       string            `json:"ReplyTo,omitempty"`
	Subject       string            `json:"Subject,omitempty"`
	HTMLBody      string            `json:"HtmlBody,omitempty"`
	TextBody      string            `json:"TextBody,omitempty"`
	Tag           string            `json:"Tag,omitempty"`
	Headers       []Header          `json:"Headers,omitempty"`
	Attachments   []Attachment      `json:"Attachments,omitempty"`
	TrackOpens    *bool             `json:"TrackOpens,omitempty"`
	TrackLinks    string            `json:"TrackLinks,omitempty"`
	Metadata      map[string]string `json:"Metadata,omitempty"`
}

// AddHeader добавляет заголовок.
func (m *Message) AddHeader(name, value string) {
	m.Headers = append(m.Headers, Header{Name: name, Value: value})
}

// AddAttachment добавляет вложение.
func (m *Message) AddAttachment(a Attachment) {
	m.Attachments = append(m.Attachments, a)
}

// AddMetadata добавляет пару ключ-значение в метаданные.
func (m *Message) AddMetadata(key, value string) {
	if m.Metadata == nil {
		m.Metadata = map[string]string{}
	}
	m.Metadata[key] = value
}

// MessageResponse — результат отправки письма.
type MessageResponse struct {
	ErrorCode   int    `json:"ErrorCode"`
	Message     string `json:"Message"`
	To          string `json:"To,omitempty"`
	Cc          string `json:"Cc,omitempty"`
	Bcc         string `json:"Bcc,omitempty"`
	SubmittedAt *Time  `json:"SubmittedAt,omitempty"`
	MessageID   string `json:"MessageID,omitempty"`
}

// RequestResponse — стандартный ответ API на операции без данных (например, удаление).
type RequestResponse struct {
	ErrorCode int    `json:"ErrorCode"`
	Message   string `json:"Message"`
}
