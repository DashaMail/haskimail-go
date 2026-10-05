package haskimail

import (
	"net/url"
	"strconv"
	"time"
)

// Params — параметры строки запроса (пагинация, фильтры).
// Нулевое значение и nil готовы к использованию.
//
//	p := haskimail.NewParams().Int("count", 50).Int("offset", 0).String("tag", "welcome")
type Params struct {
	v url.Values
}

// NewParams создаёт пустой набор параметров.
func NewParams() *Params { return &Params{v: url.Values{}} }

func (p *Params) set(k, v string) *Params {
	if p.v == nil {
		p.v = url.Values{}
	}
	p.v.Set(k, v)
	return p
}

// String добавляет строковый параметр.
func (p *Params) String(name, value string) *Params { return p.set(name, value) }

// Int добавляет целочисленный параметр.
func (p *Params) Int(name string, value int) *Params { return p.set(name, strconv.Itoa(value)) }

// Bool добавляет логический параметр.
func (p *Params) Bool(name string, value bool) *Params {
	return p.set(name, strconv.FormatBool(value))
}

// Date добавляет дату в формате yyyy-MM-dd.
func (p *Params) Date(name string, value time.Time) *Params {
	return p.set(name, value.Format("2006-01-02"))
}

// Encode возвращает закодированную строку запроса без "?".
func (p *Params) Encode() string {
	if p == nil || len(p.v) == 0 {
		return ""
	}
	return p.v.Encode()
}

func (p *Params) query() string {
	if s := p.Encode(); s != "" {
		return "?" + s
	}
	return ""
}
