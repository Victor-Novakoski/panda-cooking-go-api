package repository

import (
	"context"
	"time"

	"panda-cooking-go-api/internal/model"
)

// Os services dependem destas interfaces, não do GORM: nos testes de
// service entram os mocks de internal/service/mocks.

type UserRepo interface {
	Create(ctx context.Context, user *model.User) error
	FindByEmail(ctx context.Context, email string) (*model.User, error)
	FindByID(ctx context.Context, id string) (*model.User, error)
	Update(ctx context.Context, user *model.User) error
	Delete(ctx context.Context, id string) error
}

type RecipeRepo interface {
	// Create grava a receita com fotos, ingredientes e passos numa transação.
	Create(ctx context.Context, recipe *model.Recipe) error
	List(ctx context.Context, filter RecipeFilter, page Page) ([]model.Recipe, int64, error)
	FindByID(ctx context.Context, id string) (*model.Recipe, error)
	Exists(ctx context.Context, id string) (bool, error)
	Update(ctx context.Context, recipe *model.Recipe) error
	Replace(ctx context.Context, recipe *model.Recipe) error
	Delete(ctx context.Context, id string) error
	AddImage(ctx context.Context, image *model.ImageRecipe) error
	FindImageByID(ctx context.Context, id uint) (*model.ImageRecipe, error)
	UpdateImage(ctx context.Context, image *model.ImageRecipe) error
	DeleteImage(ctx context.Context, id uint) error
	// AddIngredient acha ou cria o ingrediente pelo nome e liga à receita.
	AddIngredient(ctx context.Context, ir *model.IngredientRecipe) error
	FindIngredientRecipeByID(ctx context.Context, id uint) (*model.IngredientRecipe, error)
	DeleteIngredientRecipe(ctx context.Context, id uint) error
	AddPreparation(ctx context.Context, p *model.Preparation) error
	FindPreparationByID(ctx context.Context, id uint) (*model.Preparation, error)
	UpdatePreparation(ctx context.Context, p *model.Preparation) error
	DeletePreparation(ctx context.Context, id uint) error
}

type CategoryRepo interface {
	FindAll(ctx context.Context) ([]model.Category, error)
	Exists(ctx context.Context, id uint) (bool, error)
}

type CommentRepo interface {
	Create(ctx context.Context, comment *model.Comment) error
	ListByRecipe(ctx context.Context, recipeID string, page Page) ([]model.Comment, int64, error)
	FindByID(ctx context.Context, id uint) (*model.Comment, error)
	Update(ctx context.Context, comment *model.Comment) error
	Delete(ctx context.Context, id uint) error
}

type FavoriteRepo interface {
	Exists(ctx context.Context, userID, recipeID string) (bool, error)
	Create(ctx context.Context, fav *model.FavoriteRecipe) error
	// Delete diz se havia o favorito para apagar.
	Delete(ctx context.Context, userID, recipeID string) (bool, error)
	ListRecipes(ctx context.Context, userID string, page Page) ([]model.Recipe, int64, error)
}

type SessionRepo interface {
	// Create grava a sessão e o primeiro refresh token dela.
	Create(ctx context.Context, session *model.Session, token *model.RefreshToken) error
	// Rotate troca o refresh token pelo próximo (ver RotateResult).
	Rotate(ctx context.Context, tokenHash string, next *model.RefreshToken, now time.Time, grace time.Duration) (RotateResult, error)
	// RevokeByToken encerra a sessão do token (logout).
	RevokeByToken(ctx context.Context, tokenHash string, now time.Time) error
	// DeleteExpired apaga as sessões vencidas ou revogadas do usuário.
	DeleteExpired(ctx context.Context, userID string, now time.Time) error
}
