package models

import (
	"time"

	"gorm.io/gorm"
)

type Shop struct {
	ID        uint           `gorm:"primaryKey;autoIncrement" json:"id"`
	Name      string         `gorm:"not null"                 json:"name"`
	CreatedAt time.Time      `gorm:"autoCreateTime"           json:"createdAt"`
	UpdatedAt time.Time      `gorm:"autoUpdateTime"           json:"updatedAt"`
	DeletedAt gorm.DeletedAt `gorm:"index"                    json:"deletedAt"`
}
