package service_test

import (
	"errors"
	"panda-cooking-go-api/internal/model"
	"panda-cooking-go-api/internal/service"
	"panda-cooking-go-api/internal/service/mocks"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gorm.io/gorm"
)

func baseRecipe() model.Recipe {
	return model.Recipe{
		ID:       "recipe-uuid",
		Name:     "Panqueca",
		UserID:   "user-uuid",
		Category: model.Category{ID: 1, Name: "Café da manhã"},
	}
}

func TestRecipeService_GetAll(t *testing.T) {
	t.Run("retorna lista de receitas", func(t *testing.T) {
		repo := &mocks.RecipeRepoMock{
			FindAllFn: func() ([]model.Recipe, error) {
				return []model.Recipe{baseRecipe(), baseRecipe()}, nil
			},
		}
		svc := service.NewRecipeService(repo)

		recipes, err := svc.GetAll()

		require.NoError(t, err)
		assert.Len(t, recipes, 2)
		assert.Equal(t, "Panqueca", recipes[0].Name)
	})

	t.Run("retorna slice vazio quando não há receitas", func(t *testing.T) {
		repo := &mocks.RecipeRepoMock{
			FindAllFn: func() ([]model.Recipe, error) {
				return []model.Recipe{}, nil
			},
		}
		svc := service.NewRecipeService(repo)

		recipes, err := svc.GetAll()

		require.NoError(t, err)
		assert.Len(t, recipes, 0)
	})
}

func TestRecipeService_GetByID(t *testing.T) {
	t.Run("retorna receita existente", func(t *testing.T) {
		r := baseRecipe()
		repo := &mocks.RecipeRepoMock{
			FindByIDFn: func(id string) (*model.Recipe, error) { return &r, nil },
		}
		svc := service.NewRecipeService(repo)

		resp, err := svc.GetByID("recipe-uuid")

		require.NoError(t, err)
		assert.Equal(t, "recipe-uuid", resp.ID)
		assert.Equal(t, "Panqueca", resp.Name)
	})

	t.Run("retorna erro quando não encontrada", func(t *testing.T) {
		repo := &mocks.RecipeRepoMock{
			FindByIDFn: func(id string) (*model.Recipe, error) {
				return nil, gorm.ErrRecordNotFound
			},
		}
		svc := service.NewRecipeService(repo)

		_, err := svc.GetByID("id-invalido")

		require.Error(t, err)
		assert.Equal(t, "receita não encontrada", err.Error())
	})
}

func TestRecipeService_Create(t *testing.T) {
	t.Run("cria receita com ingredientes", func(t *testing.T) {
		r := baseRecipe()
		repo := &mocks.RecipeRepoMock{
			CreateFn: func(recipe *model.Recipe) error {
				recipe.ID = "recipe-uuid"
				return nil
			},
			FindByIDFn: func(id string) (*model.Recipe, error) { return &r, nil },
			FindOrCreateIngredientFn: func(name string) (*model.Ingredient, error) {
				return &model.Ingredient{ID: 1, Name: name}, nil
			},
			AddIngredientFn: func(ir *model.IngredientRecipe) error { return nil },
		}
		svc := service.NewRecipeService(repo)

		resp, err := svc.Create("user-uuid", service.CreateRecipeInput{
			Name:        "Panqueca",
			Description: "Gostosa",
			Time:        "20min",
			Portions:    4,
			CategoryID:  1,
			Ingredients: []service.IngredientRecipeInput{
				{Name: "farinha", Amount: "2 xícaras"},
			},
		})

		require.NoError(t, err)
		assert.Equal(t, "recipe-uuid", resp.ID)
	})
}

func TestRecipeService_Update(t *testing.T) {
	t.Run("dono pode atualizar receita", func(t *testing.T) {
		r := baseRecipe()
		updated := baseRecipe()
		updated.Name = "Panqueca Americana"

		repo := &mocks.RecipeRepoMock{
			FindByIDFn: func(id string) (*model.Recipe, error) {
				if id == "recipe-uuid" {
					return &r, nil
				}
				return &updated, nil
			},
			UpdateFn: func(recipe *model.Recipe) error { return nil },
		}
		svc := service.NewRecipeService(repo)

		resp, err := svc.Update("recipe-uuid", "user-uuid", service.UpdateRecipeInput{
			Name: "Panqueca Americana",
		})

		require.NoError(t, err)
		assert.Equal(t, "Panqueca Americana", resp.Name)
	})

	t.Run("outro usuário não pode atualizar", func(t *testing.T) {
		r := baseRecipe() // UserID = "user-uuid"
		repo := &mocks.RecipeRepoMock{
			FindByIDFn: func(id string) (*model.Recipe, error) { return &r, nil },
		}
		svc := service.NewRecipeService(repo)

		_, err := svc.Update("recipe-uuid", "outro-user", service.UpdateRecipeInput{Name: "Hackeada"})

		require.Error(t, err)
		assert.Equal(t, "sem permissão para editar esta receita", err.Error())
	})
}

