package mocks

import "panda-cooking-go-api/internal/model"

type CommentRepoMock struct {
	CreateFn       func(comment *model.Comment) error
	FindAllFn      func() ([]model.Comment, error)
	FindByRecipeFn func(recipeID string) ([]model.Comment, error)
	FindByIDFn     func(id uint) (*model.Comment, error)
	UpdateFn       func(comment *model.Comment) error
	DeleteFn       func(id uint) error
}

func (m *CommentRepoMock) Create(comment *model.Comment) error { return m.CreateFn(comment) }
func (m *CommentRepoMock) FindAll() ([]model.Comment, error)   { return m.FindAllFn() }
func (m *CommentRepoMock) FindByRecipe(recipeID string) ([]model.Comment, error) {
	return m.FindByRecipeFn(recipeID)
}
func (m *CommentRepoMock) FindByID(id uint) (*model.Comment, error) {
	return m.FindByIDFn(id)
}
func (m *CommentRepoMock) Update(comment *model.Comment) error { return m.UpdateFn(comment) }
func (m *CommentRepoMock) Delete(id uint) error                { return m.DeleteFn(id) }
