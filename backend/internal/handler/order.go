package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"posiflora/backend/internal/dto/requests"
	_ "posiflora/backend/internal/dto/responses"
	"posiflora/backend/internal/service"
)

type OrderHandler struct {
	svc *service.OrderService
}

func NewOrderHandler(svc *service.OrderService) *OrderHandler {
	return &OrderHandler{svc: svc}
}

// @Summary      Создать заказ
// @Tags         orders
// @Accept       json
// @Produce      json
// @Param        shopId  path      int                    true  "ID магазина"
// @Param        body    body      requests.CreateOrder   true  "Данные заказа"
// @Success      201     {object}  responses.CreateOrder
// @Failure      400     {object}  map[string]string
// @Failure      404     {object}  map[string]string
// @Router       /shops/{shopId}/orders [post]
func (h *OrderHandler) Create(c *gin.Context) {
	shopID, err := parseShopID(c)
	if err != nil {
		return
	}

	var req requests.CreateOrder
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	result, err := h.svc.Create(c.Request.Context(), shopID, req)
	if err != nil {
		respondServiceError(c, err)
		return
	}

	c.JSON(http.StatusCreated, result)
}
