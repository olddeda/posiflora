package handler_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"posiflora/backend/internal/test"
)

func TestOrderCreate_Success(t *testing.T) {
	db := test.OpenDB(t)
	shop := test.SeedShopWithIntegration(t, db)
	router := test.NewTestRouter(t, db, "../../locales")

	w := test.PostJSON(router, "/shops/"+test.Itoa(shop.ID)+"/orders", map[string]any{
		"number": "H-001", "total": 1000, "customerName": "Тест",
	})
	if w.Code != http.StatusCreated {
		t.Errorf("expected 201, got %d: %s", w.Code, w.Body.String())
	}

	var resp map[string]any
	if err := json.Unmarshal(w.Body.Bytes(), &resp); err != nil {
		t.Fatalf("unmarshal response: %v", err)
	}
	if resp["notifyStatus"] != "sent" {
		t.Errorf("expected notifyStatus=sent, got %v", resp["notifyStatus"])
	}
}

func TestOrderCreate_ShopNotFound(t *testing.T) {
	db := test.OpenDB(t)
	router := test.NewTestRouter(t, db, "../../locales")

	w := test.PostJSON(router, "/shops/99999/orders", map[string]any{
		"number": "H-002", "total": 500, "customerName": "X",
	})
	if w.Code != http.StatusNotFound {
		t.Errorf("expected 404, got %d", w.Code)
	}
}

func TestOrderCreate_InvalidShopID(t *testing.T) {
	db := test.OpenDB(t)
	router := test.NewTestRouter(t, db, "../../locales")

	w := test.PostJSON(router, "/shops/abc/orders", map[string]any{
		"number": "H-003", "total": 100, "customerName": "X",
	})
	if w.Code != http.StatusBadRequest {
		t.Errorf("expected 400, got %d", w.Code)
	}
}

func TestOrderCreate_MissingFields(t *testing.T) {
	db := test.OpenDB(t)
	shop := test.SeedShop(t, db, "Handler Test Shop")
	router := test.NewTestRouter(t, db, "../../locales")

	w := test.PostJSON(router, "/shops/"+test.Itoa(shop.ID)+"/orders", map[string]any{
		"number": "H-004",
	})
	if w.Code != http.StatusBadRequest {
		t.Errorf("expected 400, got %d", w.Code)
	}
}

func TestOrderCreate_InvalidJSON(t *testing.T) {
	db := test.OpenDB(t)
	shop := test.SeedShop(t, db, "Handler Test Shop")
	router := test.NewTestRouter(t, db, "../../locales")

	req := httptest.NewRequest(http.MethodPost, "/shops/"+test.Itoa(shop.ID)+"/orders", bytes.NewBufferString("{bad json}"))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("expected 400, got %d", w.Code)
	}
}
