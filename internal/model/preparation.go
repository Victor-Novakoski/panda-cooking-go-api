package model

type Preparation struct {
	ID          uint   `gorm:"primaryKey;autoIncrement"`
	Description string `gorm:"not null"`
	RecipeID    string `gorm:"type:uuid;not null"`
}
