# Haskimail Go

Go-библиотека для работы с [API Haskimail](https://api.haskimail.ru/) — сервиса транзакционной электронной почты.

Без внешних зависимостей, только стандартная библиотека. Требуется Go 1.21+.

## Установка

```bash
go get github.com/DashaMail/haskimail-go
```

## Быстрый старт

```go
client := haskimail.NewClient("ваш-серверный-токен")

resp, err := client.DeliverMessage(ctx, &haskimail.Message{
    From:       "отправитель@домен.ru",
    To:         "получатель@домен.ru",
    Subject:    "Тема письма",
    HTMLBody:   "<h1>Привет!</h1>",
    TrackOpens: haskimail.Bool(true),
    TrackLinks: haskimail.TrackLinksHTMLAndText,
})
if err != nil {
    log.Fatal(err)
}
fmt.Println("ID сообщения:", resp.MessageID)
```

### Отправка по шаблону

```go
client.DeliverMessageWithTemplate(ctx, &haskimail.TemplatedMessage{
    TemplateID:    12345, // или TemplateAlias: "welcome"
    From:          "от@домен.ru",
    To:            "кому@домен.ru",
    // Замена буквальная: ключ — фрагмент текста шаблона вместе с разделителями
    TemplateModel: map[string]string{"{{name}}": "Иван", "{{product}}": "Haskimail"},
})
```

### Вложения

```go
att, err := haskimail.AttachmentFromFile("/путь/к/отчёт.pdf", "")
msg.AddAttachment(att)

// встроенное изображение
msg.AddAttachment(haskimail.NewAttachment("logo.png", pngBytes, "image/png", "cid:logo"))
```

### Получатели с именами

```go
msg.To = haskimail.FormatRecipients(
    haskimail.Recipient{Name: "Иван", Email: "ivan@домен.ru"},
    haskimail.Recipient{Name: "Анна", Email: "anna@домен.ru"},
)
```

### Параметры списков

```go
bounces, err := client.GetBounces(ctx, haskimail.NewParams().
    Int("count", 50).
    Int("offset", 0).
    String("type", "HardBounce").
    Date("fromdate", time.Now().AddDate(0, 0, -7)))
```

### API аккаунта

```go
account := haskimail.NewAccountClient("токен-аккаунта")
servers, err := account.GetServers(ctx, haskimail.NewParams().Int("count", 10).Int("offset", 0))
```

## Настройка клиента

```go
client := haskimail.NewClient(token,
    haskimail.WithBaseURL("api.example.ru"),       // свой хост
    haskimail.WithInsecure(),                       // HTTP вместо HTTPS
    haskimail.WithConnectTimeout(10*time.Second),   // по умолчанию 60 с
    haskimail.WithReadTimeout(30*time.Second),      // по умолчанию 60 с
    haskimail.WithHTTPClient(myHTTPClient),         // свой *http.Client
)
```

Клиенты потокобезопасны: создайте один экземпляр и используйте его из всех горутин. Каждый метод принимает `context.Context` для отмены и дедлайнов.

## Обработка ошибок

Ошибки API имеют тип `*haskimail.Error` и проверяются через `errors.Is` / `errors.As`:

```go
_, err := client.DeliverMessage(ctx, msg)

var apiErr *haskimail.Error
switch {
case errors.Is(err, haskimail.ErrInvalidAPIKey):   // 401
case errors.Is(err, haskimail.ErrInvalidMessage):  // 422
    errors.As(err, &apiErr)
    log.Printf("код %d: %s", apiErr.ErrorCode, apiErr.Message)
case errors.Is(err, haskimail.ErrTimeout):         // 408
case errors.Is(err, haskimail.ErrInternalServer):  // 500
case errors.Is(err, haskimail.ErrUnknown):         // прочие HTTP-коды
case err != nil:                                    // сетевые ошибки, отмена контекста
}
```

## Особенности API

- Серверный токен (`NewClient`) работает только с данными своего сервера: письмами, статистикой, шаблонами, каналами, вебхуками. Серверы, домены и перенос шаблонов — только через `NewAccountClient` с токеном аккаунта.
- Если в письме не указан `MessageStream`, оно уходит через транзакционный канал сервера по умолчанию.
- Перегрузок в Go нет, поэтому операции по алиасу вынесены в отдельные методы: `GetTemplateByAlias`, `SetTemplateByAlias`, `DeleteTemplateByAlias`, а открытия и клики конкретного сообщения — в `GetMessageOpensByID` и `GetMessageClicksByID`.
- Необязательные логические поля в запросах — указатели (`*bool`), чтобы можно было явно передать `false`. Используйте `haskimail.Bool(…)`, `haskimail.Int(…)`.
- Даты в ответах имеют тип `*haskimail.Time`, который встраивает `time.Time` и понимает все форматы дат API. Даты без часового пояса трактуются в `haskimail.TimeLocation` (по умолчанию локальный пояс).
- Значения параметров запроса URL-кодируются, а ID и алиасы в путях экранируются автоматически.
- Тип вложений из файла определяется по расширению, а если не удалось — по содержимому.

## Тесты

Юнит-тесты (без сети, на `httptest`): проверяют все 50 методов (HTTP-метод, путь, query, тело, токен), сериализацию, ошибки, форматы дат, вложения и отмену по контексту.

```bash
go test -race ./...
```

Интеграционные тесты против живого API:

```bash
export HASKIMAIL_SERVER_TOKEN=...       # обязательно
export HASKIMAIL_ACCOUNT_TOKEN=...      # для тестов API аккаунта
export HASKIMAIL_SENDER_EMAIL=...       # подтверждённый отправитель (для отправки писем)
export HASKIMAIL_RECIPIENT_EMAIL=...    # получатель тестовых писем
go test -tags integration -v -run Integration ./...
```

Тесты только читают данные, кроме отправки тестовых писем (если заданы отправитель и получатель) и создания временных шаблона и вебхука, которые сразу удаляются. Каждый ответ API дополнительно сверяется с Go-моделями: если API вернул поле, которого нет в модели, в выводе появится строка `НЕИЗВЕСТНЫЕ ПОЛЯ` с телом ответа.
