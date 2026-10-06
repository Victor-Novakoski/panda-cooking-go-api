package database

import (
	"log"

	"panda-cooking-go-api/internal/model"

	"gorm.io/gorm"
)

var defaultCategories = []string{
	"Café da manhã",
	"Almoço",
	"Jantar",
	"Sobremesa",
	"Lanche",
	"Bebidas",
	"Vegano",
	"Vegetariano",
	"Fitness",
	"Internacional",
}

func Seed(db *gorm.DB) {
	var count int64
	db.Model(&model.Category{}).Count(&count)
	if count > 0 {
		return // já foi populado, não faz nada
	}

	categories := make([]model.Category, len(defaultCategories))
	for i, name := range defaultCategories {
		categories[i] = model.Category{Name: name}
	}

	if err := db.Create(&categories).Error; err != nil {
		log.Printf("erro ao popular categorias: %v", err)
		return
	}

	log.Printf("%d categorias inseridas", len(categories))
}
