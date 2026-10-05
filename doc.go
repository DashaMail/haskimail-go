// Package haskimail — Go-клиент для API Haskimail (https://api.haskimail.ru/),
// сервиса транзакционной электронной почты.
//
// Пакет предоставляет два клиента:
//
//   - [Client] — серверный API (отправка писем, шаблоны, отказы, статистика,
//     вебхуки, каналы, стоп-списки). Создаётся через [NewClient] с серверным токеном.
//   - [AccountClient] — API аккаунта (серверы, домены, подписи отправителей,
//     пуш шаблонов). Создаётся через [NewAccountClient]
//     с токеном аккаунта.
//
// Пример:
//
//	client := haskimail.NewClient("ваш-серверный-токен")
//	resp, err := client.DeliverMessage(ctx, &haskimail.Message{
//		From:     "отправитель@домен.ru",
//		To:       "получатель@домен.ru",
//		Subject:  "Тема письма",
//		HTMLBody: "<h1>Привет!</h1>",
//	})
//
// Все методы принимают [context.Context] для отмены и таймаутов. Ошибки API
// возвращаются как *[Error] и проверяются через errors.Is / errors.As:
//
//	if errors.Is(err, haskimail.ErrInvalidAPIKey) { ... }
//
// Пакет не имеет внешних зависимостей — только стандартная библиотека.
package haskimail
