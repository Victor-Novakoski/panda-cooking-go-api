package middleware

import (
	"log/slog"
	"net/http"
	"runtime/debug"
	"time"

	"github.com/gin-gonic/gin"
)

// AccessLog registra cada requisição: status, rota, tempo, IP e usuário.
// 401, 403 e 429 saem como aviso (tentativa de acesso indevido ou abuso) e
// 5xx como erro. O /health respondendo 200 fica de fora, para o healthcheck
// do Docker não encher o log.
func AccessLog(log *slog.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		c.Next()

		status := c.Writer.Status()
		if c.Request.URL.Path == "/health" && status == http.StatusOK {
			return
		}

		level := slog.LevelInfo
		switch {
		case status >= 500:
			level = slog.LevelError
		case status == http.StatusUnauthorized, status == http.StatusForbidden, status == http.StatusTooManyRequests:
			level = slog.LevelWarn
		}

		attrs := []slog.Attr{
			slog.String("method", c.Request.Method),
			slog.String("path", c.Request.URL.Path),
			slog.String("route", c.FullPath()),
			slog.Int("status", status),
			slog.Float64("duration_ms", float64(time.Since(start).Microseconds())/1000),
			slog.String("ip", c.ClientIP()),
		}
		if userID := c.GetString(UserIDKey); userID != "" {
			attrs = append(attrs, slog.String("user_id", userID))
		}
		log.LogAttrs(c.Request.Context(), level, "requisição", attrs...)
	}
}

// Recovery transforma um panic em 500 com a mensagem genérica de sempre e
// registra o erro com a pilha no log, em vez de derrubar a conexão.
func Recovery(log *slog.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if rec := recover(); rec != nil {
				if rec == http.ErrAbortHandler { //nolint:errorlint // o net/http compara assim
					panic(rec)
				}
				log.ErrorContext(c.Request.Context(), "panic na requisição",
					"panic", rec, "path", c.Request.URL.Path, "stack", string(debug.Stack()))
				if !c.Writer.Written() {
					c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "erro interno, tente novamente mais tarde"})
				} else {
					c.Abort()
				}
			}
		}()
		c.Next()
	}
}
