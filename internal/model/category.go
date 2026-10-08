package model

type Category struct {
	ID   uint `gorm:"primaryKey"`
	Name string
}
