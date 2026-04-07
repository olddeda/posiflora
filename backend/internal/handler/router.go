package handler

import (
	"log"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func NewRouter(tg *TelegramHandler, order *OrderHandler, allowedOrigins string) *gin.Engine {
	r := gin.Default()

	if err := r.SetTrustedProxies(nil); err != nil {
		log.Fatalf("set trusted proxies: %v", err)
	}

	r.Use(corsMiddleware(allowedOrigins))

	r.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})

	r.GET("/docs", scalarHandler)
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	shops := r.Group("/shops/:shopId")
	{
		shops.POST("/telegram/connect", tg.Connect)
		shops.GET("/telegram/status", tg.GetStatus)
		shops.POST("/orders", order.Create)
	}

	return r
}

func scalarHandler(c *gin.Context) {
	c.Header("Content-Type", "text/html; charset=utf-8")
	c.String(http.StatusOK, `<!DOCTYPE html>
<html>
<head>
  <title>Posiflora API</title>
  <meta charset="utf-8"/>
  <meta name="viewport" content="width=device-width, initial-scale=1"/>
</head>
<body>
  <script id="api-reference" data-url="/swagger/doc.json"></script>
  <script src="https://cdn.jsdelivr.net/npm/@scalar/api-reference"></script>
</body>
</html>`)
}

func corsMiddleware(allowedOrigins string) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Header("Access-Control-Allow-Origin", allowedOrigins)
		c.Header("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		c.Header("Access-Control-Allow-Headers", "Content-Type, Authorization")

		if c.Request.Method == http.MethodOptions {
			c.AbortWithStatus(http.StatusNoContent)
			return
		}
		c.Next()
	}
}

func respondServiceError(c *gin.Context, err error) {
	msg := err.Error()
	if strings.Contains(msg, "not found") {
		c.JSON(http.StatusNotFound, gin.H{"error": msg})
		return
	}
	if strings.Contains(msg, "required") || strings.Contains(msg, "invalid") {
		c.JSON(http.StatusBadRequest, gin.H{"error": msg})
		return
	}
	log.Printf("internal error: %v", err)
	c.JSON(http.StatusInternalServerError, gin.H{"error": "internal error"})
}
