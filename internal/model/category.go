package model

type Category struct {
	ID      uint   `gorm:"primaryKey;autoIncrement"`
	Name    string `gorm:"uniqueIndex;not null"`
	Recipes []Recipe `gorm:"foreignKey:CategoryID"`
}
