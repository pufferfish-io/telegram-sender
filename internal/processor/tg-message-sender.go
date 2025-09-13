package processor

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"tg-sender/internal/contract"
	"tg-sender/internal/logger"
)

type TgMessageSender struct {
	token      string
	apiBase    string
	httpClient *http.Client
	logger     logger.Logger
}

type Option struct {
	Token      string
	ApiBase    string
	HttpClient *http.Client
	Logger     logger.Logger
}

func NewTgMessageSender(opt Option) *TgMessageSender {
	return &TgMessageSender{
		token:      opt.Token,
		apiBase:    opt.ApiBase,
		httpClient: opt.HttpClient,
		logger:     opt.Logger,
	}
}

func (t *TgMessageSender) Handle(ctx context.Context, raw []byte) error {
	var requestMessage contract.SendMessageRequest
	if err := json.Unmarshal(raw, &requestMessage); err != nil {
		return err
	}

	url := fmt.Sprintf("%s/sendMessage", t.apiBase)

	payload, err := json.Marshal(requestMessage)
	if err != nil {
		return fmt.Errorf("failed to marshal request: %w", err)
	}

	resp, err := t.httpClient.Post(url, "application/json", bytes.NewReader(payload))
	if err != nil {
		return fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 300 {
		return fmt.Errorf("telegram API responded with status: %s", resp.Status)
	}

	t.logger.Info("✅ Message sent to chat %d", requestMessage.ChatID)

	return nil
}
