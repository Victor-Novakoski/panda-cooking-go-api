package repository

import (
	"panda-cooking-go-api/internal/model"

	"gorm.io/gorm"
)

type FavoriteRepository struct {
	db *gorm.DB
}

func NewFavoriteRepository(db *gorm.DB) *FavoriteRepository {
	return &FavoriteRepository{db: db}
}

func (r *FavoriteRepository) Find(userID, recipeID string) (*model.FavoriteRecipe, error) {
	var fav model.FavoriteRecipe
	err := r.db.Where("user_id = ? AND recipe_id = ?", userID, recipeID).First(&fav).Error
	if err != nil {
		return nil, err
	}
	return &fav, nil
}

func (r *FavoriteRepository) Create(fav *model.FavoriteRecipe) error {
	return r.db.Create(fav).Error
}

func (r *FavoriteRepository) Delete(id uint) error {
	return r.db.Delete(&model.FavoriteRecipe{}, id).Error
}
