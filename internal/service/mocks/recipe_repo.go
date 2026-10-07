package mocks

import (
	"context"

	"panda-cooking-go-api/internal/model"
	"panda-cooking-go-api/internal/repository"
)

// RecipeRepoMock implementa repository.RecipeRepo. Cada teste preenche só as funções
// que espera ver chamadas; uma chamada inesperada derruba o teste com panic.
type RecipeRepoMock struct {
	CreateFn                   func(ctx context.Context, recipe *model.Recipe) error
	ListFn                     func(ctx context.Context, filter repository.RecipeFilter, page repository.Page) ([]model.Recipe, int64, error)
	FindByIDFn                 func(ctx context.Context, id string) (*model.Recipe, error)
	ExistsFn                   func(ctx context.Context, id string) (bool, error)
	UpdateFn                   func(ctx context.Context, recipe *model.Recipe) error
	ReplaceFn                  func(ctx context.Context, recipe *model.Recipe) error
	DeleteFn                   func(ctx context.Context, id string) error
	AddImageFn                 func(ctx context.Context, image *model.ImageRecipe) error
	FindImageByIDFn            func(ctx context.Context, id uint) (*model.ImageRecipe, error)
	UpdateImageFn              func(ctx context.Context, image *model.ImageRecipe) error
	DeleteImageFn              func(ctx context.Context, id uint) error
	AddIngredientFn            func(ctx context.Context, ir *model.IngredientRecipe) error
	FindIngredientRecipeByIDFn func(ctx context.Context, id uint) (*model.IngredientRecipe, error)
	DeleteIngredientRecipeFn   func(ctx context.Context, id uint) error
	AddPreparationFn           func(ctx context.Context, p *model.Preparation) error
	FindPreparationByIDFn      func(ctx context.Context, id uint) (*model.Preparation, error)
	UpdatePreparationFn        func(ctx context.Context, p *model.Preparation) error
	DeletePreparationFn        func(ctx context.Context, id uint) error
}

func (m *RecipeRepoMock) Create(ctx context.Context, recipe *model.Recipe) error {
	if m.CreateFn == nil {
		panic("mocks: RecipeRepo.Create chamado sem CreateFn")
	}
	return m.CreateFn(ctx, recipe)
}

func (m *RecipeRepoMock) List(ctx context.Context, filter repository.RecipeFilter, page repository.Page) ([]model.Recipe, int64, error) {
	if m.ListFn == nil {
		panic("mocks: RecipeRepo.List chamado sem ListFn")
	}
	return m.ListFn(ctx, filter, page)
}

func (m *RecipeRepoMock) FindByID(ctx context.Context, id string) (*model.Recipe, error) {
	if m.FindByIDFn == nil {
		panic("mocks: RecipeRepo.FindByID chamado sem FindByIDFn")
	}
	return m.FindByIDFn(ctx, id)
}

func (m *RecipeRepoMock) Exists(ctx context.Context, id string) (bool, error) {
	if m.ExistsFn == nil {
		panic("mocks: RecipeRepo.Exists chamado sem ExistsFn")
	}
	return m.ExistsFn(ctx, id)
}

func (m *RecipeRepoMock) Update(ctx context.Context, recipe *model.Recipe) error {
	if m.UpdateFn == nil {
		panic("mocks: RecipeRepo.Update chamado sem UpdateFn")
	}
	return m.UpdateFn(ctx, recipe)
}

func (m *RecipeRepoMock) Replace(ctx context.Context, recipe *model.Recipe) error {
	if m.ReplaceFn == nil {
		panic("mocks: RecipeRepo.Replace chamado sem ReplaceFn")
	}
	return m.ReplaceFn(ctx, recipe)
}

func (m *RecipeRepoMock) Delete(ctx context.Context, id string) error {
	if m.DeleteFn == nil {
		panic("mocks: RecipeRepo.Delete chamado sem DeleteFn")
	}
	return m.DeleteFn(ctx, id)
}

func (m *RecipeRepoMock) AddImage(ctx context.Context, image *model.ImageRecipe) error {
	if m.AddImageFn == nil {
		panic("mocks: RecipeRepo.AddImage chamado sem AddImageFn")
	}
	return m.AddImageFn(ctx, image)
}

func (m *RecipeRepoMock) FindImageByID(ctx context.Context, id uint) (*model.ImageRecipe, error) {
	if m.FindImageByIDFn == nil {
		panic("mocks: RecipeRepo.FindImageByID chamado sem FindImageByIDFn")
	}
	return m.FindImageByIDFn(ctx, id)
}

func (m *RecipeRepoMock) UpdateImage(ctx context.Context, image *model.ImageRecipe) error {
	if m.UpdateImageFn == nil {
		panic("mocks: RecipeRepo.UpdateImage chamado sem UpdateImageFn")
	}
	return m.UpdateImageFn(ctx, image)
}

func (m *RecipeRepoMock) DeleteImage(ctx context.Context, id uint) error {
	if m.DeleteImageFn == nil {
		panic("mocks: RecipeRepo.DeleteImage chamado sem DeleteImageFn")
	}
	return m.DeleteImageFn(ctx, id)
}

func (m *RecipeRepoMock) AddIngredient(ctx context.Context, ir *model.IngredientRecipe) error {
	if m.AddIngredientFn == nil {
		panic("mocks: RecipeRepo.AddIngredient chamado sem AddIngredientFn")
	}
	return m.AddIngredientFn(ctx, ir)
}

func (m *RecipeRepoMock) FindIngredientRecipeByID(ctx context.Context, id uint) (*model.IngredientRecipe, error) {
	if m.FindIngredientRecipeByIDFn == nil {
		panic("mocks: RecipeRepo.FindIngredientRecipeByID chamado sem FindIngredientRecipeByIDFn")
	}
	return m.FindIngredientRecipeByIDFn(ctx, id)
}

func (m *RecipeRepoMock) DeleteIngredientRecipe(ctx context.Context, id uint) error {
	if m.DeleteIngredientRecipeFn == nil {
		panic("mocks: RecipeRepo.DeleteIngredientRecipe chamado sem DeleteIngredientRecipeFn")
	}
	return m.DeleteIngredientRecipeFn(ctx, id)
}

func (m *RecipeRepoMock) AddPreparation(ctx context.Context, p *model.Preparation) error {
	if m.AddPreparationFn == nil {
		panic("mocks: RecipeRepo.AddPreparation chamado sem AddPreparationFn")
	}
	return m.AddPreparationFn(ctx, p)
}

func (m *RecipeRepoMock) FindPreparationByID(ctx context.Context, id uint) (*model.Preparation, error) {
	if m.FindPreparationByIDFn == nil {
		panic("mocks: RecipeRepo.FindPreparationByID chamado sem FindPreparationByIDFn")
	}
	return m.FindPreparationByIDFn(ctx, id)
}

func (m *RecipeRepoMock) UpdatePreparation(ctx context.Context, p *model.Preparation) error {
	if m.UpdatePreparationFn == nil {
		panic("mocks: RecipeRepo.UpdatePreparation chamado sem UpdatePreparationFn")
	}
	return m.UpdatePreparationFn(ctx, p)
}

func (m *RecipeRepoMock) DeletePreparation(ctx context.Context, id uint) error {
	if m.DeletePreparationFn == nil {
		panic("mocks: RecipeRepo.DeletePreparation chamado sem DeletePreparationFn")
	}
	return m.DeletePreparationFn(ctx, id)
}
