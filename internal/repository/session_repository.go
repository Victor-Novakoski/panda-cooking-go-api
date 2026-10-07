package repository

import (
	"context"
	"errors"
	"time"

	"panda-cooking-go-api/internal/model"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

// RotateResult é o que aconteceu ao usar um refresh token.
type RotateResult struct {
	Outcome RotateOutcome
	Session model.Session
}

type RotateOutcome int

const (
	// Rotated: token válido; o próximo foi gravado na mesma sessão.
	Rotated RotateOutcome = iota + 1
	// TokenNotFound: o token não existe (inventado ou de sessão já apagada).
	TokenNotFound
	// SessionEnded: token ou sessão vencidos, ou sessão revogada.
	SessionEnded
	// Reused: token já usado fora da janela de tolerância. A sessão foi
	// revogada, porque alguém além do dono pode ter o token.
	Reused
)

type SessionRepository struct {
	db *gorm.DB
}

func NewSessionRepository(db *gorm.DB) *SessionRepository {
	return &SessionRepository{db: db}
}

func (r *SessionRepository) Create(ctx context.Context, session *model.Session, token *model.RefreshToken) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(session).Error; err != nil {
			return err
		}
		token.SessionID = session.ID
		return tx.Create(token).Error
	})
}

// Rotate usa o refresh token e grava o próximo, numa transação com a linha
// do token travada (FOR UPDATE): duas renovações com o mesmo token esperam
// uma pela outra em vez de passarem juntas.
//
// Token já usado há menos de grace ainda vale: é o caso de duas abas
// renovando ao mesmo tempo ou de a resposta anterior não ter chegado.
func (r *SessionRepository) Rotate(ctx context.Context, tokenHash string, next *model.RefreshToken, now time.Time, grace time.Duration) (RotateResult, error) {
	var result RotateResult
	err := r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		var token model.RefreshToken
		err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).
			Where("token_hash = ?", tokenHash).First(&token).Error
		if errors.Is(err, gorm.ErrRecordNotFound) {
			result.Outcome = TokenNotFound
			return nil
		}
		if err != nil {
			return err
		}

		if err := tx.First(&result.Session, "id = ?", token.SessionID).Error; err != nil {
			return err
		}
		session := &result.Session

		switch {
		case session.RevokedAt != nil || !now.Before(session.ExpiresAt) || !now.Before(token.ExpiresAt):
			result.Outcome = SessionEnded
			return nil
		case token.UsedAt != nil && now.Sub(*token.UsedAt) > grace:
			result.Outcome = Reused
			return tx.Model(session).Update("revoked_at", now).Error
		}

		if token.UsedAt == nil {
			if err := tx.Model(&token).Update("used_at", now).Error; err != nil {
				return err
			}
		}

		// Tokens vencidos da sessão já não servem nem para detectar reúso
		// (vencido é recusado antes): saem para a tabela não crescer sem fim.
		if err := tx.Where("session_id = ? AND expires_at <= ?", session.ID, now).
			Delete(&model.RefreshToken{}).Error; err != nil {
			return err
		}

		next.SessionID = session.ID
		if next.ExpiresAt.After(session.ExpiresAt) {
			next.ExpiresAt = session.ExpiresAt
		}
		if err := tx.Create(next).Error; err != nil {
			return err
		}
		result.Outcome = Rotated
		return nil
	})
	return result, err
}

func (r *SessionRepository) RevokeByToken(ctx context.Context, tokenHash string, now time.Time) error {
	return r.db.WithContext(ctx).Model(&model.Session{}).
		Where("revoked_at IS NULL AND id = (SELECT session_id FROM refresh_tokens WHERE token_hash = ?)", tokenHash).
		Update("revoked_at", now).Error
}

func (r *SessionRepository) DeleteExpired(ctx context.Context, userID string, now time.Time) error {
	return r.db.WithContext(ctx).
		Where("user_id = ? AND (expires_at <= ? OR revoked_at IS NOT NULL)", userID, now).
		Delete(&model.Session{}).Error
}
