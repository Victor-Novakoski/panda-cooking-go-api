package middleware

import (
	"math"
	"net/http"
	"strconv"

	"panda-cooking-go-api/internal/ratelimit"

	"github.com/gin-gonic/gin"
)

// RateLimitByIP limita as requisições de cada IP. Passou do limite, responde
// 429 com Retry-After em segundos.
func RateLimitByIP(l *ratelimit.Limiter) gin.HandlerFunc {
	return func(c *gin.Context) {
		if ok, wait := l.Allow(c.ClientIP()); !ok {
			c.Header("Retry-After", strconv.Itoa(int(math.Ceil(wait.Seconds()))))
			c.AbortWithStatusJSON(http.StatusTooManyRequests, gin.H{"error": "muitas requisições, tente novamente mais tarde"})
			return
		}
		c.Next()
	}
}
