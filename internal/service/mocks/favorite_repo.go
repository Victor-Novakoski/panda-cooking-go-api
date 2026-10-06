package mocks

import "panda-cooking-go-api/internal/model"

type FavoriteRepoMock struct {
	FindFn   func(userID, recipeID string) (*model.FavoriteRecipe, error)
	CreateFn func(fav *model.FavoriteRecipe) error
	DeleteFn func(id uint) error
}

func (m *FavoriteRepoMock) Find(userID, recipeID string) (*model.FavoriteRecipe, error) {
	return m.FindFn(userID, recipeID)
}
func (m *FavoriteRepoMock) Create(fav *model.FavoriteRecipe) error { return m.CreateFn(fav) }
func (m *FavoriteRepoMock) Delete(id uint) error                   { return m.DeleteFn(id) }
