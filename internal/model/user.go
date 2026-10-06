package model

import "time"

type User struct {
	ID           string    `gorm:"primaryKey;type:uuid;default:gen_random_uuid()"`
	Name         string    `gorm:"not null"`
	Email        string    `gorm:"uniqueIndex;not null"`
	Password     string    `gorm:"not null"`
	ImageProfile string
	IsAdm        bool      `gorm:"default:false"`
	CreatedAt    time.Time
	UpdatedAt    time.Time

	Recipes         []Recipe         `gorm:"foreignKey:UserID;constraint:OnDelete:CASCADE"`
	Comments        []Comment        `gorm:"foreignKey:UserID;constraint:OnDelete:CASCADE"`
	FavoriteRecipes []FavoriteRecipe `gorm:"foreignKey:UserID;constraint:OnDelete:CASCADE"`
}
