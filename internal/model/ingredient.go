package model

type Ingredient struct {
	ID   uint `gorm:"primaryKey"`
	Name string
}

type IngredientRecipe struct {
	ID           uint `gorm:"primaryKey"`
	Amount       string
	RecipeID     string
	IngredientID uint

	Ingredient Ingredient `gorm:"foreignKey:IngredientID"`
}
