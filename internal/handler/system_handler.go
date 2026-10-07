package handler

import (
	"context"
	"log/slog"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

// Health responde se a API e o banco estão de pé. É o que o healthcheck do
// Docker e o proxy reverso consultam.
func Health(ping func(context.Context) error) gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(c.Request.Context(), 2*time.Second)
		defer cancel()
		if err := ping(ctx); err != nil {
			slog.ErrorContext(ctx, "banco fora do ar no healthcheck", "erro", err)
			c.JSON(http.StatusServiceUnavailable, gin.H{"status": "unavailable"})
			return
		}
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	}
}

// OpenAPI serve a especificação da API (api/openapi.yaml).
func OpenAPI(spec []byte) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Data(http.StatusOK, "application/yaml; charset=utf-8", spec)
	}
}

// NotFound e MethodNotAllowed mantêm o formato de erro da API também nas
// rotas que não existem.
func NotFound(c *gin.Context) {
	c.JSON(http.StatusNotFound, ErrorBody{Error: "rota não encontrada"})
}

func MethodNotAllowed(c *gin.Context) {
	c.JSON(http.StatusMethodNotAllowed, ErrorBody{Error: "método não permitido nesta rota"})
}
