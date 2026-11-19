package mocks

import (
	"context"

	"panda-cooking-go-api/internal/model"
)

// UploadRepoMock implementa repository.UploadRepo. Cada teste preenche só as funções
// que espera ver chamadas; uma chamada inesperada derruba o teste com panic.
type UploadRepoMock struct {
	CreateFn   func(ctx context.Context, upload *model.Upload) error
	FindByIDFn func(ctx context.Context, id string) (*model.Upload, error)
}

func (m *UploadRepoMock) Create(ctx context.Context, upload *model.Upload) error {
	if m.CreateFn == nil {
		panic("mocks: UploadRepo.Create chamado sem CreateFn")
	}
	return m.CreateFn(ctx, upload)
}

func (m *UploadRepoMock) FindByID(ctx context.Context, id string) (*model.Upload, error) {
	if m.FindByIDFn == nil {
		panic("mocks: UploadRepo.FindByID chamado sem FindByIDFn")
	}
	return m.FindByIDFn(ctx, id)
}
