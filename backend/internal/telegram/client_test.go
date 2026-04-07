package telegram_test

import (
	"testing"

	"posiflora/backend/internal/telegram"
)

func TestMockClient_SendMessage_ReturnsNil(t *testing.T) {
	client := telegram.NewMockClient()
	if err := client.SendMessage("token", "123456", "test message"); err != nil {
		t.Errorf("mock client should not return error, got: %v", err)
	}
}

func TestMockClient_SendMessage_EmptyArgs(t *testing.T) {
	client := telegram.NewMockClient()
	if err := client.SendMessage("", "", ""); err != nil {
		t.Errorf("mock client should not return error for empty args, got: %v", err)
	}
}
