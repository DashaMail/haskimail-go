package haskimail_test

import (
	"context"
	"errors"
	"fmt"
	"log"

	"github.com/DashaMail/haskimail-go"
)

func Example() {
	client := haskimail.NewClient("ваш-серверный-токен")

	msg := &haskimail.Message{
		From:       "отправитель@домен.ru",
		To:         "получатель@домен.ru",
		Subject:    "Тема письма",
		HTMLBody:   "<h1>Привет!</h1><p>Это письмо отправлено через Haskimail.</p>",
		TrackOpens: haskimail.Bool(true),
		TrackLinks: haskimail.TrackLinksHTMLAndText,
	}
	msg.AddMetadata("campaign", "welcome")

	resp, err := client.DeliverMessage(context.Background(), msg)
	if errors.Is(err, haskimail.ErrInvalidMessage) {
		var apiErr *haskimail.Error
		errors.As(err, &apiErr)
		log.Fatalf("письмо отклонено: код %d, %s", apiErr.ErrorCode, apiErr.Message)
	} else if err != nil {
		log.Fatal(err)
	}
	fmt.Println("ID сообщения:", resp.MessageID)
}
