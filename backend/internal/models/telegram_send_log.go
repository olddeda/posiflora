package models

import (
	"time"

	"gorm.io/gorm"
)

type SendStatus string

const (
	StatusSent   SendStatus = "SENT"
	StatusFailed SendStatus = "FAILED"
)

type TelegramSendLog struct {
	ID        uint           `gorm:"primaryKey;autoIncrement"           json:"id"`
	ShopID    uint           `gorm:"not null;uniqueIndex:idx_shop_order" json:"shopId"`
	OrderID   uint           `gorm:"not null;uniqueIndex:idx_shop_order" json:"orderId"`
	Message   string         `gorm:"not null"                            json:"message"`
	Status    SendStatus     `gorm:"type:varchar(10);not null"           json:"status"`
	Error     *string        `                                           json:"error,omitempty"`
	SentAt    time.Time      `gorm:"column:sent_at;autoCreateTime"       json:"sentAt"`
	UpdatedAt time.Time      `gorm:"autoUpdateTime"                      json:"updatedAt"`
	DeletedAt gorm.DeletedAt `gorm:"index"                               json:"deletedAt"`
}
