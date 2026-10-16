package haskimail

// Типы шаблонов.
const (
	TemplateTypeStandard = "Standard"
	TemplateTypeLayout   = "Layout"
)

// BaseTemplate — краткие сведения о шаблоне (в списке шаблонов).
type BaseTemplate struct {
	TemplateID     int64  `json:"TemplateId,omitempty"`
	Alias          string `json:"Alias,omitempty"`
	Name           string `json:"Name,omitempty"`
	Active         *bool  `json:"Active,omitempty"`
	TemplateType   string `json:"TemplateType,omitempty"`
	LayoutTemplate string `json:"LayoutTemplate,omitempty"`
}

// Template — шаблон целиком; используется и в запросах создания/изменения.
type Template struct {
	BaseTemplate
	Subject            string `json:"Subject,omitempty"`
	HTMLBody           string `json:"HtmlBody,omitempty"`
	TextBody           string `json:"TextBody,omitempty"`
	AssociatedServerID int64  `json:"AssociatedServerId,omitempty"`
}

// Templates — страница списка шаблонов.
type Templates struct {
	TotalCount int            `json:"TotalCount"`
	Templates  []BaseTemplate `json:"Templates"`
}

// TemplateContent — содержимое шаблона.
type TemplateContent struct {
	Subject  string `json:"Subject,omitempty"`
	HTMLBody string `json:"HtmlBody,omitempty"`
	TextBody string `json:"TextBody,omitempty"`
}

// TemplateToValidate — шаблон для проверки (ValidateTemplate). Нужен хотя бы один из
// Subject, HTMLBody, TextBody. TestRenderModel подставляется буквально, как TemplateModel
// при отправке: ключ — фрагмент вместе с разделителями, map[string]string{"{{name}}": "Иван"}.
type TemplateToValidate struct {
	Subject         string `json:"Subject,omitempty"`
	HTMLBody        string `json:"HtmlBody,omitempty"`
	TextBody        string `json:"TextBody,omitempty"`
	TestRenderModel any    `json:"TestRenderModel,omitempty"`
}

// TemplateValidationError — фрагмент {{…}} без значения или «{{» без пары. Line и
// CharacterPosition считаются с 1 в исходном тексте; 0 — фрагмент принесло значение модели.
type TemplateValidationError struct {
	Message           string `json:"Message"`
	Line              int    `json:"Line"`
	CharacterPosition int    `json:"CharacterPosition"`
}

// TemplateValidationField — результат проверки одной части шаблона.
type TemplateValidationField struct {
	ContentIsValid   bool                      `json:"ContentIsValid"`
	ValidationErrors []TemplateValidationError `json:"ValidationErrors"`
	// RenderedContent — часть после подстановки, как её получит получатель.
	RenderedContent string `json:"RenderedContent"`
}

// TemplateValidation — результат проверки шаблона. Части, которой не было в запросе,
// нет и в ответе (nil).
type TemplateValidation struct {
	AllContentIsValid bool                     `json:"AllContentIsValid"`
	Subject           *TemplateValidationField `json:"Subject,omitempty"`
	HTMLBody          *TemplateValidationField `json:"HtmlBody,omitempty"`
	TextBody          *TemplateValidationField `json:"TextBody,omitempty"`
	// SuggestedTemplateModel — все {{…}} шаблона с значениями из TestRenderModel
	// или заглушками «имя_Value».
	SuggestedTemplateModel map[string]string `json:"SuggestedTemplateModel"`
}

// TemplatedMessage — письмо, отправляемое по шаблону.
// Укажите TemplateID или TemplateAlias.
type TemplatedMessage struct {
	TemplateID    int64             `json:"TemplateId,omitempty"`
	TemplateAlias string            `json:"TemplateAlias,omitempty"`
	TemplateModel any               `json:"TemplateModel,omitempty"`
	InlineCSS     *bool             `json:"InlineCss,omitempty"`
	From          string            `json:"From,omitempty"`
	To            string            `json:"To,omitempty"`
	Cc            string            `json:"Cc,omitempty"`
	Bcc           string            `json:"Bcc,omitempty"`
	ReplyTo       string            `json:"ReplyTo,omitempty"`
	Tag           string            `json:"Tag,omitempty"`
	Headers       []Header          `json:"Headers,omitempty"`
	Attachments   []Attachment      `json:"Attachments,omitempty"`
	MessageStream string            `json:"MessageStream,omitempty"`
	TrackOpens    *bool             `json:"TrackOpens,omitempty"`
	TrackLinks    string            `json:"TrackLinks,omitempty"`
	Metadata      map[string]string `json:"Metadata,omitempty"`
}

// AddHeader добавляет заголовок.
func (m *TemplatedMessage) AddHeader(name, value string) {
	m.Headers = append(m.Headers, Header{Name: name, Value: value})
}

// AddAttachment добавляет вложение.
func (m *TemplatedMessage) AddAttachment(a Attachment) {
	m.Attachments = append(m.Attachments, a)
}

// AddMetadata добавляет пару ключ-значение в метаданные.
func (m *TemplatedMessage) AddMetadata(key, value string) {
	if m.Metadata == nil {
		m.Metadata = map[string]string{}
	}
	m.Metadata[key] = value
}

// TemplatesPushRequest — запрос на перенос шаблонов между серверами.
// Имена полей — как в API: сервер читает их с учётом регистра.
type TemplatesPushRequest struct {
	SourceServerID      int64 `json:"SourceServerID"`
	DestinationServerID int64 `json:"DestinationServerID"`
	PerformChanges      *bool `json:"PerformChanges,omitempty"`
}

// TemplatesPushAction — действие над шаблоном при пуше.
type TemplatesPushAction struct {
	Action     string `json:"Action"`
	TemplateID int64  `json:"TemplateId"`
	Alias      string `json:"Alias"`
	Name       string `json:"Name"`
}

// TemplatesPush — результат пуша шаблонов.
type TemplatesPush struct {
	TotalCount int                   `json:"TotalCount"`
	Templates  []TemplatesPushAction `json:"Templates"`
}
