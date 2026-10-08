package mocks

import (
	"context"

	"panda-cooking-go-api/internal/model"
)

// CategoryRepoMock implementa repository.CategoryRepo. Cada teste preenche só as funções
// que espera ver chamadas; uma chamada inesperada derruba o teste com panic.
type CategoryRepoMock struct {
	FindAllFn func(ctx context.Context) ([]model.Category, error)
	ExistsFn  func(ctx context.Context, id uint) (bool, error)
}

func (m *CategoryRepoMock) FindAll(ctx context.Context) ([]model.Category, error) {
	if m.FindAllFn == nil {
		panic("mocks: CategoryRepo.FindAll chamado sem FindAllFn")
	}
	return m.FindAllFn(ctx)
}

func (m *CategoryRepoMock) Exists(ctx context.Context, id uint) (bool, error) {
	if m.ExistsFn == nil {
		panic("mocks: CategoryRepo.Exists chamado sem ExistsFn")
	}
	return m.ExistsFn(ctx, id)
}
