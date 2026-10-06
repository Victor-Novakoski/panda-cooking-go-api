package model

type FavoriteRecipe struct {
	ID       uint   `gorm:"primaryKey;autoIncrement"`
	UserID   string `gorm:"type:uuid;not null"`
	RecipeID string `gorm:"type:uuid;not null"`

	Recipe Recipe `gorm:"foreignKey:RecipeID"`
}
