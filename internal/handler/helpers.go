package handler

import (
	"errors"
	"log"
	"net/http"

	"panda-cooking-go-api/internal/service"

	"github.com/gin-gonic/gin"
)

// msgInternal é o que o cliente vê em qualquer erro inesperado.
// O detalhe fica só no log, para não expor banco, SQL ou caminhos internos.
const msgInternal = "erro interno, tente novamente mais tarde"

var statusByKind = map[service.Kind]int{
	service.KindNotFound:        http.StatusNotFound,
	service.KindForbidden:       http.StatusForbidden,
	service.KindConflict:        http.StatusConflict,
	service.KindUnauthorized:    http.StatusUnauthorized,
	service.KindTooManyRequests: http.StatusTooManyRequests,
	service.KindInvalid:         http.StatusBadRequest,
}

// respondError responde um erro do service: erro conhecido sai com o status e
// a mensagem dele; o resto vira 500 com mensagem genérica.
func respondError(c *gin.Context, err error) {
	var appErr *service.Error
	if errors.As(err, &appErr) {
		if status, ok := statusByKind[appErr.Kind]; ok {
			c.JSON(status, gin.H{"error": appErr.Message})
			return
		}
	}

	log.Printf("erro interno em %s %s: %v", c.Request.Method, c.FullPath(), err)
	c.JSON(http.StatusInternalServerError, gin.H{"error": msgInternal})
}
