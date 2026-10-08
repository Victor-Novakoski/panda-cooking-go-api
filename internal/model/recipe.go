package model

import "time"

type Recipe struct {
	ID          string `gorm:"primaryKey;default:gen_random_uuid()"`
	Name        string
	Description string
	Time        string
	Portions    int
	UserID      string
	CategoryID  uint
	CreatedAt   time.Time
	UpdatedAt   time.Time

	User         User               `gorm:"foreignKey:UserID"`
	Category     Category           `gorm:"foreignKey:CategoryID"`
	Images       []ImageRecipe      `gorm:"foreignKey:RecipeID"`
	Ingredients  []IngredientRecipe `gorm:"foreignKey:RecipeID"`
	Preparations []Preparation      `gorm:"foreignKey:RecipeID"`
}
