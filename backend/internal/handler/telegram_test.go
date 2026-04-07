package handler_test

import (
	"encoding/json"
	"net/http"
	"testing"

	"posiflora/backend/internal/test"
)

func TestTelegramConnect_Success(t *testing.T) {
	db := test.OpenDB(t)
	shop := test.SeedShop(t, db, "Handler Test Shop")
	router := test.NewTestRouter(t, db, "../../locales")

	w := test.PostJSON(router, "/shops/"+test.Itoa(shop.ID)+"/telegram/connect", map[string]any{
		"botToken": "123:ABC", "chatId": "999888", "enabled": true,
	})
	if w.Code != http.StatusOK {
		t.Errorf("expected 200, got %d: %s", w.Code, w.Body.String())
	}

	var resp map[string]any
	if err := json.Unmarshal(w.Body.Bytes(), &resp); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	if resp["chatId"] != "999888" {
		t.Errorf("expected chatId 999888, got %v", resp["chatId"])
	}
}

func TestTelegramConnect_MissingChatID(t *testing.T) {
	db := test.OpenDB(t)
	shop := test.SeedShop(t, db, "Handler Test Shop")
	router := test.NewTestRouter(t, db, "../../locales")

	w := test.PostJSON(router, "/shops/"+test.Itoa(shop.ID)+"/telegram/connect", map[string]any{
		"botToken": "123:ABC", "enabled": true,
	})
	if w.Code != http.StatusBadRequest {
		t.Errorf("expected 400, got %d", w.Code)
	}
}

func TestTelegramConnect_ShopNotFound(t *testing.T) {
	db := test.OpenDB(t)
	router := test.NewTestRouter(t, db, "../../locales")

	w := test.PostJSON(router, "/shops/99999/telegram/connect", map[string]any{
		"botToken": "tok", "chatId": "123", "enabled": true,
	})
	if w.Code != http.StatusNotFound {
		t.Errorf("expected 404, got %d", w.Code)
	}
}

func TestTelegramGetStatus_NoIntegration(t *testing.T) {
	db := test.OpenDB(t)
	shop := test.SeedShop(t, db, "Handler Test Shop")
	router := test.NewTestRouter(t, db, "../../locales")

	w := test.GetRequest(router, "/shops/"+test.Itoa(shop.ID)+"/telegram/status")
	if w.Code != http.StatusOK {
		t.Errorf("expected 200, got %d", w.Code)
	}
}

func TestTelegramGetStatus_WithIntegration(t *testing.T) {
	db := test.OpenDB(t)
	shop := test.SeedShopWithIntegration(t, db)
	router := test.NewTestRouter(t, db, "../../locales")

	w := test.GetRequest(router, "/shops/"+test.Itoa(shop.ID)+"/telegram/status")
	if w.Code != http.StatusOK {
		t.Errorf("expected 200, got %d", w.Code)
	}

	var resp map[string]any
	if err := json.Unmarshal(w.Body.Bytes(), &resp); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	if resp["chatId"] != "123456" {
		t.Errorf("expected chatId 123456, got %v", resp["chatId"])
	}
	if resp["enabled"] != true {
		t.Errorf("expected enabled=true, got %v", resp["enabled"])
	}
}

func TestTelegramGetStatus_ShopNotFound(t *testing.T) {
	db := test.OpenDB(t)
	router := test.NewTestRouter(t, db, "../../locales")

	w := test.GetRequest(router, "/shops/99999/telegram/status")
	if w.Code != http.StatusNotFound {
		t.Errorf("expected 404, got %d", w.Code)
	}
}

func TestTelegramGetStatus_InvalidShopID(t *testing.T) {
	db := test.OpenDB(t)
	router := test.NewTestRouter(t, db, "../../locales")

	w := test.GetRequest(router, "/shops/abc/telegram/status")
	if w.Code != http.StatusBadRequest {
		t.Errorf("expected 400, got %d", w.Code)
	}
}
