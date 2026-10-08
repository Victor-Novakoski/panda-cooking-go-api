package service_test

import (
	"context"
	"testing"

	"panda-cooking-go-api/internal/model"
	"panda-cooking-go-api/internal/repository"
	"panda-cooking-go-api/internal/service"
	"panda-cooking-go-api/internal/service/mocks"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gorm.io/gorm"
)

const recipeID = "3d6f9a2c-8b1e-4f7a-9c0d-5e2b8a1f6c47"

func storedRecipe() *model.Recipe {
	return &model.Recipe{
		ID: recipeID, Name: "Bolo de fubá", Description: "Bolo simples de fubá", Time: "50 minutos", Portions: 8,
		UserID: userID, CategoryID: 12,
		User:         model.User{ID: userID, Name: "Maria", Email: "maria@email.com", PasswordHash: "hash"},
		Category:     model.Category{ID: 12, Name: "Pães e Bolos"},
		Images:       []model.ImageRecipe{{ID: 1, URL: "https://img.exemplo.com/bolo.jpg", RecipeID: recipeID}},
		Ingredients:  []model.IngredientRecipe{{ID: 1, Amount: "2 xícaras", RecipeID: recipeID, Ingredient: model.Ingredient{ID: 1, Name: "fubá"}}},
		Preparations: []model.Preparation{{ID: 1, Description: "Misture tudo.", RecipeID: recipeID}},
	}
}

func recipeInput() service.RecipeInput {
	return service.RecipeInput{
		Name: "Bolo de fubá", Description: "Bolo simples de fubá", Time: "50 minutos", Portions: 8, CategoryID: 12,
		Ingredients:  []service.IngredientRecipeInput{{Name: "  Farinha   de Trigo ", Amount: "2 xícaras"}},
		Preparations: []service.PreparationInput{{Description: "Misture tudo."}},
	}
}

func categories(exists bool) *mocks.CategoryRepoMock {
	return &mocks.CategoryRepoMock{ExistsFn: func(context.Context, uint) (bool, error) { return exists, nil }}
}

func findRecipe(r *model.Recipe) func(context.Context, string) (*model.Recipe, error) {
	return func(_ context.Context, id string) (*model.Recipe, error) {
		if r == nil || id != r.ID {
			return nil, gorm.ErrRecordNotFound
		}
		return r, nil
	}
}

func TestRecipeService_Create(t *testing.T) {
	t.Run("grava com o dono e ingrediente normalizado", func(t *testing.T) {
		var saved *model.Recipe
		repo := &mocks.RecipeRepoMock{
			CreateFn: func(_ context.Context, r *model.Recipe) error {
				saved = r
				r.ID = recipeID
				return nil
			},
			FindByIDFn: findRecipe(storedRecipe()),
		}

		resp, err := service.NewRecipeService(repo, categories(true)).Create(ctx, userID, recipeInput())

		require.NoError(t, err)
		assert.Equal(t, recipeID, resp.ID)
		assert.Equal(t, userID, saved.UserID)
		assert.Equal(t, "farinha de trigo", saved.Ingredients[0].Ingredient.Name)
		// a resposta não leva dado privado do autor
		assert.Equal(t, service.AuthorResponse{ID: userID, Name: "Maria"}, resp.Author)
	})

	t.Run("categoria que não existe é erro no campo category_id", func(t *testing.T) {
		_, err := service.NewRecipeService(&mocks.RecipeRepoMock{}, categories(false)).Create(ctx, userID, recipeInput())

		assert.ErrorIs(t, err, service.ErrCategoryNotFound)
		assert.Equal(t, "category_id", service.ErrCategoryNotFound.Field)
	})
}

func TestRecipeService_List(t *testing.T) {
	var gotFilter repository.RecipeFilter
	var gotPage repository.Page
	repo := &mocks.RecipeRepoMock{ListFn: func(_ context.Context, f repository.RecipeFilter, p repository.Page) ([]model.Recipe, int64, error) {
		gotFilter, gotPage = f, p
		return []model.Recipe{*storedRecipe()}, 25, nil
	}}

	page, err := service.NewRecipeService(repo, categories(true)).List(ctx, service.ListRecipesQuery{Search: "bolo", CategoryID: 12})

	require.NoError(t, err)
	assert.Equal(t, repository.RecipeFilter{Search: "bolo", CategoryID: 12}, gotFilter)
	assert.Equal(t, repository.Page{Number: 1, Size: service.DefaultRecipesPerPage}, gotPage)
	assert.Equal(t, int64(25), page.Total)
	assert.Equal(t, 3, page.TotalPages)
	require.Len(t, page.Items, 1)
	assert.Equal(t, "https://img.exemplo.com/bolo.jpg", page.Items[0].ImageURL)
	assert.Equal(t, "Pães e Bolos", page.Items[0].Category.Name)
}

