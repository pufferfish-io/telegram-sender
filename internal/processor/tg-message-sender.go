package processor

import (
	"context"
	"encoding/json"

	"tg-sender/internal/model"
	"tg-sender/internal/service"
)

type TgMessageSender struct {
	tgService *service.TgService
}

func NewTgMessageSender(kafkaTopic string, telegramService *service.TgService) *TgMessageSender {
	return &TgMessageSender{
		tgService: telegramService,
	}
}

func (t *TgMessageSender) Handle(ctx context.Context, raw []byte) error {
	var requestMessage model.SendMessageRequest
	if err := json.Unmarshal(raw, &requestMessage); err != nil {
		return err
	}

	err := t.tgService.SendMessage(requestMessage)
	if err != nil {
		return err
	}

	return nil
}
