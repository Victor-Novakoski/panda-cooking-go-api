package repository

import "panda-cooking-go-api/internal/model"

type UserRepo interface {
	Create(user *model.User) error
	FindByEmail(email string) (*model.User, error)
	FindByID(id string) (*model.User, error)
	Update(user *model.User) error
	Delete(id string) error
	FindFavoriteRecipes(userID string) ([]model.FavoriteRecipe, error)
}

type RecipeRepo interface {
	Create(recipe *model.Recipe) error
	FindAll() ([]model.Recipe, error)
	FindByID(id string) (*model.Recipe, error)
	Update(recipe *model.Recipe) error
	Replace(recipe *model.Recipe) error
	Delete(id string) error
	FindOrCreateIngredient(name string) (*model.Ingredient, error)
	AddIngredient(ir *model.IngredientRecipe) error
	FindIngredientRecipeByID(id uint) (*model.IngredientRecipe, error)
	DeleteIngredientRecipe(id uint) error
	AddImage(image *model.ImageRecipe) error
	FindImageByID(id uint) (*model.ImageRecipe, error)
	UpdateImage(image *model.ImageRecipe) error
	DeleteImage(id uint) error
	AddPreparation(p *model.Preparation) error
	FindPreparationByID(id uint) (*model.Preparation, error)
	UpdatePreparation(p *model.Preparation) error
	DeletePreparation(id uint) error
}

type CategoryRepo interface {
	FindAll() ([]model.Category, error)
}

type CommentRepo interface {
	Create(comment *model.Comment) error
	FindAll() ([]model.Comment, error)
	FindByRecipe(recipeID string) ([]model.Comment, error)
	FindByID(id uint) (*model.Comment, error)
	Update(comment *model.Comment) error
	Delete(id uint) error
}

type FavoriteRepo interface {
	Find(userID, recipeID string) (*model.FavoriteRecipe, error)
	Create(fav *model.FavoriteRecipe) error
	Delete(id uint) error
}
