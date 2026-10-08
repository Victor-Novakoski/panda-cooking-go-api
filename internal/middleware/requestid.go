package middleware

import (
	"crypto/rand"
	"encoding/hex"
	"regexp"

	"panda-cooking-go-api/internal/applog"

	"github.com/gin-gonic/gin"
)

const RequestIDHeader = "X-Request-ID"

// um id que veio de fora só é aceito se for curto e sem caractere estranho,
// para não sujar o log (quebra de linha, texto enorme)
var validRequestID = regexp.MustCompile(`^[A-Za-z0-9._-]{1,64}$`)

// RequestID dá um id a cada requisição: reaproveita o X-Request-ID que veio
// do front ou do proxy, ou gera um. O id volta na resposta e vai em todo log
// da requisição, para ligar o erro que o usuário viu à linha do log.
func RequestID() gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.GetHeader(RequestIDHeader)
		if !validRequestID.MatchString(id) {
			id = newRequestID()
		}
		c.Header(RequestIDHeader, id)
		c.Request = c.Request.WithContext(applog.WithRequestID(c.Request.Context(), id))
		c.Next()
	}
}

func newRequestID() string {
	b := make([]byte, 12)
	_, _ = rand.Read(b)
	return hex.EncodeToString(b)
}
