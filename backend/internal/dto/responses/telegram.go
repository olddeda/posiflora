package responses

import "time"

type TelegramStatus struct {
	Enabled       bool       `json:"enabled"`
	ChatID        string     `json:"chatId"`
	LastSentAt    *time.Time `json:"lastSentAt"`
	SentCount7d   int64      `json:"sentCount7d"`
	FailedCount7d int64      `json:"failedCount7d"`
}

type ConnectTelegram struct {
	ID        uint      `json:"id"`
	ShopID    uint      `json:"shopId"`
	ChatID    string    `json:"chatId"`
	Enabled   bool      `json:"enabled"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}
