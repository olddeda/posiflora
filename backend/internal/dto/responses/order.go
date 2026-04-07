package responses

import "posiflora/backend/internal/models"

type CreateOrder struct {
	Order        models.Order `json:"order"`
	NotifyStatus string       `json:"notifyStatus"`
}
