package models

import (
	"time"

	"gorm.io/gorm"
)

type TelegramIntegration struct {
	ID        uint           `gorm:"primaryKey;autoIncrement" json:"id"`
	ShopID    uint           `gorm:"not null;uniqueIndex"     json:"shopId"`
	BotToken  string         `gorm:"not null"                 json:"botToken"`
	ChatID    string         `gorm:"not null"                 json:"chatId"`
	Enabled   bool           `gorm:"not null"                 json:"enabled"`
	CreatedAt time.Time      `gorm:"autoCreateTime"           json:"createdAt"`
	UpdatedAt time.Time      `gorm:"autoUpdateTime"           json:"updatedAt"`
	DeletedAt gorm.DeletedAt `gorm:"index"                    json:"deletedAt"`
}
