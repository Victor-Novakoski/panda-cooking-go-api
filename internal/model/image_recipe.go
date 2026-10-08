package model

type ImageRecipe struct {
	ID       uint `gorm:"primaryKey"`
	URL      string
	RecipeID string
}
