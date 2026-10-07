package model

type Preparation struct {
	ID          uint `gorm:"primaryKey"`
	Description string
	RecipeID    string
}
