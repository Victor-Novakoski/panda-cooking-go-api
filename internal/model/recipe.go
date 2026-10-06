package model

import "time"

type Recipe struct {
	ID          string    `gorm:"primaryKey;type:uuid;default:gen_random_uuid()"`
	Name        string    `gorm:"not null"`
	Description string    `gorm:"not null"`
	Time        string    `gorm:"not null"`
	Portions    int       `gorm:"not null"`
	UserID      string    `gorm:"type:uuid;not null"`
	CategoryID  uint      `gorm:"not null"`
	CreatedAt   time.Time
	UpdatedAt   time.Time

	User         User              `gorm:"foreignKey:UserID"`
	Category     Category          `gorm:"foreignKey:CategoryID"`
	Images       []ImageRecipe     `gorm:"foreignKey:RecipeID;constraint:OnDelete:CASCADE"`
	Ingredients  []IngredientRecipe `gorm:"foreignKey:RecipeID;constraint:OnDelete:CASCADE"`
	Preparations []Preparation     `gorm:"foreignKey:RecipeID;constraint:OnDelete:CASCADE"`
	Comments     []Comment         `gorm:"foreignKey:RecipeID;constraint:OnDelete:CASCADE"`
	Favorites    []FavoriteRecipe  `gorm:"foreignKey:RecipeID;constraint:OnDelete:CASCADE"`
}
