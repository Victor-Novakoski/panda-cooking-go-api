package middleware

import (
	"context"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

// SecurityHeaders manda o navegador se proteger nas respostas da API: não
// adivinhar o tipo do conteúdo, não abrir a resposta em iframe, não executar
// nada que venha dela e não guardar em cache (as respostas podem ter dado
// pessoal). Com HTTPS, o HSTS faz o navegador nunca mais usar http.
func SecurityHeaders(https bool) gin.HandlerFunc {
	return func(c *gin.Context) {
		h := c.Writer.Header()
		h.Set("X-Content-Type-Options", "nosniff")
		h.Set("X-Frame-Options", "DENY")
		h.Set("Referrer-Policy", "no-referrer")
		h.Set("Content-Security-Policy", "default-src 'none'; frame-ancestors 'none'")
		h.Set("Cross-Origin-Resource-Policy", "same-origin")
		h.Set("Cache-Control", "no-store")
		if https {
			h.Set("Strict-Transport-Security", "max-age=31536000; includeSubDomains")
		}
		c.Next()
	}
}

// BodyLimit corta o corpo da requisição em max bytes. Quem passar recebe 413
// (ver bindJSON no handler), em vez de o servidor ler um corpo gigante.
func BodyLimit(max int64) gin.HandlerFunc {
	return func(c *gin.Context) {
		if c.Request.Body != nil {
			c.Request.Body = http.MaxBytesReader(c.Writer, c.Request.Body, max)
		}
		c.Next()
	}
}

// Timeout dá um prazo para cada requisição. O contexto chega até as
// consultas do banco, que são canceladas quando o prazo acaba.
func Timeout(d time.Duration) gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(c.Request.Context(), d)
		defer cancel()
		c.Request = c.Request.WithContext(ctx)
		c.Next()
	}
}

// CrossOrigin recusa requisição de navegador vinda de outro site nos métodos
// que mudam dados (POST, PUT, PATCH, DELETE). Usa o http.CrossOriginProtection
// da biblioteca padrão, que olha o Sec-Fetch-Site e, em navegador antigo, o
// Origin. É a proteção contra CSRF da rota de refresh, que usa cookie, junto
// com o SameSite=Strict. As origens em trusted (o front) podem chamar.
// Requisição sem esses cabeçalhos (curl, o proxy do front) passa: não há
// navegador de vítima mandando cookie.
func CrossOrigin(trusted []string) (gin.HandlerFunc, error) {
	protection := http.NewCrossOriginProtection()
	for _, origin := range trusted {
		if err := protection.AddTrustedOrigin(origin); err != nil {
			return nil, err
		}
	}
	return func(c *gin.Context) {
		if err := protection.Check(c.Request); err != nil {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "origem não permitida"})
			return
		}
		c.Next()
	}, nil
}
