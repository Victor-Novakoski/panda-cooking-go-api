package repository

import (
	"context"

	"panda-cooking-go-api/internal/model"

	"gorm.io/gorm"
)

type FavoriteRepository struct {
	db *gorm.DB
}

func NewFavoriteRepository(db *gorm.DB) *FavoriteRepository {
	return &FavoriteRepository{db: db}
}

func (r *FavoriteRepository) Exists(ctx context.Context, userID, recipeID string) (bool, error) {
	if !IsUUID(recipeID) {
		return false, nil
	}
	var count int64
	err := r.db.WithContext(ctx).Model(&model.FavoriteRecipe{}).
		Where("user_id = ? AND recipe_id = ?", userID, recipeID).
		Count(&count).Error
	return count > 0, err
}

// Create grava o favorito. Duplicado vira gorm.ErrDuplicatedKey (índice único).
func (r *FavoriteRepository) Create(ctx context.Context, fav *model.FavoriteRecipe) error {
	return r.db.WithContext(ctx).Create(fav).Error
}

func (r *FavoriteRepository) Delete(ctx context.Context, userID, recipeID string) (bool, error) {
	if !IsUUID(recipeID) {
		return false, nil
	}
	res := r.db.WithContext(ctx).
		Where("user_id = ? AND recipe_id = ?", userID, recipeID).
		Delete(&model.FavoriteRecipe{})
	return res.RowsAffected > 0, res.Error
}

// ListRecipes lista as receitas favoritas do usuário, da favoritada por
// último para a primeira.
func (r *FavoriteRepository) ListRecipes(ctx context.Context, userID string, page Page) ([]model.Recipe, int64, error) {
	db := r.db.WithContext(ctx).Model(&model.Recipe{}).
		Joins("JOIN favorite_recipes f ON f.recipe_id = recipes.id AND f.user_id = ?", userID).
		Session(&gorm.Session{})

	var total int64
	if err := db.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	var recipes []model.Recipe
	err := preloadSummary(db).
		Order("f.created_at DESC, f.id DESC").
		Limit(page.Size).Offset(page.offset()).
		Find(&recipes).Error
	return recipes, total, err
}
