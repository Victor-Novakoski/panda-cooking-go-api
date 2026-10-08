package middleware

import (
	"net/http"
	"strings"

	"panda-cooking-go-api/internal/repository"
	"panda-cooking-go-api/pkg/token"

	"github.com/gin-gonic/gin"
)

const (
	UserIDKey = "userID"
	IsAdmKey  = "isAdm"
)

// Auth exige o access token no cabeçalho Authorization: Bearer <token>.
func Auth(secretKey string) gin.HandlerFunc {
	return func(c *gin.Context) {
		scheme, raw, ok := strings.Cut(c.GetHeader("Authorization"), " ")
		if !ok || !strings.EqualFold(scheme, "bearer") || raw == "" {
			unauthorized(c, "token não fornecido")
			return
		}

		claims, err := token.Parse(strings.TrimSpace(raw), secretKey)
		if err != nil || !repository.IsUUID(claims.Subject) {
			unauthorized(c, "token inválido ou expirado")
			return
		}

		c.Set(UserIDKey, claims.Subject)
		c.Set(IsAdmKey, claims.IsAdm)
		c.Next()
	}
}

func unauthorized(c *gin.Context, msg string) {
	c.Header("WWW-Authenticate", `Bearer realm="panda-cooking"`)
	c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": msg})
}
