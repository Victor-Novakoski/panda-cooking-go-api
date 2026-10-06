package repository

import (
	"panda-cooking-go-api/internal/model"

	"gorm.io/gorm"
)

type UserRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) *UserRepository {
	return &UserRepository{db: db}
}

func (r *UserRepository) Create(user *model.User) error {
	return r.db.Create(user).Error
}

func (r *UserRepository) FindByEmail(email string) (*model.User, error) {
	var user model.User
	err := r.db.Where("email = ?", email).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *UserRepository) FindByID(id string) (*model.User, error) {
	var user model.User
	err := r.db.First(&user, "id = ?", id).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *UserRepository) Update(user *model.User) error {
	return r.db.Save(user).Error
}

func (r *UserRepository) Delete(id string) error {
	return r.db.Delete(&model.User{}, "id = ?", id).Error
}

func (r *UserRepository) FindFavoriteRecipes(userID string) ([]model.FavoriteRecipe, error) {
	var favorites []model.FavoriteRecipe
	err := r.db.
		Preload("Recipe").
		Preload("Recipe.Category").
		Preload("Recipe.Images").
		Where("user_id = ?", userID).
		Find(&favorites).Error
	return favorites, err
}
