package server_test

import (
	"context"
	"fmt"
	"net/http"
	"testing"
	"time"

	"panda-cooking-go-api/internal/handler"
	"panda-cooking-go-api/internal/model"
	"panda-cooking-go-api/internal/repository"
	"panda-cooking-go-api/internal/server"
	"panda-cooking-go-api/pkg/token"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

func (e *env) withUser(password string) *model.User {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.MinCost)
	require.NoError(e.t, err)
	user := &model.User{ID: userID, Name: "Maria", Email: "maria@email.com", PasswordHash: string(hash)}
	e.users.FindByEmailFn = func(_ context.Context, email string) (*model.User, error) {
		if email != user.Email {
			return nil, gorm.ErrRecordNotFound
		}
		return user, nil
	}
	e.users.FindByIDFn = func(context.Context, string) (*model.User, error) { return user, nil }
	e.sessions.DeleteExpiredFn = func(context.Context, string, time.Time) error { return nil }
	e.sessions.CreateFn = func(context.Context, *model.Session, *model.RefreshToken) error { return nil }
	return user
}

func TestLogin(t *testing.T) {
	t.Run("abre a sessão com o refresh token em cookie HttpOnly", func(t *testing.T) {
		e := newEnv(t)
		e.withUser("panelas-de-barro")

		w := e.do(http.MethodPost, "/api/auth/login", map[string]string{"email": " Maria@Email.com ", "password": "panelas-de-barro"})

		require.Equal(t, http.StatusOK, w.Code, w.Body.String())
		body := decode[map[string]any](t, w)
		assert.NotEmpty(t, body["access_token"])
		assert.Equal(t, "Bearer", body["token_type"])
		assert.NotContains(t, w.Body.String(), "refresh", "o refresh token não pode ir no JSON")

		refresh := cookie(w, handler.RefreshCookie)
		require.NotNil(t, refresh)
		assert.NotEmpty(t, refresh.Value)
		assert.True(t, refresh.HttpOnly)
		assert.Equal(t, http.SameSiteStrictMode, refresh.SameSite)
		assert.Equal(t, "/api/auth", refresh.Path)
		assert.InDelta(t, (7 * 24 * time.Hour).Seconds(), float64(refresh.MaxAge), 5)
		assert.False(t, refresh.Secure, "em desenvolvimento roda em http")

		session := cookie(w, handler.SessionCookie)
		require.NotNil(t, session)
		assert.Equal(t, "1", session.Value)
		assert.Equal(t, "/", session.Path)
		assert.True(t, session.HttpOnly)
	})

	t.Run("em produção os cookies são Secure", func(t *testing.T) {
		cfg := defaultConfig()
		cfg.Env = "production"
		cfg.CORSOrigins = []string{"https://panda.exemplo.com"}
		e := newEnvWith(t, cfg)
		e.withUser("panelas-de-barro")

		w := e.do(http.MethodPost, "/api/auth/login", map[string]string{"email": "maria@email.com", "password": "panelas-de-barro"})

		require.Equal(t, http.StatusOK, w.Code)
		assert.True(t, cookie(w, handler.RefreshCookie).Secure)
		assert.True(t, cookie(w, handler.SessionCookie).Secure)
		assert.Equal(t, "max-age=31536000; includeSubDomains", w.Header().Get("Strict-Transport-Security"))
	})

	t.Run("senha errada é 401 sem cookie", func(t *testing.T) {
		e := newEnv(t)
		e.withUser("panelas-de-barro")

		w := e.do(http.MethodPost, "/api/auth/login", map[string]string{"email": "maria@email.com", "password": "errada"})

		assert.Equal(t, http.StatusUnauthorized, w.Code)
		assert.Equal(t, "e-mail ou senha inválidos", decode[errorBody](t, w).Error)
		assert.Empty(t, w.Result().Cookies())
	})

	t.Run("e-mail inválido é 422 no campo", func(t *testing.T) {
		w := newEnv(t).do(http.MethodPost, "/api/auth/login", map[string]string{"email": "maria", "password": "x"})

		require.Equal(t, http.StatusUnprocessableEntity, w.Code)
		assert.Equal(t, map[string]string{"email": "e-mail inválido"}, decode[errorBody](t, w).Fields)
	})

	t.Run("passou das tentativas por IP recebe 429 com Retry-After", func(t *testing.T) {
		e := newEnv(t)
		e.withUser("panelas-de-barro")
		// e-mails diferentes, para não cair no bloqueio por e-mail antes
		for i := range server.LoginLimit {
			body := map[string]string{"email": fmt.Sprintf("pessoa%d@email.com", i), "password": "errada"}
			require.Equal(t, http.StatusUnauthorized, e.do(http.MethodPost, "/api/auth/login", body).Code)
		}

		w := e.do(http.MethodPost, "/api/auth/login", map[string]string{"email": "maria@email.com", "password": "panelas-de-barro"})

		assert.Equal(t, http.StatusTooManyRequests, w.Code)
		assert.NotEmpty(t, w.Header().Get("Retry-After"))
	})
}

func refreshRequest(value string) *http.Request {
	req := request(http.MethodPost, "/api/auth/refresh", nil)
	if value != "" {
		req.AddCookie(&http.Cookie{Name: handler.RefreshCookie, Value: value})
	}
	return req
}

