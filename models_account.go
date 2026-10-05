package haskimail

// Server — сервер (набор настроек отправки) в аккаунте.
type Server struct {
	ID                         int64    `json:"Id,omitempty"`
	Name                       string   `json:"Name,omitempty"`
	Color                      string   `json:"Color,omitempty"`
	APITokens                  []string `json:"ApiTokens,omitempty"`
	DeliveryHookURL            string   `json:"DeliveryHookUrl,omitempty"`
	BounceHookURL              string   `json:"BounceHookUrl,omitempty"`
	OpenHookURL                string   `json:"OpenHookUrl,omitempty"`
	ClickHookURL               string   `json:"ClickHookUrl,omitempty"`
	TrackOpens                 *bool    `json:"TrackOpens,omitempty"`
	TrackLinks                 string   `json:"TrackLinks,omitempty"`
	PostFirstOpenOnly          *bool    `json:"PostFirstOpenOnly,omitempty"`
	InboundAddress             string   `json:"InboundAddress,omitempty"`
	InboundDomain              string   `json:"InboundDomain,omitempty"`
	InboundHash                string   `json:"InboundHash,omitempty"`
	InboundSpamThreshold       *int     `json:"InboundSpamThreshold,omitempty"`
	SMTPAPIActivated           *bool    `json:"SmtpApiActivated,omitempty"`
	RawEmailEnabled            *bool    `json:"RawEmailEnabled,omitempty"`
	EnableSMTPAPIErrorHooks    *bool    `json:"EnableSmtpApiErrorHooks,omitempty"`
	IncludeBounceContentInHook *bool    `json:"IncludeBounceContentInHook,omitempty"`
	ServerLink                 string   `json:"ServerLink,omitempty"`
	DeliveryType               string   `json:"DeliveryType,omitempty"`
}

// Servers — страница списка серверов.
type Servers struct {
	TotalCount int      `json:"TotalCount"`
	Servers    []Server `json:"Servers"`
}

// Domain — домен отправителя. В запросах создания/изменения
// значимо в основном поле Name; флаги проверки только для чтения.
type Domain struct {
	ID                       int64  `json:"Id,omitempty"`
	Name                     string `json:"Name,omitempty"`
	SPFVerified              bool   `json:"SpfVerified,omitempty"`
	DKIMVerified             bool   `json:"DkimVerified,omitempty"`
	WeakDKIM                 bool   `json:"WeakDKIM,omitempty"`
	ReturnPathDomainVerified bool   `json:"ReturnPathDomainVerified,omitempty"`
}

// DomainDetails — полные сведения о домене, включая DNS-записи.
type DomainDetails struct {
	Domain
	DKIMHost                      string `json:"DkimHost"`
	DKIMTextValue                 string `json:"DkimTextValue"`
	DKIMPendingHost               string `json:"DkimPendingHost"`
	DKIMPendingTextValue          string `json:"DkimPendingTextValue"`
	DKIMRevokedHost               string `json:"DkimRevokedHost"`
	DKIMRevokedTextValue          string `json:"DkimRevokedTextValue"`
	DKIMUpdateStatus              string `json:"DkimUpdateStatus"`
	ReturnPathDomain              string `json:"ReturnPathDomain"`
	ReturnPathDomainCNAMEValue    string `json:"ReturnPathDomainCNAMEValue"`
	SafeToRemoveRevokedKeyFromDNS bool   `json:"SafeToRemoveRevokedKeyFromDNS"`
}

// Domains — список доменов.
type Domains struct {
	TotalCount int      `json:"TotalCount"`
	Domains    []Domain `json:"Domains"`
}

// Signature — подпись отправителя.
type Signature struct {
	ID                  int64  `json:"Id"`
	Domain              string `json:"Domain"`
	EmailAddress        string `json:"EmailAddress"`
	ReplyToEmailAddress string `json:"ReplyToEmailAddress"`
	Name                string `json:"Name"`
	Confirmed           bool   `json:"Confirmed"`
}

// SignatureDetails — полные сведения о подписи.
type SignatureDetails struct {
	Signature
	DKIMHost                   string `json:"DkimHost"`
	DKIMTextValue              string `json:"DkimTextValue"`
	DKIMVerified               bool   `json:"DkimVerified"`
	SPFVerified                bool   `json:"SpfVerified"`
	ReturnPathDomain           string `json:"ReturnPathDomain"`
	ReturnPathDomainCNAMEValue string `json:"ReturnPathDomainCNAMEValue"`
	ReturnPathDomainVerified   bool   `json:"ReturnPathDomainVerified"`
	ConfirmationPersonalNote   string `json:"ConfirmationPersonalNote"`
}

// Signatures — страница списка подписей.
type Signatures struct {
	TotalCount       int         `json:"TotalCount"`
	SenderSignatures []Signature `json:"SenderSignatures"`
}

// SignatureToCreate — запрос создания/изменения подписи.
type SignatureToCreate struct {
	FromEmail                string `json:"FromEmail,omitempty"`
	Name                     string `json:"Name,omitempty"`
	ReplyToEmail             string `json:"ReplyToEmail,omitempty"`
	ReturnPathDomain         string `json:"ReturnPathDomain,omitempty"`
	ConfirmationPersonalNote string `json:"ConfirmationPersonalNote,omitempty"`
}
