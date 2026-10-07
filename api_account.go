package haskimail

import (
	"context"
	"net/http"
)

// AccountClient — клиент API аккаунта Haskimail. Безопасен для конкурентного использования.
type AccountClient struct {
	*base
}

// ========== Серверы ==========

// GetServer возвращает сервер по ID.
func (c *AccountClient) GetServer(ctx context.Context, id int64) (*Server, error) {
	return do[Server](ctx, c.base, http.MethodGet, "/servers/"+itoa(id), nil, nil)
}

// CreateServer создаёт сервер.
func (c *AccountClient) CreateServer(ctx context.Context, s *Server) (*Server, error) {
	return do[Server](ctx, c.base, http.MethodPost, "/servers", nil, s)
}

// SetServer изменяет сервер.
func (c *AccountClient) SetServer(ctx context.Context, id int64, s *Server) (*Server, error) {
	return do[Server](ctx, c.base, http.MethodPut, "/servers/"+itoa(id), nil, s)
}

// GetServers возвращает список серверов.
func (c *AccountClient) GetServers(ctx context.Context, p *Params) (*Servers, error) {
	return do[Servers](ctx, c.base, http.MethodGet, "/servers", p, nil)
}

// DeleteServer удаляет сервер.
func (c *AccountClient) DeleteServer(ctx context.Context, id int64) (*RequestResponse, error) {
	return do[RequestResponse](ctx, c.base, http.MethodDelete, "/servers/"+itoa(id), nil, nil)
}

// ========== Домены ==========

// GetDomains возвращает список доменов.
func (c *AccountClient) GetDomains(ctx context.Context) (*Domains, error) {
	return do[Domains](ctx, c.base, http.MethodGet, "/domains", nil, nil)
}

// GetDomainDetails возвращает подробности домена.
func (c *AccountClient) GetDomainDetails(ctx context.Context, id int64) (*DomainDetails, error) {
	return do[DomainDetails](ctx, c.base, http.MethodGet, "/domains/"+itoa(id), nil, nil)
}

// CreateDomain добавляет домен.
func (c *AccountClient) CreateDomain(ctx context.Context, d *Domain) (*DomainDetails, error) {
	return do[DomainDetails](ctx, c.base, http.MethodPost, "/domains", nil, d)
}

// DeleteDomain удаляет домен.
func (c *AccountClient) DeleteDomain(ctx context.Context, id int64) (*RequestResponse, error) {
	return do[RequestResponse](ctx, c.base, http.MethodDelete, "/domains/"+itoa(id), nil, nil)
}

// VerifyDomainSPF проверяет SPF-запись домена.
func (c *AccountClient) VerifyDomainSPF(ctx context.Context, id int64) (*DomainDetails, error) {
	return do[DomainDetails](ctx, c.base, http.MethodPut, "/domains/"+itoa(id)+"/verifyspf", nil, nil)
}

// VerifyDomainDKIM проверяет DKIM-запись домена.
func (c *AccountClient) VerifyDomainDKIM(ctx context.Context, id int64) (*DomainDetails, error) {
	return do[DomainDetails](ctx, c.base, http.MethodPut, "/domains/"+itoa(id)+"/verifyDkim", nil, nil)
}

// ========== Пуш шаблонов ==========

// PushTemplates копирует шаблоны между серверами.
func (c *AccountClient) PushTemplates(ctx context.Context, r *TemplatesPushRequest) (*TemplatesPush, error) {
	return do[TemplatesPush](ctx, c.base, http.MethodPut, "/templates/push", nil, r)
}
