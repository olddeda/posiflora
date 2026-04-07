package repository_test

import (
	"context"
	"testing"

	"posiflora/backend/internal/models"
	"posiflora/backend/internal/repository"
	"posiflora/backend/internal/test"
)

func TestOrderRepository_Create(t *testing.T) {
	db := test.OpenDB(t)
	shop := test.SeedShop(t, db, "Repo Test Shop")

	repo := repository.NewOrderRepository(db)
	order := &models.Order{
		ShopID:       shop.ID,
		Number:       "R-001",
		Total:        1234.56,
		CustomerName: "Тест",
	}

	if err := repo.Create(context.Background(), order); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if order.ID == 0 {
		t.Error("expected non-zero ID after create")
	}

	var persisted models.Order
	if err := db.First(&persisted, order.ID).Error; err != nil {
		t.Fatalf("order not found in DB: %v", err)
	}
	if persisted.Number != "R-001" {
		t.Errorf("expected number R-001, got %q", persisted.Number)
	}
	if persisted.Total != 1234.56 {
		t.Errorf("expected total 1234.56, got %v", persisted.Total)
	}
}
