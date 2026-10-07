package service_test

import (
	"bytes"
	"context"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"log/slog"
	"testing"
	"time"

	"panda-cooking-go-api/internal/model"
	"panda-cooking-go-api/internal/repository"
	"panda-cooking-go-api/internal/service"
	"panda-cooking-go-api/internal/service/mocks"
	"panda-cooking-go-api/pkg/token"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

const (
	secretKey = "chave-de-teste-com-mais-de-32-caracteres"
	userID    = "6f1c2b8e-1d7a-4c55-9a0e-2f4b7c9d1e30"
	otherID   = "0b9e7c41-5a3d-4e21-8f6a-7c2d9b1e4f58"
)

var (
	ctx = context.Background()
	// relógio fixo do service; perto do real para o JWT emitido ainda valer
	now = time.Now()
)

func discardLog() *slog.Logger { return slog.New(slog.DiscardHandler) }

func sha(s string) string {
	sum := sha256.Sum256([]byte(s))
	return hex.EncodeToString(sum[:])
}

// userWithPassword cria um usuário com hash de custo baixo, para o teste
// não gastar tempo com bcrypt.
func userWithPassword(t *testing.T, password string) *model.User {
	t.Helper()
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.MinCost)
	require.NoError(t, err)
	return &model.User{ID: userID, Name: "Maria", Email: "maria@email.com", PasswordHash: string(hash)}
}

func newAuth(users *mocks.UserRepoMock, sessions *mocks.SessionRepoMock, log *slog.Logger) *service.AuthService {
	s := service.NewAuthService(users, sessions, secretKey, log)
	s.Now = func() time.Time { return now }
	return s
}

func TestAuthService_Login(t *testing.T) {
	t.Run("senha certa abre sessão e devolve access token", func(t *testing.T) {
		user := userWithPassword(t, "panelas-de-barro")
		var session *model.Session
		var refresh *model.RefreshToken
		users := &mocks.UserRepoMock{FindByEmailFn: func(_ context.Context, email string) (*model.User, error) {
			assert.Equal(t, "maria@email.com", email)
			return user, nil
		}}
		sessions := &mocks.SessionRepoMock{
			DeleteExpiredFn: func(_ context.Context, id string, at time.Time) error {
				assert.Equal(t, userID, id)
				return nil
			},
			CreateFn: func(_ context.Context, s *model.Session, rt *model.RefreshToken) error {
				session, refresh = s, rt
				return nil
			},
		}

		res, err := newAuth(users, sessions, discardLog()).Login(ctx, service.LoginInput{Email: "maria@email.com", Password: "panelas-de-barro"})

		require.NoError(t, err)
		assert.Equal(t, "Bearer", res.TokenType)
		assert.Equal(t, int(service.AccessTokenTTL.Seconds()), res.ExpiresIn)
		assert.Equal(t, userID, res.User.ID)

		claims, err := token.Parse(res.AccessToken, secretKey)
		require.NoError(t, err)
		assert.Equal(t, userID, claims.Subject)
		assert.False(t, claims.IsAdm)

		require.NotNil(t, session)
		assert.Equal(t, userID, session.UserID)
		assert.Equal(t, now.Add(service.SessionMaxAge), session.ExpiresAt)
		require.NotNil(t, refresh)
		assert.Equal(t, now.Add(service.RefreshTokenTTL), refresh.ExpiresAt)
		assert.Equal(t, refresh.ExpiresAt, res.RefreshExpiresAt)
		// no banco vai só o hash do token
		assert.NotEmpty(t, res.RefreshToken)
		assert.Equal(t, sha(res.RefreshToken), refresh.TokenHash)
	})

	t.Run("senha errada é 401 e não abre sessão", func(t *testing.T) {
		user := userWithPassword(t, "panelas-de-barro")
		users := &mocks.UserRepoMock{FindByEmailFn: func(context.Context, string) (*model.User, error) { return user, nil }}

		_, err := newAuth(users, &mocks.SessionRepoMock{}, discardLog()).Login(ctx, service.LoginInput{Email: "maria@email.com", Password: "outra-senha"})

		assert.ErrorIs(t, err, service.ErrInvalidCredentials)
	})

	t.Run("e-mail sem conta responde igual a senha errada", func(t *testing.T) {
		users := &mocks.UserRepoMock{FindByEmailFn: func(context.Context, string) (*model.User, error) { return nil, gorm.ErrRecordNotFound }}

		_, err := newAuth(users, &mocks.SessionRepoMock{}, discardLog()).Login(ctx, service.LoginInput{Email: "ninguem@email.com", Password: "qualquer-senha"})

		assert.ErrorIs(t, err, service.ErrInvalidCredentials)
	})

	t.Run("erro do banco não vira 401", func(t *testing.T) {
		users := &mocks.UserRepoMock{FindByEmailFn: func(context.Context, string) (*model.User, error) { return nil, errors.New("conexão recusada") }}

		_, err := newAuth(users, &mocks.SessionRepoMock{}, discardLog()).Login(ctx, service.LoginInput{Email: "maria@email.com", Password: "panelas-de-barro"})

		require.Error(t, err)
		assert.NotErrorIs(t, err, service.ErrInvalidCredentials)
	})
}

