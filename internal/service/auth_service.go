package service

import (
	"context"
	"crypto/hmac"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"errors"
	"fmt"
	"log/slog"
	"sync"
	"time"

	"panda-cooking-go-api/internal/model"
	"panda-cooking-go-api/internal/ratelimit"
	"panda-cooking-go-api/internal/repository"
	"panda-cooking-go-api/pkg/token"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

const (
	// Depois de MaxLoginFailures senhas erradas para o mesmo e-mail, o login
	// dele fica bloqueado até a janela de LoginLockout acabar.
	MaxLoginFailures = 5
	LoginLockout     = 15 * time.Minute

	// AccessTokenTTL é a validade do JWT. Curta porque um JWT não tem como
	// ser revogado; quem segura a sessão é o refresh token.
	AccessTokenTTL = 15 * time.Minute
	// RefreshTokenTTL é a validade de cada refresh token (ele roda a cada uso).
	RefreshTokenTTL = 7 * 24 * time.Hour
	// SessionMaxAge é o limite da sessão: depois disso é preciso entrar de novo.
	SessionMaxAge = 30 * 24 * time.Hour
	// RefreshGrace é quanto tempo um refresh token já usado ainda é aceito
	// (duas abas renovando juntas, resposta que não chegou ao navegador).
	RefreshGrace = 20 * time.Second
)

// AuthResult é o que login e renovação devolvem. O refresh token não vai no
// JSON: o handler põe em cookie HttpOnly.
type AuthResult struct {
	AccessToken      string       `json:"access_token"`
	TokenType        string       `json:"token_type"`
	ExpiresIn        int          `json:"expires_in"`
	User             UserResponse `json:"user"`
	RefreshToken     string       `json:"-"`
	RefreshExpiresAt time.Time    `json:"-"`
}

type AuthService struct {
	users         repository.UserRepo
	sessions      repository.SessionRepo
	secretKey     string
	loginFailures *ratelimit.Limiter
	log           *slog.Logger

	// Now pode ser trocado nos testes para controlar o relógio.
	Now func() time.Time
}

func NewAuthService(users repository.UserRepo, sessions repository.SessionRepo, secretKey string, log *slog.Logger) *AuthService {
	return &AuthService{
		users:         users,
		sessions:      sessions,
		secretKey:     secretKey,
		loginFailures: ratelimit.New(MaxLoginFailures, LoginLockout),
		log:           log,
		Now:           time.Now,
	}
}

// Login confere e-mail e senha e abre uma sessão.
//
// E-mail inexistente responde igual a senha errada, conta tentativa para o
// bloqueio e também passa pelo bcrypt, para o tempo de resposta não revelar
// quem tem conta.
func (s *AuthService) Login(ctx context.Context, input LoginInput) (*AuthResult, error) {
	email := input.Email
	// A tentativa conta antes do bcrypt: conferir e contar depois deixaria
	// passar várias tentativas simultâneas com o mesmo e-mail. O login certo
	// zera o contador.
	if ok, wait := s.loginFailures.Allow(email); !ok {
		s.log.WarnContext(ctx, "login bloqueado por excesso de tentativas", "email_hash", s.emailHash(email))
		blocked := *ErrTooManyAttempts
		blocked.RetryAfter = wait
		return nil, &blocked
	}

	user, err := s.users.FindByEmail(ctx, email)
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}

	hash := dummyHash()
	if user != nil {
		hash = []byte(user.PasswordHash)
	}
	if bcrypt.CompareHashAndPassword(hash, []byte(input.Password)) != nil || user == nil {
		s.log.WarnContext(ctx, "login falhou", "email_hash", s.emailHash(email), "conta_existe", user != nil)
		return nil, ErrInvalidCredentials
	}
	s.loginFailures.Reset(email)

	now := s.Now()
	// aproveita o login para limpar as sessões velhas do usuário
	if err := s.sessions.DeleteExpired(ctx, user.ID, now); err != nil {
		s.log.WarnContext(ctx, "não foi possível apagar sessões vencidas", "user_id", user.ID, "erro", err)
	}

	refresh, refreshHash, err := newRefreshToken()
	if err != nil {
		return nil, err
	}
	session := &model.Session{UserID: user.ID, ExpiresAt: now.Add(SessionMaxAge)}
	rt := &model.RefreshToken{TokenHash: refreshHash, ExpiresAt: now.Add(RefreshTokenTTL)}
	if err := s.sessions.Create(ctx, session, rt); err != nil {
		return nil, fmt.Errorf("criar sessão: %w", err)
	}

	return s.result(user, refresh, rt.ExpiresAt, now)
}

