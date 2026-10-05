package haskimail

import (
	"net"
	"net/http"
	"time"
)

const (
	// Version — версия библиотеки.
	Version = "1.0.0"

	// DefaultAPIURL — хост API по умолчанию.
	DefaultAPIURL = "api.haskimail.ru"

	// ServerTokenHeader — заголовок для серверного токена.
	ServerTokenHeader = "X-Haski-Server-Token"

	// AccountTokenHeader — заголовок для токена аккаунта.
	AccountTokenHeader = "X-Haski-Account-Token"

	// DefaultConnectTimeout — таймаут установки соединения по умолчанию.
	DefaultConnectTimeout = 60 * time.Second

	// DefaultReadTimeout — таймаут ожидания ответа по умолчанию.
	DefaultReadTimeout = 60 * time.Second
)

// Option настраивает клиент при создании.
type Option func(*config)

type config struct {
	baseURL        string
	secure         bool
	httpClient     *http.Client
	userAgent      string
	connectTimeout time.Duration
	readTimeout    time.Duration
	headers        map[string]string
}

func defaultConfig() *config {
	return &config{
		baseURL:        DefaultAPIURL,
		secure:         true,
		userAgent:      "Haskimail-Go (version: " + Version + ")",
		connectTimeout: DefaultConnectTimeout,
		readTimeout:    DefaultReadTimeout,
	}
}

// WithBaseURL задаёт пользовательский хост API (без схемы), например "api.example.ru".
func WithBaseURL(host string) Option {
	return func(c *config) { c.baseURL = host }
}

// WithInsecure переключает клиент на HTTP вместо HTTPS.
func WithInsecure() Option {
	return func(c *config) { c.secure = false }
}

// WithSecureConnection явно задаёт протокол: true — HTTPS, false — HTTP.
func WithSecureConnection(secure bool) Option {
	return func(c *config) { c.secure = secure }
}

// WithHTTPClient задаёт собственный *http.Client. В этом случае
// WithConnectTimeout и WithReadTimeout не применяются.
func WithHTTPClient(hc *http.Client) Option {
	return func(c *config) { c.httpClient = hc }
}

// WithConnectTimeout задаёт таймаут установки TCP-соединения.
func WithConnectTimeout(d time.Duration) Option {
	return func(c *config) { c.connectTimeout = d }
}

// WithReadTimeout задаёт таймаут ожидания заголовков ответа.
func WithReadTimeout(d time.Duration) Option {
	return func(c *config) { c.readTimeout = d }
}

// WithUserAgent переопределяет заголовок User-Agent.
func WithUserAgent(ua string) Option {
	return func(c *config) { c.userAgent = ua }
}

// WithHeader добавляет произвольный заголовок ко всем запросам.
func WithHeader(name, value string) Option {
	return func(c *config) {
		if c.headers == nil {
			c.headers = map[string]string{}
		}
		c.headers[name] = value
	}
}

func (c *config) buildHTTPClient() *http.Client {
	if c.httpClient != nil {
		return c.httpClient
	}
	tr := http.DefaultTransport.(*http.Transport).Clone()
	tr.DialContext = (&net.Dialer{Timeout: c.connectTimeout, KeepAlive: 30 * time.Second}).DialContext
	tr.ResponseHeaderTimeout = c.readTimeout
	return &http.Client{Transport: tr}
}

// NewClient создаёт клиент серверного API.
func NewClient(serverToken string, opts ...Option) *Client {
	return &Client{base: newBase(ServerTokenHeader, serverToken, opts)}
}

// NewAccountClient создаёт клиент API аккаунта.
func NewAccountClient(accountToken string, opts ...Option) *AccountClient {
	return &AccountClient{base: newBase(AccountTokenHeader, accountToken, opts)}
}

func newBase(authHeader, token string, opts []Option) *base {
	cfg := defaultConfig()
	for _, o := range opts {
		o(cfg)
	}
	scheme := "https"
	if !cfg.secure {
		scheme = "http"
	}
	headers := map[string]string{
		authHeader:     token,
		"User-Agent":   cfg.userAgent,
		"Accept":       "application/json",
		"Content-Type": "application/json",
	}
	for k, v := range cfg.headers {
		headers[k] = v
	}
	return &base{
		baseURL: scheme + "://" + cfg.baseURL,
		headers: headers,
		http:    cfg.buildHTTPClient(),
	}
}

// Bool возвращает указатель на значение. Удобно для необязательных полей запросов.
func Bool(v bool) *bool { return &v }

// Int возвращает указатель на значение.
func Int(v int) *int { return &v }

// String возвращает указатель на значение.
func String(v string) *string { return &v }