func TestRecipeService_Delete(t *testing.T) {
	t.Run("dono pode deletar receita", func(t *testing.T) {
		r := baseRecipe()
		repo := &mocks.RecipeRepoMock{
			FindByIDFn: func(id string) (*model.Recipe, error) { return &r, nil },
			DeleteFn:   func(id string) error { return nil },
		}
		svc := service.NewRecipeService(repo)

		err := svc.Delete("recipe-uuid", "user-uuid")
		require.NoError(t, err)
	})

	t.Run("outro usuário recebe erro de permissão", func(t *testing.T) {
		r := baseRecipe()
		repo := &mocks.RecipeRepoMock{
			FindByIDFn: func(id string) (*model.Recipe, error) { return &r, nil },
		}
		svc := service.NewRecipeService(repo)

		err := svc.Delete("recipe-uuid", "invasor")
		require.Error(t, err)
		assert.Equal(t, "sem permissão para deletar esta receita", err.Error())
	})

	t.Run("receita inexistente retorna erro", func(t *testing.T) {
		repo := &mocks.RecipeRepoMock{
			FindByIDFn: func(id string) (*model.Recipe, error) {
				return nil, gorm.ErrRecordNotFound
			},
		}
		svc := service.NewRecipeService(repo)

		err := svc.Delete("id-fantasma", "user-uuid")
		require.Error(t, err)
		assert.Equal(t, "receita não encontrada", err.Error())
	})
}

func TestRecipeService_AddIngredient(t *testing.T) {
	t.Run("adiciona ingrediente normalizando o nome", func(t *testing.T) {
		r := baseRecipe()
		var nomeGravado string

		repo := &mocks.RecipeRepoMock{
			FindByIDFn: func(id string) (*model.Recipe, error) { return &r, nil },
			FindOrCreateIngredientFn: func(name string) (*model.Ingredient, error) {
				nomeGravado = name // captura o nome que chegou ao repo
				return &model.Ingredient{ID: 1, Name: name}, nil
			},
			AddIngredientFn: func(ir *model.IngredientRecipe) error { return nil },
		}
		svc := service.NewRecipeService(repo)

		_, err := svc.AddIngredient("recipe-uuid", "user-uuid", service.IngredientRecipeInput{
			Name:   "  Farinha de Trigo  ", // maiúscula e com espaços
			Amount: "2 xícaras",
		})

		require.NoError(t, err)
		assert.Equal(t, "farinha de trigo", nomeGravado) // deve chegar normalizado
	})

	t.Run("outro usuário não pode adicionar ingrediente", func(t *testing.T) {
		r := baseRecipe()
		repo := &mocks.RecipeRepoMock{
			FindByIDFn: func(id string) (*model.Recipe, error) { return &r, nil },
		}
		svc := service.NewRecipeService(repo)

		_, err := svc.AddIngredient("recipe-uuid", "invasor", service.IngredientRecipeInput{
			Name: "sal", Amount: "1 pitada",
		})

		require.Error(t, err)
		assert.Equal(t, "sem permissão", err.Error())
	})
}

func TestRecipeService_Ingredient_ErrorPropagation(t *testing.T) {
	t.Run("erro no banco ao criar ingrediente é propagado", func(t *testing.T) {
		r := baseRecipe()
		repo := &mocks.RecipeRepoMock{
			FindByIDFn: func(id string) (*model.Recipe, error) { return &r, nil },
			FindOrCreateIngredientFn: func(name string) (*model.Ingredient, error) {
				return nil, errors.New("db error")
			},
		}
		svc := service.NewRecipeService(repo)

		_, err := svc.AddIngredient("recipe-uuid", "user-uuid", service.IngredientRecipeInput{
			Name: "sal", Amount: "1 pitada",
		})

		require.Error(t, err)
		assert.Equal(t, "db error", err.Error())
	})
}
