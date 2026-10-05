package haskimail

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"reflect"
	"strconv"
)

// base содержит общую для обоих клиентов логику HTTP.
type base struct {
	baseURL string
	headers map[string]string
	http    *http.Client
}

// BaseURL возвращает базовый URL клиента, например "https://api.haskimail.ru".
func (b *base) BaseURL() string { return b.baseURL }

// do выполняет запрос и декодирует ответ в T.
func do[T any](ctx context.Context, b *base, method, path string, params *Params, body any) (*T, error) {
	u := b.baseURL + path + params.query()

	var rdr io.Reader
	if body != nil {
		buf, err := json.Marshal(body)
		if err != nil {
			return nil, fmt.Errorf("haskimail: сериализация запроса: %w", err)
		}
		rdr = bytes.NewReader(buf)
	}

	req, err := http.NewRequestWithContext(ctx, method, u, rdr)
	if err != nil {
		return nil, fmt.Errorf("haskimail: создание запроса: %w", err)
	}
	for k, v := range b.headers {
		req.Header.Set(k, v)
	}

	resp, err := b.http.Do(req)
	if err != nil {
		return nil, fmt.Errorf("haskimail: %s %s: %w", method, path, err)
	}
	defer resp.Body.Close()

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("haskimail: чтение ответа: %w", err)
	}

	if resp.StatusCode < 200 || resp.StatusCode > 299 {
		return nil, newError(resp.StatusCode, data)
	}

	out := new(T)
	if err := decode(data, out); err != nil {
		return nil, fmt.Errorf("haskimail: разбор ответа: %w", err)
	}
	return out, nil
}

// decode разбирает JSON. Если API вернул
// массив, а ожидается одиночный объект, берётся первый элемент.
func decode(data []byte, out any) error {
	trimmed := bytes.TrimSpace(data)
	if len(trimmed) == 0 {
		return nil
	}
	if trimmed[0] == '[' {
		k := reflect.TypeOf(out).Elem().Kind()
		if k != reflect.Slice && k != reflect.Array && k != reflect.Interface {
			var items []json.RawMessage
			if err := json.Unmarshal(trimmed, &items); err != nil {
				return err
			}
			if len(items) == 0 {
				return nil
			}
			return json.Unmarshal(items[0], out)
		}
	}
	return json.Unmarshal(trimmed, out)
}

func itoa(id int64) string { return strconv.FormatInt(id, 10) }