func TestRecipeService_ListEmpty(t *testing.T) {
	repo := &mocks.RecipeRepoMock{ListFn: func(context.Context, repository.RecipeFilter, repository.Page) ([]model.Recipe, int64, error) {
		return nil, 0, nil
	}}

	page, err := service.NewRecipeService(repo, categories(true)).List(ctx, service.ListRecipesQuery{Page: 3, PerPage: 5})

	require.NoError(t, err)
	assert.NotNil(t, page.Items, "items sai como [] no JSON, não null")
	assert.Equal(t, 3, page.Page)
	assert.Equal(t, 5, page.PerPage)
	assert.Equal(t, 0, page.TotalPages)
}

func TestRecipeService_GetByID_NotFound(t *testing.T) {
	repo := &mocks.RecipeRepoMock{FindByIDFn: findRecipe(nil)}

	_, err := service.NewRecipeService(repo, categories(true)).GetByID(ctx, recipeID)

	assert.ErrorIs(t, err, service.ErrRecipeNotFound)
}

func TestRecipeService_Update(t *testing.T) {
	t.Run("muda só os campos enviados", func(t *testing.T) {
		var saved model.Recipe
		repo := &mocks.RecipeRepoMock{
			FindByIDFn: findRecipe(storedRecipe()),
			UpdateFn: func(_ context.Context, r *model.Recipe) error {
				saved = *r
				return nil
			},
		}

		_, err := service.NewRecipeService(repo, &mocks.CategoryRepoMock{}).Update(ctx, recipeID, userID, service.UpdateRecipeInput{Portions: ptr(10)})

		require.NoError(t, err)
		assert.Equal(t, 10, saved.Portions)
		assert.Equal(t, "Bolo de fubá", saved.Name)
		assert.Equal(t, uint(12), saved.CategoryID)
	})

	t.Run("categoria nova é conferida", func(t *testing.T) {
		repo := &mocks.RecipeRepoMock{FindByIDFn: findRecipe(storedRecipe())}

		_, err := service.NewRecipeService(repo, categories(false)).Update(ctx, recipeID, userID, service.UpdateRecipeInput{CategoryID: ptr(uint(99))})

		assert.ErrorIs(t, err, service.ErrCategoryNotFound)
	})

	t.Run("receita de outra pessoa é 403", func(t *testing.T) {
		repo := &mocks.RecipeRepoMock{FindByIDFn: findRecipe(storedRecipe())}

		_, err := service.NewRecipeService(repo, categories(true)).Update(ctx, recipeID, otherID, service.UpdateRecipeInput{Name: ptr("Meu bolo")})

		assert.ErrorIs(t, err, service.ErrForbiddenEditRecipe)
	})
}

func TestRecipeService_Replace(t *testing.T) {
	t.Run("troca tudo mantendo id e dono", func(t *testing.T) {
		var saved *model.Recipe
		repo := &mocks.RecipeRepoMock{
			FindByIDFn: findRecipe(storedRecipe()),
			ReplaceFn: func(_ context.Context, r *model.Recipe) error {
				saved = r
				return nil
			},
		}

		_, err := service.NewRecipeService(repo, categories(true)).Replace(ctx, recipeID, userID, recipeInput())

		require.NoError(t, err)
		assert.Equal(t, recipeID, saved.ID)
		assert.Equal(t, userID, saved.UserID)
		assert.Len(t, saved.Ingredients, 1)
		assert.Empty(t, saved.Images)
	})

	t.Run("receita de outra pessoa é 403", func(t *testing.T) {
		repo := &mocks.RecipeRepoMock{FindByIDFn: findRecipe(storedRecipe())}

		_, err := service.NewRecipeService(repo, categories(true)).Replace(ctx, recipeID, otherID, recipeInput())

		assert.ErrorIs(t, err, service.ErrForbiddenEditRecipe)
	})
}

func TestRecipeService_Delete(t *testing.T) {
	t.Run("dono apaga", func(t *testing.T) {
		deleted := false
		repo := &mocks.RecipeRepoMock{
			FindByIDFn: findRecipe(storedRecipe()),
			DeleteFn: func(context.Context, string) error {
				deleted = true
				return nil
			},
		}

		require.NoError(t, service.NewRecipeService(repo, categories(true)).Delete(ctx, recipeID, userID))
		assert.True(t, deleted)
	})

	t.Run("outra pessoa recebe 403", func(t *testing.T) {
		repo := &mocks.RecipeRepoMock{FindByIDFn: findRecipe(storedRecipe())}

		err := service.NewRecipeService(repo, categories(true)).Delete(ctx, recipeID, otherID)

		assert.ErrorIs(t, err, service.ErrForbiddenDelRecipe)
	})

	t.Run("receita que não existe é 404", func(t *testing.T) {
		repo := &mocks.RecipeRepoMock{FindByIDFn: findRecipe(nil)}

		assert.ErrorIs(t, service.NewRecipeService(repo, categories(true)).Delete(ctx, recipeID, userID), service.ErrRecipeNotFound)
	})
}

