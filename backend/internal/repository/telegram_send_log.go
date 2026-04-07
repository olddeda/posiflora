package repository

import (
	"context"
	"errors"
	"time"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"posiflora/backend/internal/models"
)

type TelegramSendLogRepository struct {
	db *gorm.DB
}

func NewTelegramSendLogRepository(db *gorm.DB) *TelegramSendLogRepository {
	return &TelegramSendLogRepository{db: db}
}

func (r *TelegramSendLogRepository) FindByShopAndOrder(ctx context.Context, shopID, orderID uint) (*models.TelegramSendLog, error) {
	var entry models.TelegramSendLog
	result := r.db.WithContext(ctx).
		Where("shop_id = ? AND order_id = ?", shopID, orderID).
		First(&entry)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	if result.Error != nil {
		return nil, result.Error
	}
	return &entry, nil
}

func (r *TelegramSendLogRepository) Create(ctx context.Context, entry *models.TelegramSendLog) error {
	return r.db.WithContext(ctx).
		Clauses(clause.OnConflict{DoNothing: true}).
		Create(entry).Error
}

func (r *TelegramSendLogRepository) GetCounts(ctx context.Context, shopID uint, since time.Time) (sentCount, failedCount int64, lastSentAt *time.Time, err error) {
	if err = r.db.WithContext(ctx).Model(&models.TelegramSendLog{}).
		Where("shop_id = ? AND status = ? AND sent_at >= ?", shopID, models.StatusSent, since).
		Count(&sentCount).Error; err != nil {
		return
	}

	if err = r.db.WithContext(ctx).Model(&models.TelegramSendLog{}).
		Where("shop_id = ? AND status = ? AND sent_at >= ?", shopID, models.StatusFailed, since).
		Count(&failedCount).Error; err != nil {
		return
	}

	var latest models.TelegramSendLog
	result := r.db.WithContext(ctx).
		Where("shop_id = ? AND status = ?", shopID, models.StatusSent).
		Order("sent_at DESC").
		First(&latest)
	if result.Error == nil {
		t := latest.SentAt
		lastSentAt = &t
	} else if !errors.Is(result.Error, gorm.ErrRecordNotFound) {
		err = result.Error
	}
	return
}
