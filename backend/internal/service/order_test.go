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

func TestCreateOrder_SendsNotificationAndLogsSuccess(t *testing.T) {
	db := test.OpenDB(t)
	shop := test.SeedShopWithIntegration(t, db)

	tg := &test.SuccessClient{}
	repos := repository.New(db)
	svc := service.NewOrderService(repos, tg, test.LoadTranslator(t, "../../locales", "ru"))

	resp, err := svc.Create(context.Background(), shop.ID, requests.CreateOrder{
		Number: "T-001", Total: 1500, CustomerName: "Иван",
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(tg.Calls) != 1 {
		t.Errorf("expected 1 telegram call, got %d", len(tg.Calls))
	}
	if resp.NotifyStatus != "sent" {
		t.Errorf("expected notifyStatus=sent, got %q", resp.NotifyStatus)
	}

	var logEntry models.TelegramSendLog
	if err := db.Where("shop_id = ? AND order_id = ?", shop.ID, resp.Order.ID).First(&logEntry).Error; err != nil {
		t.Fatalf("log not found: %v", err)
	}
	if logEntry.Status != models.StatusSent {
		t.Errorf("expected log status SENT, got %q", logEntry.Status)
	}
}

func TestCreateOrder_Idempotent_NoDuplicateLogs(t *testing.T) {
	db := test.OpenDB(t)
	shop := test.SeedShopWithIntegration(t, db)

	tg := &test.SuccessClient{}
	repos := repository.New(db)
	svc := service.NewOrderService(repos, tg, test.LoadTranslator(t, "../../locales", "ru"))

	ctx := context.Background()

	resp, err := svc.Create(ctx, shop.ID, requests.CreateOrder{
		Number: "T-002", Total: 999, CustomerName: "Мария",
	})
	if err != nil {
		t.Fatalf("first call error: %v", err)
	}

	var count int64
	db.Model(&models.TelegramSendLog{}).
		Where("shop_id = ? AND order_id = ?", shop.ID, resp.Order.ID).
		Count(&count)
	if count != 1 {
		t.Errorf("expected 1 log row after first call, got %d", count)
	}

	dup := &models.TelegramSendLog{
		ShopID: shop.ID, OrderID: resp.Order.ID, Message: "duplicate", Status: models.StatusSent,
	}
	if err := repos.SendLog.Create(ctx, dup); err != nil {
		t.Fatalf("duplicate insert: %v", err)
	}

	db.Model(&models.TelegramSendLog{}).
		Where("shop_id = ? AND order_id = ?", shop.ID, resp.Order.ID).
		Count(&count)
	if count != 1 {
		t.Errorf("expected still 1 log row after duplicate insert, got %d", count)
	}
	if len(tg.Calls) != 1 {
		t.Errorf("expected 1 total telegram call, got %d", len(tg.Calls))
	}
}

func TestCreateOrder_TelegramFailure_OrderStillCreated(t *testing.T) {
	db := test.OpenDB(t)
	shop := test.SeedShopWithIntegration(t, db)

	repos := repository.New(db)
	svc := service.NewOrderService(repos, &test.FailClient{}, test.LoadTranslator(t, "../../locales", "ru"))

	resp, err := svc.Create(context.Background(), shop.ID, requests.CreateOrder{
		Number: "T-003", Total: 2500, CustomerName: "Олег",
	})
	if err != nil {
		t.Fatalf("unexpected service error: %v", err)
	}
	if resp.Order.ID == 0 {
		t.Error("expected order to be created (non-zero ID)")
	}
	if resp.NotifyStatus != "failed" {
		t.Errorf("expected notifyStatus=failed, got %q", resp.NotifyStatus)
	}

	var logEntry models.TelegramSendLog
	if err := db.Where("shop_id = ? AND order_id = ?", shop.ID, resp.Order.ID).First(&logEntry).Error; err != nil {
		t.Fatalf("log not found: %v", err)
	}
	if logEntry.Status != models.StatusFailed {
		t.Errorf("expected log status FAILED, got %q", logEntry.Status)
	}
	if logEntry.Error == nil || *logEntry.Error == "" {
		t.Error("expected error text in log")
	}
}

func TestCreateOrder_ShopNotFound(t *testing.T) {
	db := test.OpenDB(t)
	repos := repository.New(db)
	svc := service.NewOrderService(repos, &test.SuccessClient{}, test.LoadTranslator(t, "../../locales", "ru"))

	_, err := svc.Create(context.Background(), 99999, requests.CreateOrder{
		Number: "T-004", Total: 100, CustomerName: "X",
	})
	if err == nil {
		t.Fatal("expected error for non-existent shop")
	}
}

func TestCreateOrder_NoIntegration_Skipped(t *testing.T) {
	db := test.OpenDB(t)
	shop := &models.Shop{Name: "Shop Without Integration"}
	if err := db.Create(shop).Error; err != nil {
		t.Fatalf("create shop: %v", err)
	}

	tg := &test.SuccessClient{}
	repos := repository.New(db)
	svc := service.NewOrderService(repos, tg, test.LoadTranslator(t, "../../locales", "ru"))

	resp, err := svc.Create(context.Background(), shop.ID, requests.CreateOrder{
		Number: "T-005", Total: 500, CustomerName: "Петр",
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if resp.NotifyStatus != "skipped" {
		t.Errorf("expected notifyStatus=skipped, got %q", resp.NotifyStatus)
	}
	if len(tg.Calls) != 0 {
		t.Errorf("expected 0 telegram calls, got %d", len(tg.Calls))
	}
}

func TestCreateOrder_IntegrationDisabled_Skipped(t *testing.T) {
	db := test.OpenDB(t)
	shop := &models.Shop{Name: "Shop Disabled"}
	if err := db.Create(shop).Error; err != nil {
		t.Fatalf("create shop: %v", err)
	}
	ti := &models.TelegramIntegration{
		ShopID: shop.ID, BotToken: "tok", ChatID: "111", Enabled: false,
	}
	if err := db.Select("shop_id", "bot_token", "chat_id", "enabled", "created_at", "updated_at").Create(ti).Error; err != nil {
		t.Fatalf("create integration: %v", err)
	}

	tg := &test.SuccessClient{}
	repos := repository.New(db)
	svc := service.NewOrderService(repos, tg, test.LoadTranslator(t, "../../locales", "ru"))

	resp, err := svc.Create(context.Background(), shop.ID, requests.CreateOrder{
		Number: "T-006", Total: 200, CustomerName: "Анна",
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if resp.NotifyStatus != "skipped" {
		t.Errorf("expected notifyStatus=skipped, got %q", resp.NotifyStatus)
	}
	if len(tg.Calls) != 0 {
		t.Errorf("expected 0 telegram calls, got %d", len(tg.Calls))
	}
}

func TestCreateOrder_OrderFieldsPersisted(t *testing.T) {
	db := test.OpenDB(t)
	shop := test.SeedShopWithIntegration(t, db)
	repos := repository.New(db)
	svc := service.NewOrderService(repos, &test.SuccessClient{}, test.LoadTranslator(t, "../../locales", "ru"))

	resp, err := svc.Create(context.Background(), shop.ID, requests.CreateOrder{
		Number: "A-999", Total: 3300.50, CustomerName: "Тест",
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if resp.Order.Number != "A-999" {
		t.Errorf("expected number A-999, got %q", resp.Order.Number)
	}
	if resp.Order.Total != 3300.50 {
		t.Errorf("expected total 3300.50, got %v", resp.Order.Total)
	}
	if resp.Order.CustomerName != "Тест" {
		t.Errorf("expected customerName Тест, got %q", resp.Order.CustomerName)
	}
	if resp.Order.ShopID != shop.ID {
		t.Errorf("expected shopID %d, got %d", shop.ID, resp.Order.ShopID)
	}
}
