package telegram

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"time"
)

type Client interface {
	SendMessage(token, chatID, text string) error
}

type HTTPClient struct {
	client *http.Client
}

func NewHTTPClient() *HTTPClient {
	return &HTTPClient{
		client: &http.Client{Timeout: 10 * time.Second},
	}
}

type sendMessagePayload struct {
	ChatID string `json:"chat_id"`
	Text   string `json:"text"`
}

type telegramResponse struct {
	OK          bool   `json:"ok"`
	Description string `json:"description"`
}

func (c *HTTPClient) SendMessage(token, chatID, text string) error {
	url := fmt.Sprintf("https://api.telegram.org/bot%s/sendMessage", token)

	body, err := json.Marshal(sendMessagePayload{ChatID: chatID, Text: text})
	if err != nil {
		return fmt.Errorf("marshal payload: %w", err)
	}

	resp, err := c.client.Post(url, "application/json", bytes.NewReader(body))
	if err != nil {
		return fmt.Errorf("telegram http error: %w", err)
	}
	defer func() { _ = resp.Body.Close() }()

	if resp.StatusCode >= 500 {
		return fmt.Errorf("telegram server error: status %d", resp.StatusCode)
	}

	raw, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("read response body: %w", err)
	}

	var tgResp telegramResponse
	if err := json.Unmarshal(raw, &tgResp); err != nil {
		return fmt.Errorf("parse telegram response: %w", err)
	}
	if !tgResp.OK {
		return fmt.Errorf("telegram API error: %s", tgResp.Description)
	}
	return nil
}

type MockClient struct{}

func NewMockClient() *MockClient {
	return &MockClient{}
}

func (m *MockClient) SendMessage(token, chatID, text string) error {
	log.Printf("[TELEGRAM MOCK] chat_id=%s | %s", chatID, text)
	return nil
}
