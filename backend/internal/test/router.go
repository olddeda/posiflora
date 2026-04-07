package test

import (
	"testing"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"posiflora/backend/internal/handler"
	"posiflora/backend/internal/repository"
	"posiflora/backend/internal/service"
	"posiflora/backend/internal/telegram"
)

func NewTestRouter(t *testing.T, db *gorm.DB, localesDir string) *gin.Engine {
	t.Helper()
	translator := LoadTranslator(t, localesDir, "ru")
	repos := repository.New(db)
	tgSvc := service.NewTelegramService(repos)
	orderSvc := service.NewOrderService(repos, telegram.NewMockClient(), translator)
	tgHandler := handler.NewTelegramHandler(tgSvc)
	orderHandler := handler.NewOrderHandler(orderSvc)
	return handler.NewRouter(tgHandler, orderHandler, "*")
}
