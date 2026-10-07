package mocks

import (
	"context"

	"panda-cooking-go-api/internal/model"
)

// UserRepoMock implementa repository.UserRepo. Cada teste preenche só as funções
// que espera ver chamadas; uma chamada inesperada derruba o teste com panic.
type UserRepoMock struct {
	CreateFn      func(ctx context.Context, user *model.User) error
	FindByEmailFn func(ctx context.Context, email string) (*model.User, error)
	FindByIDFn    func(ctx context.Context, id string) (*model.User, error)
	UpdateFn      func(ctx context.Context, user *model.User) error
	DeleteFn      func(ctx context.Context, id string) error
}

func (m *UserRepoMock) Create(ctx context.Context, user *model.User) error {
	if m.CreateFn == nil {
		panic("mocks: UserRepo.Create chamado sem CreateFn")
	}
	return m.CreateFn(ctx, user)
}

func (m *UserRepoMock) FindByEmail(ctx context.Context, email string) (*model.User, error) {
	if m.FindByEmailFn == nil {
		panic("mocks: UserRepo.FindByEmail chamado sem FindByEmailFn")
	}
	return m.FindByEmailFn(ctx, email)
}

func (m *UserRepoMock) FindByID(ctx context.Context, id string) (*model.User, error) {
	if m.FindByIDFn == nil {
		panic("mocks: UserRepo.FindByID chamado sem FindByIDFn")
	}
	return m.FindByIDFn(ctx, id)
}

func (m *UserRepoMock) Update(ctx context.Context, user *model.User) error {
	if m.UpdateFn == nil {
		panic("mocks: UserRepo.Update chamado sem UpdateFn")
	}
	return m.UpdateFn(ctx, user)
}

func (m *UserRepoMock) Delete(ctx context.Context, id string) error {
	if m.DeleteFn == nil {
		panic("mocks: UserRepo.Delete chamado sem DeleteFn")
	}
	return m.DeleteFn(ctx, id)
}
