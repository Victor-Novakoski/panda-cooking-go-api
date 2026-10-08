package token_test

import (
	"testing"
	"time"

	"panda-cooking-go-api/pkg/token"

	"github.com/golang-jwt/jwt/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const secret = "segredo-de-teste-com-32-caracteres!!"

func TestGenerateAndParse(t *testing.T) {
	signed, err := token.Generate("user-uuid", true, secret, time.Now(), time.Minute)
	require.NoError(t, err)

	claims, err := token.Parse(signed, secret)

	require.NoError(t, err)
	assert.Equal(t, "user-uuid", claims.Subject)
	assert.True(t, claims.IsAdm)
}

func TestParse_Recusa(t *testing.T) {
	valid := func(mutate func(*token.Claims)) string {
		c := token.Claims{RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    token.Issuer,
			Subject:   "user-uuid",
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Minute)),
		}}
		mutate(&c)
		s, err := jwt.NewWithClaims(jwt.SigningMethodHS256, c).SignedString([]byte(secret))
		require.NoError(t, err)
		return s
	}

	expired, err := token.Generate("user-uuid", false, secret, time.Now().Add(-time.Hour), 15*time.Minute)
	require.NoError(t, err)
	none, err := jwt.NewWithClaims(jwt.SigningMethodNone, token.Claims{}).SignedString(jwt.UnsafeAllowNoneSignatureType)
	require.NoError(t, err)

	tests := map[string]string{
		"vencido":           expired,
		"outra chave":       mustSign(t, "outra-chave-de-teste-com-32-caracteres"),
		"algoritmo none":    none,
		"sem validade":      valid(func(c *token.Claims) { c.ExpiresAt = nil }),
		"outro emissor":     valid(func(c *token.Claims) { c.Issuer = "outra-api" }),
		"sem usuário":       valid(func(c *token.Claims) { c.Subject = "" }),
		"texto qualquer":    "nao-e-um-jwt",
		"assinatura mexida": expired[:len(expired)-2] + "xx",
	}
	for name, tok := range tests {
		t.Run(name, func(t *testing.T) {
			_, err := token.Parse(tok, secret)
			assert.Error(t, err)
		})
	}
}

func mustSign(t *testing.T, key string) string {
	s, err := token.Generate("user-uuid", false, key, time.Now(), time.Minute)
	require.NoError(t, err)
	return s
}
