package model

import "time"

type Comment struct {
	ID          uint `gorm:"primaryKey"`
	Description string
	UserID      string
	RecipeID    string
	CreatedAt   time.Time
	UpdatedAt   time.Time

	User User `gorm:"foreignKey:UserID"`
}
