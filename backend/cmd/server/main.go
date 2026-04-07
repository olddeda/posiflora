// @title           Posiflora API
// @version         1.0
// @description     Telegram-интеграция для магазинов Posiflora
// @host            localhost:8080
// @BasePath        /
// @accept          json
// @produce         json

package main

import (
	"fmt"
	"log"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"

	_ "posiflora/backend/docs"
	"posiflora/backend/internal/config"
	"posiflora/backend/internal/db"
	"posiflora/backend/internal/handler"
	"posiflora/backend/internal/i18n"
	"posiflora/backend/internal/repository"
	"posiflora/backend/internal/service"
	"posiflora/backend/internal/telegram"
)

func main() {
	cfg := config.Load()

	gormDB, err := openDB(cfg.DatabaseURL)
	if err != nil {
		log.Fatalf("db: %v", err)
	}

	if err := db.RunMigrations(gormDB, "migrations"); err != nil {
		log.Fatalf("migrations: %v", err)
	}

	var tgClient telegram.Client
	if cfg.TelegramEnabled {
		log.Println("Telegram: REAL mode")
		tgClient = telegram.NewHTTPClient()
	} else {
		log.Println("Telegram: MOCK mode")
		tgClient = telegram.NewMockClient()
	}

	translator, err := i18n.Load(cfg.LocalesDir, cfg.Locale)
	if err != nil {
		log.Fatalf("i18n: %v", err)
	}

	repos := repository.New(gormDB)

	tgService := service.NewTelegramService(repos)
	orderService := service.NewOrderService(repos, tgClient, translator)

	tgHandler := handler.NewTelegramHandler(tgService)
	orderHandler := handler.NewOrderHandler(orderService)

	router := handler.NewRouter(tgHandler, orderHandler, cfg.AllowedOrigins)

	addr := fmt.Sprintf(":%s", cfg.Port)
	log.Printf("Listening on %s", addr)
	if err := router.Run(addr); err != nil {
		log.Fatalf("server: %v", err)
	}
}

func openDB(dsn string) (*gorm.DB, error) {
	cfg := &gorm.Config{
		Logger: logger.Default.LogMode(logger.Warn),
	}
	var gormDB *gorm.DB
	var err error
	for i := range 5 {
		gormDB, err = gorm.Open(postgres.Open(dsn), cfg)
		if err == nil {
			return gormDB, nil
		}
		log.Printf("db not ready (attempt %d/5): %v", i+1, err)
		time.Sleep(2 * time.Second)
	}
	return nil, err
}
