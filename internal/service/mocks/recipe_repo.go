package mocks

import "panda-cooking-go-api/internal/model"

type RecipeRepoMock struct {
	CreateFn                   func(recipe *model.Recipe) error
	FindAllFn                  func() ([]model.Recipe, error)
	FindByIDFn                 func(id string) (*model.Recipe, error)
	UpdateFn                   func(recipe *model.Recipe) error
	DeleteFn                   func(id string) error
	FindOrCreateIngredientFn   func(name string) (*model.Ingredient, error)
	AddIngredientFn            func(ir *model.IngredientRecipe) error
	FindIngredientRecipeByIDFn func(id uint) (*model.IngredientRecipe, error)
	DeleteIngredientRecipeFn   func(id uint) error
	AddImageFn                 func(image *model.ImageRecipe) error
	FindImageByIDFn            func(id uint) (*model.ImageRecipe, error)
	UpdateImageFn              func(image *model.ImageRecipe) error
	DeleteImageFn              func(id uint) error
	AddPreparationFn           func(p *model.Preparation) error
	FindPreparationByIDFn      func(id uint) (*model.Preparation, error)
	UpdatePreparationFn        func(p *model.Preparation) error
	DeletePreparationFn        func(id uint) error
}

func (m *RecipeRepoMock) Create(recipe *model.Recipe) error { return m.CreateFn(recipe) }
func (m *RecipeRepoMock) FindAll() ([]model.Recipe, error)  { return m.FindAllFn() }
func (m *RecipeRepoMock) FindByID(id string) (*model.Recipe, error) {
	return m.FindByIDFn(id)
}
func (m *RecipeRepoMock) Update(recipe *model.Recipe) error { return m.UpdateFn(recipe) }
func (m *RecipeRepoMock) Delete(id string) error            { return m.DeleteFn(id) }
func (m *RecipeRepoMock) FindOrCreateIngredient(name string) (*model.Ingredient, error) {
	return m.FindOrCreateIngredientFn(name)
}
func (m *RecipeRepoMock) AddIngredient(ir *model.IngredientRecipe) error {
	return m.AddIngredientFn(ir)
}
func (m *RecipeRepoMock) FindIngredientRecipeByID(id uint) (*model.IngredientRecipe, error) {
	return m.FindIngredientRecipeByIDFn(id)
}
func (m *RecipeRepoMock) DeleteIngredientRecipe(id uint) error    { return m.DeleteIngredientRecipeFn(id) }
func (m *RecipeRepoMock) AddImage(image *model.ImageRecipe) error { return m.AddImageFn(image) }
func (m *RecipeRepoMock) FindImageByID(id uint) (*model.ImageRecipe, error) {
	return m.FindImageByIDFn(id)
}
func (m *RecipeRepoMock) UpdateImage(image *model.ImageRecipe) error { return m.UpdateImageFn(image) }
func (m *RecipeRepoMock) DeleteImage(id uint) error                  { return m.DeleteImageFn(id) }
func (m *RecipeRepoMock) AddPreparation(p *model.Preparation) error  { return m.AddPreparationFn(p) }
func (m *RecipeRepoMock) FindPreparationByID(id uint) (*model.Preparation, error) {
	return m.FindPreparationByIDFn(id)
}
func (m *RecipeRepoMock) UpdatePreparation(p *model.Preparation) error {
	return m.UpdatePreparationFn(p)
}
func (m *RecipeRepoMock) DeletePreparation(id uint) error { return m.DeletePreparationFn(id) }
