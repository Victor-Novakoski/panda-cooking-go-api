package service

import (
	"context"
	"errors"

	"panda-cooking-go-api/internal/model"
	"panda-cooking-go-api/internal/repository"

	"gorm.io/gorm"
)

type FavoriteService struct {
	repo    repository.FavoriteRepo
	recipes repository.RecipeRepo
}

func NewFavoriteService(repo repository.FavoriteRepo, recipes repository.RecipeRepo) *FavoriteService {
	return &FavoriteService{repo: repo, recipes: recipes}
}

// FavoriteStatus diz se a receita está nos favoritos de quem pede.
type FavoriteStatus struct {
	Favorite bool `json:"favorite"`
}

func (s *FavoriteService) Status(ctx context.Context, userID, recipeID string) (*FavoriteStatus, error) {
	if err := s.checkRecipe(ctx, recipeID); err != nil {
		return nil, err
	}
	ok, err := s.repo.Exists(ctx, userID, recipeID)
	if err != nil {
		return nil, err
	}
	return &FavoriteStatus{Favorite: ok}, nil
}

// Add favorita a receita. O índice único do banco impede o duplicado, mesmo
// com dois cliques ao mesmo tempo.
func (s *FavoriteService) Add(ctx context.Context, userID, recipeID string) (*FavoriteStatus, error) {
	if err := s.checkRecipe(ctx, recipeID); err != nil {
		return nil, err
	}

	err := s.repo.Create(ctx, &model.FavoriteRecipe{UserID: userID, RecipeID: recipeID})
	switch {
	case errors.Is(err, gorm.ErrDuplicatedKey):
		return nil, ErrAlreadyFavorite
	case errors.Is(err, gorm.ErrForeignKeyViolated):
		return nil, missingRecipeOrUser(s.checkRecipe(ctx, recipeID))
	case err != nil:
		return nil, err
	}
	return &FavoriteStatus{Favorite: true}, nil
}

func (s *FavoriteService) Remove(ctx context.Context, userID, recipeID string) error {
	deleted, err := s.repo.Delete(ctx, userID, recipeID)
	if err != nil {
		return err
	}
	if !deleted {
		return ErrFavoriteNotFound
	}
	return nil
}

// List devolve as receitas favoritas do usuário, da favoritada por último
// para a primeira.
func (s *FavoriteService) List(ctx context.Context, userID string, q PageQuery) (Page[RecipeSummary], error) {
	page := q.repo(DefaultRecipesPerPage)
	recipes, total, err := s.repo.ListRecipes(ctx, userID, page)
	if err != nil {
		return Page[RecipeSummary]{}, err
	}
	return newPage(toRecipeSummaries(recipes), page, total), nil
}

func (s *FavoriteService) checkRecipe(ctx context.Context, recipeID string) error {
	ok, err := s.recipes.Exists(ctx, recipeID)
	if err != nil {
		return err
	}
	if !ok {
		return ErrRecipeNotFound
	}
	return nil
}

// missingRecipeOrUser decide o erro de uma chave estrangeira violada ao gravar
// algo de um usuário numa receita: se a receita ainda existe (recipeErr nil),
// quem sumiu foi a conta.
func missingRecipeOrUser(recipeErr error) error {
	if recipeErr == nil {
		return ErrUserNotFound
	}
	return recipeErr
}
