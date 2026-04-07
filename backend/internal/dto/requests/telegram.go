package requests

type ConnectTelegram struct {
	BotToken string `json:"botToken"`
	ChatID   string `json:"chatId"   binding:"required"`
	Enabled  bool   `json:"enabled"`
}
