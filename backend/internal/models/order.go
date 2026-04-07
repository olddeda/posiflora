package models

import (
	"time"

	"gorm.io/gorm"
)

type Order struct {
	ID           uint           `gorm:"primaryKey;autoIncrement"    json:"id"`
	ShopID       uint           `gorm:"not null;index"              json:"shopId"`
	Number       string         `gorm:"not null"                    json:"number"`
	Total        float64        `gorm:"type:numeric(12,2);not null" json:"total"`
	CustomerName string         `gorm:"not null"                    json:"customerName"`
	CreatedAt    time.Time      `gorm:"autoCreateTime"              json:"createdAt"`
	UpdatedAt    time.Time      `gorm:"autoUpdateTime"              json:"updatedAt"`
	DeletedAt    gorm.DeletedAt `gorm:"index"                       json:"deletedAt"`
}
