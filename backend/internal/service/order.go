package service

import (
	"context"
	"fmt"
	"log"

	"posiflora/backend/internal/dto/requests"
	"posiflora/backend/internal/dto/responses"
	"posiflora/backend/internal/i18n"
	"posiflora/backend/internal/models"
	"posiflora/backend/internal/repository"
	"posiflora/backend/internal/telegram"
)

type OrderService struct {
	repos      *repository.Repositories
	tgClient   telegram.Client
	translator *i18n.Translator
}

func NewOrderService(repos *repository.Repositories, tgClient telegram.Client, translator *i18n.Translator) *OrderService {
	return &OrderService{repos: repos, tgClient: tgClient, translator: translator}
}

func (s *OrderService) Create(ctx context.Context, shopID uint, req requests.CreateOrder) (*responses.CreateOrder, error) {
	shop, err := s.repos.Shop.FindByID(ctx, shopID)
	if err != nil {
		return nil, err
	}
	if shop == nil {
		return nil, fmt.Errorf("shop %d not found", shopID)
	}

	order := &models.Order{
		ShopID:       shopID,
		Number:       req.Number,
		Total:        req.Total,
		CustomerName: req.CustomerName,
	}
	if err := s.repos.Order.Create(ctx, order); err != nil {
		return nil, err
	}

	notifyStatus := s.sendNotification(ctx, shopID, order)

	return &responses.CreateOrder{
		Order:        *order,
		NotifyStatus: notifyStatus,
	}, nil
}

func (s *OrderService) sendNotification(ctx context.Context, shopID uint, order *models.Order) string {
	integration, err := s.repos.Integration.FindByShopID(ctx, shopID)
	if err != nil || integration == nil || !integration.Enabled {
		return "skipped"
	}

	existing, err := s.repos.SendLog.FindByShopAndOrder(ctx, shopID, order.ID)
	if err != nil {
		return "skipped"
	}
	if existing != nil {
		return "skipped"
	}

	msg, err := s.translator.Render("order.notification", struct {
		Number       string
		Total        float64
		CustomerName string
	}{
		Number:       order.Number,
		Total:        order.Total,
		CustomerName: order.CustomerName,
	})
	if err != nil {
		msg = fmt.Sprintf("%s | %.0f ₽ | %s", order.Number, order.Total, order.CustomerName)
	}

	sendErr := s.tgClient.SendMessage(integration.BotToken, integration.ChatID, msg)

	logEntry := &models.TelegramSendLog{
		ShopID:  shopID,
		OrderID: order.ID,
		Message: msg,
		Status:  models.StatusSent,
	}
	if sendErr != nil {
		errStr := sendErr.Error()
		logEntry.Status = models.StatusFailed
		logEntry.Error = &errStr
	}

	if err := s.repos.SendLog.Create(ctx, logEntry); err != nil {
		log.Printf("send log create: %v", err)
	}

	if sendErr != nil {
		return "failed"
	}
	return "sent"
}
