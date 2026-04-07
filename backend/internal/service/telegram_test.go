package service_test

import (
	"context"
	"testing"

	"posiflora/backend/internal/dto/requests"
	"posiflora/backend/internal/models"
	"posiflora/backend/internal/repository"
	"posiflora/backend/internal/service"
	"posiflora/backend/internal/test"
)

func TestTelegramConnect_InitialSetup(t *testing.T) {
	db := test.OpenDB(t)
	shop := &models.Shop{Name: "TG Shop"}
	if err := db.Create(shop).Error; err != nil {
		t.Fatalf("create shop: %v", err)
	}

	repos := repository.New(db)
	svc := service.NewTelegramService(repos)

	resp, err := svc.Connect(context.Background(), shop.ID, requests.ConnectTelegram{
		BotToken: "123:ABC",
		ChatID:   "987654",
		Enabled:  true,
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if resp.ChatID != "987654" {
		t.Errorf("expected chatId 987654, got %q", resp.ChatID)
	}
	if !resp.Enabled {
		t.Error("expected enabled=true")
	}
	if resp.ShopID != shop.ID {
		t.Errorf("expected shopID %d, got %d", shop.ID, resp.ShopID)
	}
}

func TestTelegramConnect_EnabledFalse_Persisted(t *testing.T) {
	db := test.OpenDB(t)
	shop := &models.Shop{Name: "TG Shop Disabled"}
	if err := db.Create(shop).Error; err != nil {
		t.Fatalf("create shop: %v", err)
	}

	repos := repository.New(db)
	svc := service.NewTelegramService(repos)

	resp, err := svc.Connect(context.Background(), shop.ID, requests.ConnectTelegram{
		BotToken: "123:ABC",
		ChatID:   "111222",
		Enabled:  false,
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if resp.Enabled {
		t.Error("expected enabled=false, got true")
	}

	var ti models.TelegramIntegration
	if err := db.Where("shop_id = ?", shop.ID).First(&ti).Error; err != nil {
		t.Fatalf("integration not found: %v", err)
	}
	if ti.Enabled {
		t.Error("expected enabled=false in DB, got true")
	}
}

func TestTelegramConnect_UpdateExisting(t *testing.T) {
	db := test.OpenDB(t)
	shop := &models.Shop{Name: "TG Update Shop"}
	if err := db.Create(shop).Error; err != nil {
		t.Fatalf("create shop: %v", err)
	}

	repos := repository.New(db)
	svc := service.NewTelegramService(repos)
	ctx := context.Background()

	_, err := svc.Connect(ctx, shop.ID, requests.ConnectTelegram{
		BotToken: "old-token", ChatID: "old-chat", Enabled: true,
	})
	if err != nil {
		t.Fatalf("first connect error: %v", err)
	}

	resp, err := svc.Connect(ctx, shop.ID, requests.ConnectTelegram{
		BotToken: "new-token", ChatID: "new-chat", Enabled: false,
	})
	if err != nil {
		t.Fatalf("second connect error: %v", err)
	}
	if resp.ChatID != "new-chat" {
		t.Errorf("expected chatId new-chat, got %q", resp.ChatID)
	}
	if resp.Enabled {
		t.Error("expected enabled=false after update")
	}

	var count int64
	db.Model(&models.TelegramIntegration{}).Where("shop_id = ?", shop.ID).Count(&count)
	if count != 1 {
		t.Errorf("expected exactly 1 integration row, got %d", count)
	}
}

func TestTelegramConnect_UpdateWithoutToken_UsesExisting(t *testing.T) {
	db := test.OpenDB(t)
	shop := &models.Shop{Name: "TG Token Reuse"}
	if err := db.Create(shop).Error; err != nil {
		t.Fatalf("create shop: %v", err)
	}

	repos := repository.New(db)
	svc := service.NewTelegramService(repos)
	ctx := context.Background()

	_, err := svc.Connect(ctx, shop.ID, requests.ConnectTelegram{
		BotToken: "original-token", ChatID: "chat-1", Enabled: true,
	})
	if err != nil {
		t.Fatalf("initial connect error: %v", err)
	}

	_, err = svc.Connect(ctx, shop.ID, requests.ConnectTelegram{
		BotToken: "", ChatID: "chat-2", Enabled: false,
	})
	if err != nil {
		t.Fatalf("update without token error: %v", err)
	}

	var ti models.TelegramIntegration
	if err := db.Where("shop_id = ?", shop.ID).First(&ti).Error; err != nil {
		t.Fatalf("integration not found: %v", err)
	}
	if ti.BotToken != "original-token" {
		t.Errorf("expected original-token preserved, got %q", ti.BotToken)
	}
	if ti.ChatID != "chat-2" {
		t.Errorf("expected chat-2, got %q", ti.ChatID)
	}
}

func TestTelegramConnect_InitialSetup_EmptyToken_ReturnsError(t *testing.T) {
	db := test.OpenDB(t)
	shop := &models.Shop{Name: "TG No Token"}
	if err := db.Create(shop).Error; err != nil {
		t.Fatalf("create shop: %v", err)
	}

	repos := repository.New(db)
	svc := service.NewTelegramService(repos)

	_, err := svc.Connect(context.Background(), shop.ID, requests.ConnectTelegram{
		BotToken: "", ChatID: "123", Enabled: true,
	})
	if err == nil {
		t.Fatal("expected error when botToken is empty for initial setup")
	}
}

func TestTelegramConnect_ShopNotFound(t *testing.T) {
	db := test.OpenDB(t)
	repos := repository.New(db)
	svc := service.NewTelegramService(repos)

	_, err := svc.Connect(context.Background(), 99999, requests.ConnectTelegram{
		BotToken: "tok", ChatID: "123", Enabled: true,
	})
	if err == nil {
		t.Fatal("expected error for non-existent shop")
	}
}

func TestTelegramGetStatus_NoIntegration(t *testing.T) {
	db := test.OpenDB(t)
	shop := &models.Shop{Name: "Status Shop Empty"}
	if err := db.Create(shop).Error; err != nil {
		t.Fatalf("create shop: %v", err)
	}

	repos := repository.New(db)
	svc := service.NewTelegramService(repos)

	status, err := svc.GetStatus(context.Background(), shop.ID)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if status.ChatID != "" {
		t.Errorf("expected empty chatId, got %q", status.ChatID)
	}
	if status.SentCount7d != 0 || status.FailedCount7d != 0 {
		t.Error("expected zero counts for shop with no integration")
	}
}

func TestTelegramGetStatus_WithIntegration(t *testing.T) {
	db := test.OpenDB(t)
	shop := test.SeedShopWithIntegration(t, db)

	repos := repository.New(db)
	svc := service.NewTelegramService(repos)

	status, err := svc.GetStatus(context.Background(), shop.ID)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if status.ChatID != "123456" {
		t.Errorf("expected chatId 123456, got %q", status.ChatID)
	}
	if !status.Enabled {
		t.Error("expected enabled=true")
	}
}

func TestTelegramGetStatus_ShopNotFound(t *testing.T) {
	db := test.OpenDB(t)
	repos := repository.New(db)
	svc := service.NewTelegramService(repos)

	_, err := svc.GetStatus(context.Background(), 99999)
	if err == nil {
		t.Fatal("expected error for non-existent shop")
	}
}
