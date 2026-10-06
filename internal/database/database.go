package database

import (
	"panda-cooking-go-api/internal/config"
	"panda-cooking-go-api/internal/model"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func Connect(cfg config.DBConfig) (*gorm.DB, error) {
	db, err := gorm.Open(postgres.Open(cfg.DSN()), &gorm.Config{
		// converte erros do Postgres (ex.: chave duplicada) nos erros do GORM
		TranslateError: true,
	})
	if err != nil {
		return nil, err
	}

	err = db.AutoMigrate(
		&model.User{},
		&model.Category{},
		&model.Recipe{},
		&model.Ingredient{},
		&model.IngredientRecipe{},
		&model.ImageRecipe{},
		&model.Preparation{},
		&model.Comment{},
		&model.FavoriteRecipe{},
	)
	if err != nil {
		return nil, err
	}

	return db, nil
}
