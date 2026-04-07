package service

import (
	"context"
	"fmt"
	"time"

	"posiflora/backend/internal/dto/requests"
	"posiflora/backend/internal/dto/responses"
	"posiflora/backend/internal/models"
	"posiflora/backend/internal/repository"
)

type TelegramService struct {
	repos *repository.Repositories
}

func NewTelegramService(repos *repository.Repositories) *TelegramService {
	return &TelegramService{repos: repos}
}

func (s *TelegramService) Connect(ctx context.Context, shopID uint, req requests.ConnectTelegram) (*responses.ConnectTelegram, error) {
	if _, err := s.requireShop(ctx, shopID); err != nil {
		return nil, err
	}

	botToken := req.BotToken
	if botToken == "" {
		existing, err := s.repos.Integration.FindByShopID(ctx, shopID)
		if err != nil {
			return nil, err
		}
		if existing == nil {
			return nil, fmt.Errorf("botToken is required for initial setup")
		}
		botToken = existing.BotToken
	}

	ti := &models.TelegramIntegration{
		ShopID:   shopID,
		BotToken: botToken,
		ChatID:   req.ChatID,
		Enabled:  req.Enabled,
	}
	if err := s.repos.Integration.Upsert(ctx, ti); err != nil {
		return nil, err
	}
	return &responses.ConnectTelegram{
		ID:        ti.ID,
		ShopID:    ti.ShopID,
		ChatID:    ti.ChatID,
		Enabled:   ti.Enabled,
		CreatedAt: ti.CreatedAt,
		UpdatedAt: ti.UpdatedAt,
	}, nil
}

func (s *TelegramService) GetStatus(ctx context.Context, shopID uint) (*responses.TelegramStatus, error) {
	if _, err := s.requireShop(ctx, shopID); err != nil {
		return nil, err
	}

	ti, err := s.repos.Integration.FindByShopID(ctx, shopID)
	if err != nil {
		return nil, err
	}
	if ti == nil {
		return &responses.TelegramStatus{}, nil
	}

	since := time.Now().AddDate(0, 0, -7)
	sentCount, failedCount, lastSentAt, err := s.repos.SendLog.GetCounts(ctx, shopID, since)
	if err != nil {
		return nil, err
	}

	return &responses.TelegramStatus{
		Enabled:       ti.Enabled,
		ChatID:        ti.ChatID,
		LastSentAt:    lastSentAt,
		SentCount7d:   sentCount,
		FailedCount7d: failedCount,
	}, nil
}

func (s *TelegramService) requireShop(ctx context.Context, shopID uint) (*models.Shop, error) {
	shop, err := s.repos.Shop.FindByID(ctx, shopID)
	if err != nil {
		return nil, err
	}
	if shop == nil {
		return nil, fmt.Errorf("shop %d not found", shopID)
	}
	return shop, nil
}
