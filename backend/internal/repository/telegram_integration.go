package repository

import (
	"context"

	"gorm.io/gorm"
	"posiflora/backend/internal/models"
)

type IntegrationRepository struct {
	db *gorm.DB
}

func NewIntegrationRepository(db *gorm.DB) *IntegrationRepository {
	return &IntegrationRepository{db: db}
}

func (r *IntegrationRepository) Upsert(ctx context.Context, ti *models.TelegramIntegration) error {
	existing, err := r.FindByShopID(ctx, ti.ShopID)
	if err != nil {
		return err
	}
	if existing == nil {
		return r.db.WithContext(ctx).
			Select("shop_id", "bot_token", "chat_id", "enabled", "created_at", "updated_at").
			Create(ti).Error
	}
	if err := r.db.WithContext(ctx).
		Model(existing).
		Updates(map[string]interface{}{
			"bot_token": ti.BotToken,
			"chat_id":   ti.ChatID,
			"enabled":   ti.Enabled,
		}).Error; err != nil {
		return err
	}
	ti.ID = existing.ID
	ti.CreatedAt = existing.CreatedAt
	ti.UpdatedAt = existing.UpdatedAt
	return nil
}

func (r *IntegrationRepository) FindByShopID(ctx context.Context, shopID uint) (*models.TelegramIntegration, error) {
	var ti models.TelegramIntegration
	result := r.db.WithContext(ctx).Where("shop_id = ?", shopID).First(&ti)
	if result.Error == gorm.ErrRecordNotFound {
		return nil, nil
	}
	return &ti, result.Error
}
