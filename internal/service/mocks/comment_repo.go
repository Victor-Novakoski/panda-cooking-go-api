package mocks

import (
	"context"

	"panda-cooking-go-api/internal/model"
	"panda-cooking-go-api/internal/repository"
)

// CommentRepoMock implementa repository.CommentRepo. Cada teste preenche só as funções
// que espera ver chamadas; uma chamada inesperada derruba o teste com panic.
type CommentRepoMock struct {
	CreateFn       func(ctx context.Context, comment *model.Comment) error
	ListByRecipeFn func(ctx context.Context, recipeID string, page repository.Page) ([]model.Comment, int64, error)
	FindByIDFn     func(ctx context.Context, id uint) (*model.Comment, error)
	UpdateFn       func(ctx context.Context, comment *model.Comment) error
	DeleteFn       func(ctx context.Context, id uint) error
}

func (m *CommentRepoMock) Create(ctx context.Context, comment *model.Comment) error {
	if m.CreateFn == nil {
		panic("mocks: CommentRepo.Create chamado sem CreateFn")
	}
	return m.CreateFn(ctx, comment)
}

func (m *CommentRepoMock) ListByRecipe(ctx context.Context, recipeID string, page repository.Page) ([]model.Comment, int64, error) {
	if m.ListByRecipeFn == nil {
		panic("mocks: CommentRepo.ListByRecipe chamado sem ListByRecipeFn")
	}
	return m.ListByRecipeFn(ctx, recipeID, page)
}

func (m *CommentRepoMock) FindByID(ctx context.Context, id uint) (*model.Comment, error) {
	if m.FindByIDFn == nil {
		panic("mocks: CommentRepo.FindByID chamado sem FindByIDFn")
	}
	return m.FindByIDFn(ctx, id)
}

func (m *CommentRepoMock) Update(ctx context.Context, comment *model.Comment) error {
	if m.UpdateFn == nil {
		panic("mocks: CommentRepo.Update chamado sem UpdateFn")
	}
	return m.UpdateFn(ctx, comment)
}

func (m *CommentRepoMock) Delete(ctx context.Context, id uint) error {
	if m.DeleteFn == nil {
		panic("mocks: CommentRepo.Delete chamado sem DeleteFn")
	}
	return m.DeleteFn(ctx, id)
}