// Refresh troca o refresh token por um novo e devolve um access token novo.
// Token já usado fora da janela de tolerância derruba a sessão inteira:
// alguém além do dono pode ter uma cópia dele.
func (s *AuthService) Refresh(ctx context.Context, refreshToken string) (*AuthResult, error) {
	if refreshToken == "" {
		return nil, ErrSessionInvalid
	}

	next, nextHash, err := newRefreshToken()
	if err != nil {
		return nil, err
	}
	now := s.Now()
	rt := &model.RefreshToken{TokenHash: nextHash, ExpiresAt: now.Add(RefreshTokenTTL)}

	res, err := s.sessions.Rotate(ctx, hashToken(refreshToken), rt, now, RefreshGrace)
	if err != nil {
		return nil, fmt.Errorf("renovar sessão: %w", err)
	}

	switch res.Outcome {
	case repository.Rotated:
	case repository.Reused:
		s.log.WarnContext(ctx, "refresh token reutilizado, sessão revogada",
			"user_id", res.Session.UserID, "session_id", res.Session.ID)
		return nil, ErrSessionInvalid
	default:
		return nil, ErrSessionInvalid
	}

	user, err := s.users.FindByID(ctx, res.Session.UserID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrSessionInvalid
		}
		return nil, err
	}

	return s.result(user, next, rt.ExpiresAt, now)
}

// Logout encerra a sessão do refresh token. Token desconhecido não é erro:
// o objetivo (não ter sessão) já está cumprido.
func (s *AuthService) Logout(ctx context.Context, refreshToken string) error {
	if refreshToken == "" {
		return nil
	}
	return s.sessions.RevokeByToken(ctx, hashToken(refreshToken), s.Now())
}

func (s *AuthService) result(user *model.User, refresh string, refreshExpiresAt, now time.Time) (*AuthResult, error) {
	access, err := token.Generate(user.ID, user.IsAdm, s.secretKey, now, AccessTokenTTL)
	if err != nil {
		return nil, fmt.Errorf("gerar access token: %w", err)
	}
	return &AuthResult{
		AccessToken:      access,
		TokenType:        "Bearer",
		ExpiresIn:        int(AccessTokenTTL.Seconds()),
		User:             toUserResponse(user),
		RefreshToken:     refresh,
		RefreshExpiresAt: refreshExpiresAt,
	}, nil
}

// emailHash identifica o e-mail no log sem gravá-lo em claro. É HMAC com a
// SECRET_KEY: sem a chave não dá para testar e-mails conhecidos contra o log.
func (s *AuthService) emailHash(email string) string {
	mac := hmac.New(sha256.New, []byte(s.secretKey))
	mac.Write([]byte(email))
	return hex.EncodeToString(mac.Sum(nil))[:16]
}

// newRefreshToken gera 32 bytes aleatórios (o valor do cookie) e o hash
// que vai para o banco.
func newRefreshToken() (plain, hash string, err error) {
	b := make([]byte, 32)
	if _, err := rand.Read(b); err != nil {
		return "", "", fmt.Errorf("gerar refresh token: %w", err)
	}
	plain = base64.RawURLEncoding.EncodeToString(b)
	return plain, hashToken(plain), nil
}

// hashToken usa SHA-256 puro: o token já é aleatório e longo, não precisa
// de bcrypt (que existe para senhas fracas).
func hashToken(t string) string {
	sum := sha256.Sum256([]byte(t))
	return hex.EncodeToString(sum[:])
}

// dummyHash é comparado quando o e-mail não existe, para o login levar o
// mesmo tempo nos dois casos. Gerado uma vez, com o mesmo custo das senhas.
var dummyHash = sync.OnceValue(func() []byte {
	h, err := bcrypt.GenerateFromPassword([]byte("senha-que-ninguem-tem"), bcrypt.DefaultCost)
	if err != nil {
		panic(err)
	}
	return h
})
