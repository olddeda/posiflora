package repository_test

import (
	"context"
	"testing"
	"time"

	"posiflora/backend/internal/models"
	"posiflora/backend/internal/repository"
	"posiflora/backend/internal/test"
)

func TestSendLogRepository_Create_And_FindByShopAndOrder(t *testing.T) {
	db := test.OpenDB(t)
	shop := test.SeedShop(t, db, "Repo Test Shop")
	order := &models.Order{ShopID: shop.ID, Number: "L-001", Total: 100, CustomerName: "X"}
	if err := db.Create(order).Error; err != nil {
		t.Fatalf("create order: %v", err)
	}

	repo := repository.NewTelegramSendLogRepository(db)
	ctx := context.Background()

	existing, err := repo.FindByShopAndOrder(ctx, shop.ID, order.ID)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if existing != nil {
		t.Error("expected nil before create")
	}

	entry := &models.TelegramSendLog{
		ShopID: shop.ID, OrderID: order.ID, Message: "test msg", Status: models.StatusSent,
	}
	if err := repo.Create(ctx, entry); err != nil {
		t.Fatalf("create log: %v", err)
	}

	found, err := repo.FindByShopAndOrder(ctx, shop.ID, order.ID)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if found == nil {
		t.Fatal("expected log entry, got nil")
		return
	}
	if found.Status != models.StatusSent {
		t.Errorf("expected SENT, got %q", found.Status)
	}
}

func TestSendLogRepository_Create_OnConflict_DoNothing(t *testing.T) {
	db := test.OpenDB(t)
	shop := test.SeedShop(t, db, "Repo Test Shop")
	order := &models.Order{ShopID: shop.ID, Number: "L-002", Total: 200, CustomerName: "Y"}
	if err := db.Create(order).Error; err != nil {
		t.Fatalf("create order: %v", err)
	}

	repo := repository.NewTelegramSendLogRepository(db)
	ctx := context.Background()

	if err := repo.Create(ctx, &models.TelegramSendLog{
		ShopID: shop.ID, OrderID: order.ID, Message: "first", Status: models.StatusSent,
	}); err != nil {
		t.Fatalf("first create: %v", err)
	}
	if err := repo.Create(ctx, &models.TelegramSendLog{
		ShopID: shop.ID, OrderID: order.ID, Message: "dup", Status: models.StatusFailed,
	}); err != nil {
		t.Fatalf("duplicate create should not error: %v", err)
	}

	var count int64
	db.Model(&models.TelegramSendLog{}).Where("shop_id = ? AND order_id = ?", shop.ID, order.ID).Count(&count)
	if count != 1 {
		t.Errorf("expected 1 row after duplicate insert, got %d", count)
	}
}

func TestSendLogRepository_GetCounts_WithinWindow(t *testing.T) {
	db := test.OpenDB(t)
	shop := test.SeedShop(t, db, "Repo Test Shop")

	repo := repository.NewTelegramSendLogRepository(db)
	ctx := context.Background()

	orders := make([]*models.Order, 3)
	for i := range orders {
		o := &models.Order{ShopID: shop.ID, Number: "G-00" + string(rune('1'+i)), Total: 100, CustomerName: "Z"}
		if err := db.Create(o).Error; err != nil {
			t.Fatalf("create order: %v", err)
		}
		orders[i] = o
	}

	statuses := []models.SendStatus{models.StatusSent, models.StatusSent, models.StatusFailed}
	for i, o := range orders {
		if err := repo.Create(ctx, &models.TelegramSendLog{
			ShopID: shop.ID, OrderID: o.ID, Message: "m", Status: statuses[i],
		}); err != nil {
			t.Fatalf("create log: %v", err)
		}
	}

	since := time.Now().AddDate(0, 0, -7)
	sentCount, failedCount, lastSentAt, err := repo.GetCounts(ctx, shop.ID, since)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if sentCount != 2 {
		t.Errorf("expected sentCount=2, got %d", sentCount)
	}
	if failedCount != 1 {
		t.Errorf("expected failedCount=1, got %d", failedCount)
	}
	if lastSentAt == nil {
		t.Error("expected non-nil lastSentAt")
	}
}

func TestSendLogRepository_GetCounts_ExcludesOldRecords(t *testing.T) {
	db := test.OpenDB(t)
	shop := test.SeedShop(t, db, "Repo Test Shop")

	order := &models.Order{ShopID: shop.ID, Number: "G-OLD", Total: 100, CustomerName: "Z"}
	if err := db.Create(order).Error; err != nil {
		t.Fatalf("create order: %v", err)
	}

	oldLog := &models.TelegramSendLog{
		ShopID: shop.ID, OrderID: order.ID, Message: "old", Status: models.StatusSent,
	}
	if err := db.Create(oldLog).Error; err != nil {
		t.Fatalf("create old log: %v", err)
	}
	if err := db.Model(oldLog).Update("sent_at", time.Now().AddDate(0, 0, -10)).Error; err != nil {
		t.Fatalf("update sent_at: %v", err)
	}

	since := time.Now().AddDate(0, 0, -7)
	sentCount, failedCount, _, err := repository.NewTelegramSendLogRepository(db).GetCounts(context.Background(), shop.ID, since)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if sentCount != 0 {
		t.Errorf("expected sentCount=0 for old records, got %d", sentCount)
	}
	if failedCount != 0 {
		t.Errorf("expected failedCount=0 for old records, got %d", failedCount)
	}
}

func TestSendLogRepository_GetCounts_Empty(t *testing.T) {
	db := test.OpenDB(t)
	shop := test.SeedShop(t, db, "Repo Test Shop")

	since := time.Now().AddDate(0, 0, -7)
	sentCount, failedCount, lastSentAt, err := repository.NewTelegramSendLogRepository(db).GetCounts(context.Background(), shop.ID, since)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if sentCount != 0 || failedCount != 0 {
		t.Errorf("expected zero counts, got sent=%d failed=%d", sentCount, failedCount)
	}
	if lastSentAt != nil {
		t.Error("expected nil lastSentAt for empty log")
	}
}
