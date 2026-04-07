package handler

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"posiflora/backend/internal/dto/requests"
	_ "posiflora/backend/internal/dto/responses"
	"posiflora/backend/internal/service"
)

type TelegramHandler struct {
	svc *service.TelegramService
}

func NewTelegramHandler(svc *service.TelegramService) *TelegramHandler {
	return &TelegramHandler{svc: svc}
}

// @Summary      Подключить Telegram-бота
// @Tags         telegram
// @Accept       json
// @Produce      json
// @Param        shopId  path      int                           true  "ID магазина"
// @Param        body    body      requests.ConnectTelegram      true  "Параметры подключения"
// @Success      200     {object}  responses.ConnectTelegram
// @Failure      400     {object}  map[string]string
// @Failure      404     {object}  map[string]string
// @Router       /shops/{shopId}/telegram/connect [post]
func (h *TelegramHandler) Connect(c *gin.Context) {
	shopID, err := parseShopID(c)
	if err != nil {
		return
	}

	var req requests.ConnectTelegram
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	integration, err := h.svc.Connect(c.Request.Context(), shopID, req)
	if err != nil {
		respondServiceError(c, err)
		return
	}

	c.JSON(http.StatusOK, integration)
}

// @Summary      Статус Telegram-интеграции
// @Tags         telegram
// @Produce      json
// @Param        shopId  path      int  true  "ID магазина"
// @Success      200     {object}  responses.TelegramStatus
// @Failure      404     {object}  map[string]string
// @Router       /shops/{shopId}/telegram/status [get]
func (h *TelegramHandler) GetStatus(c *gin.Context) {
	shopID, err := parseShopID(c)
	if err != nil {
		return
	}

	status, err := h.svc.GetStatus(c.Request.Context(), shopID)
	if err != nil {
		respondServiceError(c, err)
		return
	}

	c.JSON(http.StatusOK, status)
}

func parseShopID(c *gin.Context) (uint, error) {
	id, err := strconv.ParseUint(c.Param("shopId"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid shopId"})
		return 0, err
	}
	return uint(id), nil
}
