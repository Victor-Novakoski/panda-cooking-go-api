package mocks

import (
	"context"

	"panda-cooking-go-api/internal/model"
	"panda-cooking-go-api/internal/repository"
)

// FavoriteRepoMock implementa repository.FavoriteRepo. Cada teste preenche só as funções
// que espera ver chamadas; uma chamada inesperada derruba o teste com panic.
type FavoriteRepoMock struct {
	ExistsFn      func(ctx context.Context, userID, recipeID string) (bool, error)
	CreateFn      func(ctx context.Context, fav *model.FavoriteRecipe) error
	DeleteFn      func(ctx context.Context, userID, recipeID string) (bool, error)
	ListRecipesFn func(ctx context.Context, userID string, page repository.Page) ([]model.Recipe, int64, error)
}

func (m *FavoriteRepoMock) Exists(ctx context.Context, userID, recipeID string) (bool, error) {
	if m.ExistsFn == nil {
		panic("mocks: FavoriteRepo.Exists chamado sem ExistsFn")
	}
	return m.ExistsFn(ctx, userID, recipeID)
}

func (m *FavoriteRepoMock) Create(ctx context.Context, fav *model.FavoriteRecipe) error {
	if m.CreateFn == nil {
		panic("mocks: FavoriteRepo.Create chamado sem CreateFn")
	}
	return m.CreateFn(ctx, fav)
}

func (m *FavoriteRepoMock) Delete(ctx context.Context, userID, recipeID string) (bool, error) {
	if m.DeleteFn == nil {
		panic("mocks: FavoriteRepo.Delete chamado sem DeleteFn")
	}
	return m.DeleteFn(ctx, userID, recipeID)
}

func (m *FavoriteRepoMock) ListRecipes(ctx context.Context, userID string, page repository.Page) ([]model.Recipe, int64, error) {
	if m.ListRecipesFn == nil {
		panic("mocks: FavoriteRepo.ListRecipes chamado sem ListRecipesFn")
	}
	return m.ListRecipesFn(ctx, userID, page)
}
