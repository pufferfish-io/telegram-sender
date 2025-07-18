package service

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"tg-sender/internal/model"
)

type TgService struct {
	token      string
	apiBase    string
	httpClient *http.Client
}

func NewTelegramSenderService(token string) *TgService {
	return &TgService{
		token:      token,
		apiBase:    fmt.Sprintf("https://api.telegram.org/bot%s", token),
		httpClient: &http.Client{},
	}
}

func (s *TgService) SendMessage(req model.SendMessageRequest) error {
	url := fmt.Sprintf("%s/sendMessage", s.apiBase)

	payload, err := json.Marshal(req)
	if err != nil {
		return fmt.Errorf("failed to marshal request: %w", err)
	}

	resp, err := s.httpClient.Post(url, "application/json", bytes.NewReader(payload))
	if err != nil {
		return fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 300 {
		return fmt.Errorf("telegram API responded with status: %s", resp.Status)
	}

	log.Printf("✅ Message sent to chat %d", req.ChatID)
	return nil
}
