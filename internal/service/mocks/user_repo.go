package mocks

import (
	"panda-cooking-go-api/internal/model"
)

// UserRepoMock implementa repository.UserRepo sem banco de dados.
// Cada campo "Fn" é uma função que o teste injeta para controlar o comportamento.
type UserRepoMock struct {
	CreateFn              func(user *model.User) error
	FindByEmailFn         func(email string) (*model.User, error)
	FindByIDFn            func(id string) (*model.User, error)
	UpdateFn              func(user *model.User) error
	DeleteFn              func(id string) error
	FindFavoriteRecipesFn func(userID string) ([]model.FavoriteRecipe, error)
}

func (m *UserRepoMock) Create(user *model.User) error {
	return m.CreateFn(user)
}

func (m *UserRepoMock) FindByEmail(email string) (*model.User, error) {
	return m.FindByEmailFn(email)
}

func (m *UserRepoMock) FindByID(id string) (*model.User, error) {
	return m.FindByIDFn(id)
}

func (m *UserRepoMock) Update(user *model.User) error {
	return m.UpdateFn(user)
}

func (m *UserRepoMock) Delete(id string) error {
	return m.DeleteFn(id)
}

func (m *UserRepoMock) FindFavoriteRecipes(userID string) ([]model.FavoriteRecipe, error) {
	return m.FindFavoriteRecipesFn(userID)
}
