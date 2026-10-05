package haskimail

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
)

// Сентинел-ошибки для проверки через errors.Is.
var (
	// ErrInvalidAPIKey — HTTP 401, неверный или отсутствующий токен.
	ErrInvalidAPIKey = errors.New("haskimail: неверный API-токен")
	// ErrTimeout — HTTP 408.
	ErrTimeout = errors.New("haskimail: таймаут запроса")
	// ErrInvalidMessage — HTTP 422, ошибка валидации запроса.
	ErrInvalidMessage = errors.New("haskimail: некорректный запрос")
	// ErrInternalServer — HTTP 500.
	ErrInternalServer = errors.New("haskimail: внутренняя ошибка сервера")
	// ErrUnknown — любой другой неуспешный HTTP-статус.
	ErrUnknown = errors.New("haskimail: неизвестная ошибка")
)

// Error — ошибка, возвращённая API Haskimail.
type Error struct {
	// StatusCode — HTTP-код ответа.
	StatusCode int
	// ErrorCode — код ошибки API из тела ответа (0, если тело не разобрано).
	ErrorCode int
	// Message — описание ошибки из тела ответа.
	Message string
	// Body — сырое тело ответа.
	Body []byte
}

func (e *Error) Error() string {
	return fmt.Sprintf("haskimail: HTTP %d, код %d: %s", e.StatusCode, e.ErrorCode, e.Message)
}

// Unwrap возвращает сентинел-ошибку, соответствующую HTTP-коду,
// поэтому работает errors.Is(err, haskimail.ErrInvalidMessage).
func (e *Error) Unwrap() error {
	switch e.StatusCode {
	case http.StatusUnauthorized:
		return ErrInvalidAPIKey
	case http.StatusRequestTimeout:
		return ErrTimeout
	case http.StatusUnprocessableEntity:
		return ErrInvalidMessage
	case http.StatusInternalServerError:
		return ErrInternalServer
	default:
		return ErrUnknown
	}
}

func newError(status int, body []byte) *Error {
	e := &Error{StatusCode: status, Message: "Неизвестная ошибка", Body: body}
	var apiErr struct {
		ErrorCode int
		Message   string
	}
	if json.Unmarshal(body, &apiErr) == nil {
		e.ErrorCode = apiErr.ErrorCode
		if apiErr.Message != "" {
			e.Message = apiErr.Message
		}
	}
	return e
}
