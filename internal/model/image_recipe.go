package model

type ImageRecipe struct {
	ID       uint   `gorm:"primaryKey;autoIncrement"`
	URL      string `gorm:"not null"`
	RecipeID string `gorm:"type:uuid;not null"`
}
