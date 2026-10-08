package model

import "time"

type User struct {
	ID           string `gorm:"primaryKey;default:gen_random_uuid()"`
	Name         string
	Email        string
	PasswordHash string
	ImageProfile string
	IsAdm        bool
	CreatedAt    time.Time
	UpdatedAt    time.Time
}