func TestRefresh(t *testing.T) {
	t.Run("troca o cookie e devolve access token novo", func(t *testing.T) {
		e := newEnv(t)
		e.withUser("panelas-de-barro")
		e.sessions.RotateFn = func(context.Context, string, *model.RefreshToken, time.Time, time.Duration) (repository.RotateResult, error) {
			return repository.RotateResult{Outcome: repository.Rotated, Session: model.Session{UserID: userID}}, nil
		}

		w := e.serve(refreshRequest("token-atual"))

		require.Equal(t, http.StatusOK, w.Code, w.Body.String())
		assert.NotEmpty(t, decode[map[string]any](t, w)["access_token"])
		refresh := cookie(w, handler.RefreshCookie)
		require.NotNil(t, refresh)
		assert.NotEqual(t, "token-atual", refresh.Value)
	})

	t.Run("token recusado apaga os cookies", func(t *testing.T) {
		e := newEnv(t)
		e.sessions.RotateFn = func(context.Context, string, *model.RefreshToken, time.Time, time.Duration) (repository.RotateResult, error) {
			return repository.RotateResult{Outcome: repository.Reused}, nil
		}

		w := e.serve(refreshRequest("token-roubado"))

		assert.Equal(t, http.StatusUnauthorized, w.Code)
		assert.Equal(t, "sessão expirada, entre de novo", decode[errorBody](t, w).Error)
		for _, name := range []string{handler.RefreshCookie, handler.SessionCookie} {
			c := cookie(w, name)
			require.NotNil(t, c, name)
			assert.Negative(t, c.MaxAge, name)
		}
	})

	t.Run("sem cookie é 401", func(t *testing.T) {
		assert.Equal(t, http.StatusUnauthorized, newEnv(t).serve(refreshRequest("")).Code)
	})

	t.Run("navegador de outro site é barrado (CSRF)", func(t *testing.T) {
		e := newEnv(t)
		req := refreshRequest("token-atual")
		req.Header.Set("Sec-Fetch-Site", "cross-site")
		req.Header.Set("Origin", "https://site-malicioso.exemplo")

		w := e.serve(req)

		assert.Equal(t, http.StatusForbidden, w.Code)
		assert.Equal(t, "origem não permitida", decode[errorBody](t, w).Error)
	})

	t.Run("navegador antigo de outra origem é barrado pelo Origin", func(t *testing.T) {
		req := refreshRequest("token-atual")
		req.Header.Set("Origin", "https://site-malicioso.exemplo")

		assert.Equal(t, http.StatusForbidden, newEnv(t).serve(req).Code)
	})

	t.Run("o front, origem confiável, passa", func(t *testing.T) {
		e := newEnv(t)
		e.withUser("panelas-de-barro")
		e.sessions.RotateFn = func(context.Context, string, *model.RefreshToken, time.Time, time.Duration) (repository.RotateResult, error) {
			return repository.RotateResult{Outcome: repository.Rotated, Session: model.Session{UserID: userID}}, nil
		}
		req := refreshRequest("token-atual")
		req.Header.Set("Sec-Fetch-Site", "same-site")
		req.Header.Set("Origin", frontURL)

		assert.Equal(t, http.StatusOK, e.serve(req).Code)
	})
}

func TestLogout(t *testing.T) {
	e := newEnv(t)
	revoked := false
	e.sessions.RevokeByTokenFn = func(context.Context, string, time.Time) error {
		revoked = true
		return nil
	}
	req := request(http.MethodPost, "/api/auth/logout", nil)
	req.AddCookie(&http.Cookie{Name: handler.RefreshCookie, Value: "token-atual"})

	w := e.serve(req)

	assert.Equal(t, http.StatusNoContent, w.Code)
	assert.True(t, revoked)
	assert.Negative(t, cookie(w, handler.RefreshCookie).MaxAge)
	assert.Negative(t, cookie(w, handler.SessionCookie).MaxAge)
}

func TestAuthMiddleware(t *testing.T) {
	e := newEnv(t)
	expired, err := token.Generate(userID, false, secretKey, time.Now().Add(-time.Hour), 15*time.Minute)
	require.NoError(t, err)

	tests := []struct {
		name   string
		header string
	}{
		{name: "sem token", header: ""},
		{name: "esquema errado", header: "Basic dXNlcjpzZW5oYQ=="},
		{name: "token inválido", header: "Bearer abc.def.ghi"},
		{name: "token vencido", header: "Bearer " + expired},
		{name: "alg none", header: "Bearer eyJhbGciOiJub25lIiwidHlwIjoiSldUIn0.eyJzdWIiOiI2ZjFjMmI4ZS0xZDdhLTRjNTUtOWEwZS0yZjRiN2M5ZDFlMzAiLCJpc3MiOiJwYW5kYS1jb29raW5nLWFwaSJ9."},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := request(http.MethodGet, "/api/users/profile", nil)
			if tt.header != "" {
				req.Header.Set("Authorization", tt.header)
			}

			w := e.serve(req)

			assert.Equal(t, http.StatusUnauthorized, w.Code)
			assert.Equal(t, `Bearer realm="panda-cooking"`, w.Header().Get("WWW-Authenticate"))
		})
	}
}
