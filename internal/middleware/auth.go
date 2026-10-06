package middleware

import (
	"net/http"
	"strings"

	"panda-cooking-go-api/pkg/token"

	"github.com/gin-gonic/gin"
)

const UserIDKey = "userID"
const IsAdmKey = "isAdm"

func Auth(secretKey string) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "token não fornecido"})
			return
		}

		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) != 2 || strings.ToLower(parts[0]) != "bearer" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "formato do token inválido"})
			return
		}

		claims, err := token.Parse(parts[1], secretKey)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "token inválido ou expirado"})
			return
		}

		c.Set(UserIDKey, claims.UserID)
		c.Set(IsAdmKey, claims.IsAdm)
		c.Next()
	}
}
