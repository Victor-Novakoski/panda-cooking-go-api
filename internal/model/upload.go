package model

import "time"

// Upload é uma foto enviada do computador, servida em /api/uploads/{id}.
type Upload struct {
	ID          string `gorm:"primaryKey;default:gen_random_uuid()"`
	UserID      string
	ContentType string
	Data        []byte
	CreatedAt   time.Time
}