func TestAuthService_LoginLockout(t *testing.T) {
	user := userWithPassword(t, "panelas-de-barro")
	sessions := &mocks.SessionRepoMock{
		DeleteExpiredFn: func(context.Context, string, time.Time) error { return nil },
		CreateFn:        func(context.Context, *model.Session, *model.RefreshToken) error { return nil },
	}
	users := &mocks.UserRepoMock{FindByEmailFn: func(_ context.Context, email string) (*model.User, error) {
		if email != user.Email {
			return nil, gorm.ErrRecordNotFound
		}
		return user, nil
	}}
	auth := newAuth(users, sessions, discardLog())
	login := func(email, password string) error {
		_, err := auth.Login(ctx, service.LoginInput{Email: email, Password: password})
		return err
	}

	// errar e acertar zera a contagem
	for range service.MaxLoginFailures - 1 {
		require.ErrorIs(t, login(user.Email, "errada"), service.ErrInvalidCredentials)
	}
	require.NoError(t, login(user.Email, "panelas-de-barro"))

	for range service.MaxLoginFailures {
		require.ErrorIs(t, login(user.Email, "errada"), service.ErrInvalidCredentials)
	}
	// bloqueado mesmo com a senha certa
	assert.ErrorIs(t, login(user.Email, "panelas-de-barro"), service.ErrTooManyAttempts)
	// o bloqueio é por e-mail
	assert.ErrorIs(t, login("outra@email.com", "qualquer-senha"), service.ErrInvalidCredentials)
}

func TestAuthService_LoginLogDoesNotExposeEmail(t *testing.T) {
	var buf bytes.Buffer
	log := slog.New(slog.NewTextHandler(&buf, nil))
	users := &mocks.UserRepoMock{FindByEmailFn: func(context.Context, string) (*model.User, error) { return nil, gorm.ErrRecordNotFound }}

	_, _ = newAuth(users, &mocks.SessionRepoMock{}, log).Login(ctx, service.LoginInput{Email: "maria@email.com", Password: "qualquer-senha"})

	assert.Contains(t, buf.String(), "login falhou")
	assert.Contains(t, buf.String(), "email_hash=")
	assert.NotContains(t, buf.String(), "maria@email.com")
}

