package repository

import "gorm.io/gorm"

type Repositories struct {
	Shop        *ShopRepository
	Order       *OrderRepository
	Integration *IntegrationRepository
	SendLog     *TelegramSendLogRepository
}

func New(db *gorm.DB) *Repositories {
	return &Repositories{
		Shop:        NewShopRepository(db),
		Order:       NewOrderRepository(db),
		Integration: NewIntegrationRepository(db),
		SendLog:     NewTelegramSendLogRepository(db),
	}
}
