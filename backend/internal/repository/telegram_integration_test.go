package repository_test

import (
	"context"
	"testing"

	"posiflora/backend/internal/models"
	"posiflora/backend/internal/repository"
	"posiflora/backend/internal/test"
)

func TestIntegrationRepository_Upsert_Create(t *testing.T) {
	db := test.OpenDB(t)
	shop := test.SeedShop(t, db, "Repo Test Shop")

	repo := repository.NewIntegrationRepository(db)
	ti := &models.TelegramIntegration{
		ShopID: shop.ID, BotToken: "tok-1", ChatID: "chat-1", Enabled: true,
	}

	if err := repo.Upsert(context.Background(), ti); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if ti.ID == 0 {
		t.Error("expected non-zero ID after create")
	}

	var count int64
	db.Model(&models.TelegramIntegration{}).Where("shop_id = ?", shop.ID).Count(&count)
	if count != 1 {
		t.Errorf("expected 1 integration, got %d", count)
	}
}

func TestIntegrationRepository_Upsert_Update(t *testing.T) {
	db := test.OpenDB(t)
	shop := test.SeedShop(t, db, "Repo Test Shop")

	repo := repository.NewIntegrationRepository(db)
	ctx := context.Background()

	first := &models.TelegramIntegration{ShopID: shop.ID, BotToken: "old", ChatID: "old-chat", Enabled: true}
	if err := repo.Upsert(ctx, first); err != nil {
		t.Fatalf("first upsert error: %v", err)
	}

	second := &models.TelegramIntegration{ShopID: shop.ID, BotToken: "new", ChatID: "new-chat", Enabled: false}
	if err := repo.Upsert(ctx, second); err != nil {
		t.Fatalf("second upsert error: %v", err)
	}

	var ti models.TelegramIntegration
	if err := db.Where("shop_id = ?", shop.ID).First(&ti).Error; err != nil {
		t.Fatalf("integration not found: %v", err)
	}
	if ti.BotToken != "new" {
		t.Errorf("expected BotToken new, got %q", ti.BotToken)
	}
	if ti.ChatID != "new-chat" {
		t.Errorf("expected ChatID new-chat, got %q", ti.ChatID)
	}
	if ti.Enabled {
		t.Error("expected Enabled=false after update")
	}

	var count int64
	db.Model(&models.TelegramIntegration{}).Where("shop_id = ?", shop.ID).Count(&count)
	if count != 1 {
		t.Errorf("expected exactly 1 row, got %d", count)
	}
}

func TestIntegrationRepository_Upsert_EnabledFalsePersisted(t *testing.T) {
	db := test.OpenDB(t)
	shop := test.SeedShop(t, db, "Repo Test Shop")

	repo := repository.NewIntegrationRepository(db)
	ti := &models.TelegramIntegration{ShopID: shop.ID, BotToken: "tok", ChatID: "ch", Enabled: false}
	if err := repo.Upsert(context.Background(), ti); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	var persisted models.TelegramIntegration
	if err := db.Where("shop_id = ?", shop.ID).First(&persisted).Error; err != nil {
		t.Fatalf("not found: %v", err)
	}
	if persisted.Enabled {
		t.Error("expected Enabled=false in DB, got true")
	}
}

func TestIntegrationRepository_FindByShopID_Found(t *testing.T) {
	db := test.OpenDB(t)
	shop := test.SeedShop(t, db, "Repo Test Shop")

	repo := repository.NewIntegrationRepository(db)
	if err := repo.Upsert(context.Background(), &models.TelegramIntegration{
		ShopID: shop.ID, BotToken: "t", ChatID: "c", Enabled: true,
	}); err != nil {
		t.Fatalf("upsert: %v", err)
	}

	result, err := repo.FindByShopID(context.Background(), shop.ID)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result == nil {
		t.Fatal("expected integration, got nil")
		return
	}
	if result.ShopID != shop.ID {
		t.Errorf("expected shopID %d, got %d", shop.ID, result.ShopID)
	}
}

func TestIntegrationRepository_FindByShopID_NotFound(t *testing.T) {
	db := test.OpenDB(t)
	repo := repository.NewIntegrationRepository(db)

	result, err := repo.FindByShopID(context.Background(), 99999)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result != nil {
		t.Error("expected nil, got integration")
	}
}
