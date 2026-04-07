package repository

import (
	"context"

	"gorm.io/gorm"
	"posiflora/backend/internal/models"
)

type ShopRepository struct {
	db *gorm.DB
}

func NewShopRepository(db *gorm.DB) *ShopRepository {
	return &ShopRepository{db: db}
}

func (r *ShopRepository) FindByID(ctx context.Context, id uint) (*models.Shop, error) {
	var shop models.Shop
	result := r.db.WithContext(ctx).First(&shop, id)
	if result.Error == gorm.ErrRecordNotFound {
		return nil, nil
	}
	return &shop, result.Error
}
