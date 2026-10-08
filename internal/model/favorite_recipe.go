package model

import "time"

type FavoriteRecipe struct {
	ID        uint `gorm:"primaryKey"`
	UserID    string
	RecipeID  string
	CreatedAt time.Time
}
