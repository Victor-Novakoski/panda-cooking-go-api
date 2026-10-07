package mocks

import (
	"context"
	"time"

	"panda-cooking-go-api/internal/model"
	"panda-cooking-go-api/internal/repository"
)

// SessionRepoMock implementa repository.SessionRepo. Cada teste preenche só as funções
// que espera ver chamadas; uma chamada inesperada derruba o teste com panic.
type SessionRepoMock struct {
	CreateFn        func(ctx context.Context, session *model.Session, token *model.RefreshToken) error
	RotateFn        func(ctx context.Context, tokenHash string, next *model.RefreshToken, now time.Time, grace time.Duration) (repository.RotateResult, error)
	RevokeByTokenFn func(ctx context.Context, tokenHash string, now time.Time) error
	DeleteExpiredFn func(ctx context.Context, userID string, now time.Time) error
}

func (m *SessionRepoMock) Create(ctx context.Context, session *model.Session, token *model.RefreshToken) error {
	if m.CreateFn == nil {
		panic("mocks: SessionRepo.Create chamado sem CreateFn")
	}
	return m.CreateFn(ctx, session, token)
}

func (m *SessionRepoMock) Rotate(ctx context.Context, tokenHash string, next *model.RefreshToken, now time.Time, grace time.Duration) (repository.RotateResult, error) {
	if m.RotateFn == nil {
		panic("mocks: SessionRepo.Rotate chamado sem RotateFn")
	}
	return m.RotateFn(ctx, tokenHash, next, now, grace)
}

func (m *SessionRepoMock) RevokeByToken(ctx context.Context, tokenHash string, now time.Time) error {
	if m.RevokeByTokenFn == nil {
		panic("mocks: SessionRepo.RevokeByToken chamado sem RevokeByTokenFn")
	}
	return m.RevokeByTokenFn(ctx, tokenHash, now)
}

func (m *SessionRepoMock) DeleteExpired(ctx context.Context, userID string, now time.Time) error {
	if m.DeleteExpiredFn == nil {
		panic("mocks: SessionRepo.DeleteExpired chamado sem DeleteExpiredFn")
	}
	return m.DeleteExpiredFn(ctx, userID, now)
}
