package handler_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"posiflora/backend/internal/test"
)

func TestHealth(t *testing.T) {
	db := test.OpenDB(t)
	router := test.NewTestRouter(t, db, "../../locales")

	w := test.GetRequest(router, "/health")
	if w.Code != http.StatusOK {
		t.Errorf("expected 200, got %d", w.Code)
	}
}

func TestCORS_Preflight(t *testing.T) {
	db := test.OpenDB(t)
	router := test.NewTestRouter(t, db, "../../locales")

	req := httptest.NewRequest(http.MethodOptions, "/shops/1/orders", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusNoContent {
		t.Errorf("expected 204, got %d", w.Code)
	}
	if w.Header().Get("Access-Control-Allow-Origin") == "" {
		t.Error("expected Access-Control-Allow-Origin header")
	}
}
