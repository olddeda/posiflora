package test

import (
	"os"
	"testing"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"

	"posiflora/backend/internal/models"
)

func OpenDB(t testing.TB) *gorm.DB {
	t.Helper()
	dsn := os.Getenv("TEST_DATABASE_URL")
	if dsn == "" {
		dsn = "postgres://posiflora:posiflora@localhost:5432/posiflora_test?sslmode=disable"
	}
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
	})
	if err != nil {
		t.Skipf("postgres unavailable (%v), skipping integration tests", err)
	}
	if err := db.AutoMigrate(
		&models.Shop{},
		&models.TelegramIntegration{},
		&models.Order{},
		&models.TelegramSendLog{},
	); err != nil {
		t.Fatalf("migrate: %v", err)
	}
	t.Cleanup(func() {
		db.Exec("DELETE FROM telegram_send_logs")
		db.Exec("DELETE FROM orders")
		db.Exec("DELETE FROM telegram_integrations")
		db.Exec("DELETE FROM shops")
	})
	return db
}

func SeedShop(t testing.TB, db *gorm.DB, name string) *models.Shop {
	t.Helper()
	shop := &models.Shop{Name: name}
	if err := db.Create(shop).Error; err != nil {
		t.Fatalf("create shop: %v", err)
	}
	return shop
}

func SeedShopWithIntegration(t testing.TB, db *gorm.DB) *models.Shop {
	t.Helper()
	shop := SeedShop(t, db, "Test Shop")
	ti := &models.TelegramIntegration{
		ShopID:   shop.ID,
		BotToken: "test-token",
		ChatID:   "123456",
		Enabled:  true,
	}
	if err := db.Create(ti).Error; err != nil {
		t.Fatalf("create integration: %v", err)
	}
	return shop
}