func TestAuthService_Refresh(t *testing.T) {
	session := model.Session{ID: "sessao-1", UserID: userID}
	user := &model.User{ID: userID, Name: "Maria", Email: "maria@email.com"}

	t.Run("token válido roda e devolve outro", func(t *testing.T) {
		sessions := &mocks.SessionRepoMock{RotateFn: func(_ context.Context, hash string, next *model.RefreshToken, at time.Time, grace time.Duration) (repository.RotateResult, error) {
			assert.Equal(t, sha("token-atual"), hash)
			assert.Equal(t, now, at)
			assert.Equal(t, service.RefreshGrace, grace)
			assert.Equal(t, now.Add(service.RefreshTokenTTL), next.ExpiresAt)
			return repository.RotateResult{Outcome: repository.Rotated, Session: session}, nil
		}}
		users := &mocks.UserRepoMock{FindByIDFn: func(context.Context, string) (*model.User, error) { return user, nil }}

		res, err := newAuth(users, sessions, discardLog()).Refresh(ctx, "token-atual")

		require.NoError(t, err)
		assert.NotEmpty(t, res.AccessToken)
		assert.NotEqual(t, "token-atual", res.RefreshToken)
		assert.Equal(t, userID, res.User.ID)
	})

	for _, outcome := range []repository.RotateOutcome{repository.TokenNotFound, repository.SessionEnded, repository.Reused} {
		t.Run("token recusado vira sessão inválida", func(t *testing.T) {
			sessions := &mocks.SessionRepoMock{RotateFn: func(context.Context, string, *model.RefreshToken, time.Time, time.Duration) (repository.RotateResult, error) {
				return repository.RotateResult{Outcome: outcome, Session: session}, nil
			}}

			_, err := newAuth(&mocks.UserRepoMock{}, sessions, discardLog()).Refresh(ctx, "token-velho")

			assert.ErrorIs(t, err, service.ErrSessionInvalid)
		})
	}

	t.Run("reuso fica registrado no log", func(t *testing.T) {
		var buf bytes.Buffer
		sessions := &mocks.SessionRepoMock{RotateFn: func(context.Context, string, *model.RefreshToken, time.Time, time.Duration) (repository.RotateResult, error) {
			return repository.RotateResult{Outcome: repository.Reused, Session: session}, nil
		}}

		_, _ = newAuth(&mocks.UserRepoMock{}, sessions, slog.New(slog.NewTextHandler(&buf, nil))).Refresh(ctx, "token-roubado")

		assert.Contains(t, buf.String(), "refresh token reutilizado")
		assert.Contains(t, buf.String(), "session_id=sessao-1")
	})

	t.Run("sem cookie nem consulta o banco", func(t *testing.T) {
		_, err := newAuth(&mocks.UserRepoMock{}, &mocks.SessionRepoMock{}, discardLog()).Refresh(ctx, "")

		assert.ErrorIs(t, err, service.ErrSessionInvalid)
	})

	t.Run("usuário apagado encerra a sessão", func(t *testing.T) {
		sessions := &mocks.SessionRepoMock{RotateFn: func(context.Context, string, *model.RefreshToken, time.Time, time.Duration) (repository.RotateResult, error) {
			return repository.RotateResult{Outcome: repository.Rotated, Session: session}, nil
		}}
		users := &mocks.UserRepoMock{FindByIDFn: func(context.Context, string) (*model.User, error) { return nil, gorm.ErrRecordNotFound }}

		_, err := newAuth(users, sessions, discardLog()).Refresh(ctx, "token-atual")

		assert.ErrorIs(t, err, service.ErrSessionInvalid)
	})
}

func TestAuthService_Logout(t *testing.T) {
	t.Run("revoga a sessão pelo hash do token", func(t *testing.T) {
		var revoked string
		sessions := &mocks.SessionRepoMock{RevokeByTokenFn: func(_ context.Context, hash string, at time.Time) error {
			revoked = hash
			return nil
		}}

		require.NoError(t, newAuth(&mocks.UserRepoMock{}, sessions, discardLog()).Logout(ctx, "token-atual"))

		assert.Equal(t, sha("token-atual"), revoked)
	})

	t.Run("sem cookie não faz nada", func(t *testing.T) {
		assert.NoError(t, newAuth(&mocks.UserRepoMock{}, &mocks.SessionRepoMock{}, discardLog()).Logout(ctx, ""))
	})
}
