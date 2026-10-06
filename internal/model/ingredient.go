package model

type Ingredient struct {
	ID      uint   `gorm:"primaryKey;autoIncrement"`
	Name    string `gorm:"uniqueIndex;not null"`
}

type IngredientRecipe struct {
	ID           uint   `gorm:"primaryKey;autoIncrement"`
	Amount       string `gorm:"not null"`
	RecipeID     string `gorm:"type:uuid;not null"`
	IngredientID uint   `gorm:"not null"`

	Ingredient Ingredient `gorm:"foreignKey:IngredientID"`
}
