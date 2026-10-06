package model

import "time"

type Comment struct {
	ID          uint      `gorm:"primaryKey;autoIncrement"`
	Description string    `gorm:"not null"`
	UserID      string    `gorm:"type:uuid;not null"`
	RecipeID    string    `gorm:"type:uuid;not null"`
	CreatedAt   time.Time
	UpdatedAt   time.Time

	User User `gorm:"foreignKey:UserID"`
}
