package repository_test

import (
	"context"
	"testing"

	"posiflora/backend/internal/repository"
	"posiflora/backend/internal/test"
)

func TestShopRepository_FindByID_Found(t *testing.T) {
	db := test.OpenDB(t)
	shop := test.SeedShop(t, db, "Repo Test Shop")

	repo := repository.NewShopRepository(db)
	result, err := repo.FindByID(context.Background(), shop.ID)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result == nil {
		t.Fatal("expected shop, got nil")
		return
	}
	if result.ID != shop.ID {
		t.Errorf("expected ID %d, got %d", shop.ID, result.ID)
	}
	if result.Name != shop.Name {
		t.Errorf("expected name %q, got %q", shop.Name, result.Name)
	}
}

func TestShopRepository_FindByID_NotFound(t *testing.T) {
	db := test.OpenDB(t)
	repo := repository.NewShopRepository(db)

	result, err := repo.FindByID(context.Background(), 99999)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result != nil {
		t.Errorf("expected nil, got shop with ID %d", result.ID)
	}
}
