// Package token gera e valida o access token (JWT HS256).
package token

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

// Issuer identifica os tokens desta API; token de outro emissor é recusado.
const Issuer = "panda-cooking-api"

// Claims: Subject é o id do usuário; IsAdm vai no token para o middleware
// não consultar o banco a cada requisição.
type Claims struct {
	jwt.RegisteredClaims
	IsAdm bool `json:"adm"`
}

func Generate(userID string, isAdm bool, secretKey string, now time.Time, ttl time.Duration) (string, error) {
	claims := Claims{
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    Issuer,
			Subject:   userID,
			IssuedAt:  jwt.NewNumericDate(now),
			ExpiresAt: jwt.NewNumericDate(now.Add(ttl)),
		},
		IsAdm: isAdm,
	}
	return jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString([]byte(secretKey))
}

// Parse valida assinatura, algoritmo (só HS256: nada de "none" ou troca por
// RS256), emissor e validade.
func Parse(tokenStr, secretKey string) (*Claims, error) {
	claims := &Claims{}
	_, err := jwt.ParseWithClaims(tokenStr, claims,
		func(*jwt.Token) (any, error) { return []byte(secretKey), nil },
		jwt.WithValidMethods([]string{jwt.SigningMethodHS256.Alg()}),
		jwt.WithIssuer(Issuer),
		jwt.WithExpirationRequired(),
		jwt.WithLeeway(5*time.Second),
	)
	if err != nil {
		return nil, err
	}
	if claims.Subject == "" {
		return nil, errors.New("token sem usuário")
	}
	return claims, nil
}
