package model

import "time"

// Session é um login. Dura até ExpiresAt (limite absoluto) ou até ser
// revogada no logout ou por reuso de refresh token.
type Session struct {
	ID        string `gorm:"primaryKey;default:gen_random_uuid()"`
	UserID    string
	CreatedAt time.Time
	ExpiresAt time.Time
	RevokedAt *time.Time
}

// RefreshToken guarda só o hash do token entregue ao navegador.
type RefreshToken struct {
	ID        uint `gorm:"primaryKey"`
	SessionID string
	TokenHash string
	ExpiresAt time.Time
	UsedAt    *time.Time
	CreatedAt time.Time
}
