package haskimail

// Server — сервер (набор настроек отправки) в аккаунте.
type Server struct {
	ID                         int64    `json:"ID,omitempty"`
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
	ID                       int64  `json:"ID,omitempty"`
	Name                     string `json:"Name,omitempty"`
	SPFVerified              bool   `json:"SPFVerified,omitempty"`
	DKIMVerified             bool   `json:"DKIMVerified,omitempty"`
	WeakDKIM                 bool   `json:"WeakDKIM,omitempty"`
	ReturnPathDomainVerified bool   `json:"ReturnPathDomainVerified,omitempty"`
}

// DNSRecord — DNS-запись, которую нужно добавить домену.
type DNSRecord struct {
	// Name — имя записи с точкой на конце, например "dm._domainkey.example.com.".
	Name string `json:"name"`
	// RecordType — тип записи, всегда "TXT".
	RecordType string `json:"record_type"`
	// Value — значение записи.
	Value string `json:"value"`
	// Valid — 1, если запись в DNS найдена и верна.
	Valid int `json:"valid"`
}

// DomainDetails — полные сведения о домене, включая DNS-записи.
// SPFTextValue и DKIMTextValue API отдаёт объектами DNS-записи.
type DomainDetails struct {
	Domain
	SPFHost                       string     `json:"SPFHost"`
	SPFTextValue                  *DNSRecord `json:"SPFTextValue"`
	DKIMHost                      string     `json:"DKIMHost"`
	DKIMTextValue                 *DNSRecord `json:"DKIMTextValue"`
	DKIMPendingHost               string     `json:"DkimPendingHost"`
	DKIMPendingTextValue          string     `json:"DkimPendingTextValue"`
	DKIMRevokedHost               string     `json:"DkimRevokedHost"`
	DKIMRevokedTextValue          string     `json:"DkimRevokedTextValue"`
	DKIMUpdateStatus              string     `json:"DkimUpdateStatus"`
	ReturnPathDomain              string     `json:"ReturnPathDomain"`
	ReturnPathDomainCNAMEValue    string     `json:"ReturnPathDomainCNAMEValue"`
	SafeToRemoveRevokedKeyFromDNS bool       `json:"SafeToRemoveRevokedKeyFromDNS"`
}

// Domains — список доменов.
type Domains struct {
	TotalCount int      `json:"TotalCount"`
	Domains    []Domain `json:"Domains"`
}