// O dono de uma receita não pode mexer em item de outra receita passando o
// id do item na URL da sua (IDOR).
func TestRecipeService_ItemOfAnotherRecipe(t *testing.T) {
	const otherRecipe = "receita-de-outra-pessoa"
	repo := &mocks.RecipeRepoMock{
		FindByIDFn: findRecipe(storedRecipe()),
		FindImageByIDFn: func(_ context.Context, id uint) (*model.ImageRecipe, error) {
			return &model.ImageRecipe{ID: id, RecipeID: otherRecipe}, nil
		},
		FindIngredientRecipeByIDFn: func(_ context.Context, id uint) (*model.IngredientRecipe, error) {
			return &model.IngredientRecipe{ID: id, RecipeID: otherRecipe}, nil
		},
		FindPreparationByIDFn: func(_ context.Context, id uint) (*model.Preparation, error) {
			return &model.Preparation{ID: id, RecipeID: otherRecipe}, nil
		},
	}
	svc := service.NewRecipeService(repo, categories(true))

	_, err := svc.UpdateImage(ctx, recipeID, userID, 7, service.ImageRecipeInput{URL: "https://img.exemplo.com/x.jpg"})
	assert.ErrorIs(t, err, service.ErrImageNotFound)
	assert.ErrorIs(t, svc.DeleteImage(ctx, recipeID, userID, 7), service.ErrImageNotFound)
	assert.ErrorIs(t, svc.DeleteIngredient(ctx, recipeID, userID, 7), service.ErrIngredientNotFound)
	_, err = svc.UpdatePreparation(ctx, recipeID, userID, 7, service.PreparationInput{Description: "Outro passo."})
	assert.ErrorIs(t, err, service.ErrPreparationNotFound)
	assert.ErrorIs(t, svc.DeletePreparation(ctx, recipeID, userID, 7), service.ErrPreparationNotFound)
}

func TestRecipeService_Items(t *testing.T) {
	t.Run("ingrediente novo é normalizado", func(t *testing.T) {
		repo := &mocks.RecipeRepoMock{
			FindByIDFn: findRecipe(storedRecipe()),
			AddIngredientFn: func(_ context.Context, ir *model.IngredientRecipe) error {
				ir.ID = 9
				return nil
			},
		}

		resp, err := service.NewRecipeService(repo, categories(true)).AddIngredient(ctx, recipeID, userID, service.IngredientRecipeInput{Name: "Açúcar  Mascavo", Amount: "1 xícara"})

		require.NoError(t, err)
		assert.Equal(t, service.IngredientRecipeResponse{ID: 9, Name: "açúcar mascavo", Amount: "1 xícara"}, *resp)
	})

	t.Run("décima primeira foto é recusada", func(t *testing.T) {
		full := storedRecipe()
		full.Images = make([]model.ImageRecipe, service.MaxImages)
		repo := &mocks.RecipeRepoMock{FindByIDFn: findRecipe(full)}

		_, err := service.NewRecipeService(repo, categories(true)).AddImage(ctx, recipeID, userID, service.ImageRecipeInput{URL: "https://img.exemplo.com/y.jpg"})

		assert.ErrorIs(t, err, service.ErrTooManyImages)
	})

	t.Run("passo além do limite é recusado", func(t *testing.T) {
		full := storedRecipe()
		full.Preparations = make([]model.Preparation, service.MaxPreparations)
		repo := &mocks.RecipeRepoMock{FindByIDFn: findRecipe(full)}

		_, err := service.NewRecipeService(repo, categories(true)).AddPreparation(ctx, recipeID, userID, service.PreparationInput{Description: "Mais um."})

		assert.ErrorIs(t, err, service.ErrTooManyItems)
	})

	t.Run("item de receita alheia é 403", func(t *testing.T) {
		repo := &mocks.RecipeRepoMock{FindByIDFn: findRecipe(storedRecipe())}

		_, err := service.NewRecipeService(repo, categories(true)).AddImage(ctx, recipeID, otherID, service.ImageRecipeInput{URL: "https://img.exemplo.com/y.jpg"})

		assert.ErrorIs(t, err, service.ErrForbidden)
	})

	t.Run("troca a foto da própria receita", func(t *testing.T) {
		var saved *model.ImageRecipe
		repo := &mocks.RecipeRepoMock{
			FindByIDFn: findRecipe(storedRecipe()),
			FindImageByIDFn: func(_ context.Context, id uint) (*model.ImageRecipe, error) {
				return &model.ImageRecipe{ID: id, URL: "https://img.exemplo.com/velha.jpg", RecipeID: recipeID}, nil
			},
			UpdateImageFn: func(_ context.Context, img *model.ImageRecipe) error {
				saved = img
				return nil
			},
		}

		resp, err := service.NewRecipeService(repo, categories(true)).UpdateImage(ctx, recipeID, userID, 1, service.ImageRecipeInput{URL: "https://img.exemplo.com/nova.jpg"})

		require.NoError(t, err)
		assert.Equal(t, "https://img.exemplo.com/nova.jpg", resp.URL)
		assert.Equal(t, "https://img.exemplo.com/nova.jpg", saved.URL)
	})
}
