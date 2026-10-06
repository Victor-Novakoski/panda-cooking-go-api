package service_test

import (
	"errors"
	"testing"

	"panda-cooking-go-api/internal/model"
	"panda-cooking-go-api/internal/service"
	"panda-cooking-go-api/internal/service/mocks"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gorm.io/gorm"
)

func ptr(s string) *string { return &s }

func TestRecipeService_Replace(t *testing.T) {
	input := service.ReplaceRecipeInput{
		Name:         "Panqueca Americana",
		Description:  "Fofinha e alta",
		Time:         "20 minutos",
		Portions:     4,
		CategoryID:   3,
		Images:       []service.ImageRecipeInput{{URL: "https://exemplo.com/panqueca.jpg"}},
		Ingredients:  []service.IngredientRecipeInput{{Name: "  Farinha ", Amount: "2 xícaras"}},
		Preparations: []service.PreparationInput{{Description: "Misture tudo"}, {Description: "Frite"}},
	}

	t.Run("dono troca a receita inteira", func(t *testing.T) {
		r := baseRecipe()
		var saved *model.Recipe
		repo := &mocks.RecipeRepoMock{
			FindByIDFn: func(string) (*model.Recipe, error) { return &r, nil },
			ReplaceFn:  func(recipe *model.Recipe) error { saved = recipe; return nil },
		}

		_, err := service.NewRecipeService(repo).Replace("recipe-uuid", "user-uuid", input)

		require.NoError(t, err)
		require.NotNil(t, saved)
		assert.Equal(t, "recipe-uuid", saved.ID)
		assert.Equal(t, "user-uuid", saved.UserID, "o dono não muda")
		assert.Equal(t, uint(3), saved.CategoryID)
		assert.Equal(t, "farinha", saved.Ingredients[0].Ingredient.Name, "nome do ingrediente normalizado")
		assert.Equal(t, "2 xícaras", saved.Ingredients[0].Amount)
		require.Len(t, saved.Preparations, 2)
		assert.Equal(t, "Frite", saved.Preparations[1].Description, "passos na ordem recebida")
		assert.Equal(t, "https://exemplo.com/panqueca.jpg", saved.Images[0].URL)
	})

	t.Run("outro usuário não pode trocar", func(t *testing.T) {
		r := baseRecipe()
		repo := &mocks.RecipeRepoMock{
			FindByIDFn: func(string) (*model.Recipe, error) { return &r, nil },
			ReplaceFn:  func(*model.Recipe) error { t.Fatal("não deveria gravar"); return nil },
		}

		_, err := service.NewRecipeService(repo).Replace("recipe-uuid", "outro-uuid", input)

		assert.ErrorIs(t, err, service.ErrForbiddenEditRecipe)
	})

	t.Run("receita inexistente", func(t *testing.T) {
		repo := &mocks.RecipeRepoMock{
			FindByIDFn: func(string) (*model.Recipe, error) { return nil, gorm.ErrRecordNotFound },
		}

		_, err := service.NewRecipeService(repo).Replace("nao-existe", "user-uuid", input)

		assert.ErrorIs(t, err, service.ErrRecipeNotFound)
	})

	t.Run("categoria inexistente vira erro de validação", func(t *testing.T) {
		r := baseRecipe()
		repo := &mocks.RecipeRepoMock{
			FindByIDFn: func(string) (*model.Recipe, error) { return &r, nil },
			ReplaceFn:  func(*model.Recipe) error { return gorm.ErrForeignKeyViolated },
		}

		_, err := service.NewRecipeService(repo).Replace("recipe-uuid", "user-uuid", input)

		assert.ErrorIs(t, err, service.ErrCategoryNotFound)
	})

	t.Run("falha ao gravar é propagada", func(t *testing.T) {
		r := baseRecipe()
		repo := &mocks.RecipeRepoMock{
			FindByIDFn: func(string) (*model.Recipe, error) { return &r, nil },
			ReplaceFn:  func(*model.Recipe) error { return errors.New("tx falhou") },
		}

		_, err := service.NewRecipeService(repo).Replace("recipe-uuid", "user-uuid", input)

		assert.EqualError(t, err, "tx falhou")
	})
}

func TestCommentService_GetByRecipe(t *testing.T) {
	t.Run("lista os comentários sem o e-mail do autor", func(t *testing.T) {
		r := baseRecipe()
		commentRepo := &mocks.CommentRepoMock{
			FindByRecipeFn: func(recipeID string) ([]model.Comment, error) {
				assert.Equal(t, "recipe-uuid", recipeID)
				return []model.Comment{{
					ID: 1, Description: "Muito bom", RecipeID: recipeID, UserID: "user-uuid",
					User: model.User{ID: "user-uuid", Name: "Maria", Email: "maria@email.com"},
				}}, nil
			},
		}
		recipeRepo := &mocks.RecipeRepoMock{FindByIDFn: func(string) (*model.Recipe, error) { return &r, nil }}

		comments, err := setupCommentService(commentRepo, recipeRepo).GetByRecipe("recipe-uuid")

		require.NoError(t, err)
		require.Len(t, comments, 1)
		assert.Equal(t, "Maria", comments[0].User.Name)
		assert.Equal(t, "user-uuid", comments[0].User.ID)
	})

	t.Run("receita inexistente é 404", func(t *testing.T) {
		recipeRepo := &mocks.RecipeRepoMock{
			FindByIDFn: func(string) (*model.Recipe, error) { return nil, gorm.ErrRecordNotFound },
		}

		_, err := setupCommentService(&mocks.CommentRepoMock{}, recipeRepo).GetByRecipe("x")

		assert.ErrorIs(t, err, service.ErrRecipeNotFound)
	})
}

func TestUserService_Update(t *testing.T) {
	newRepo := func(saved **model.User) *mocks.UserRepoMock {
		return &mocks.UserRepoMock{
			FindByIDFn: func(id string) (*model.User, error) {
				return &model.User{ID: id, Name: "Maria", ImageProfile: "https://exemplo.com/maria.jpg"}, nil
			},
			UpdateFn: func(u *model.User) error { *saved = u; return nil },
		}
	}

	t.Run("campo que não veio fica como está", func(t *testing.T) {
		var saved *model.User
		svc := service.NewUserService(newRepo(&saved), "test-secret")

		resp, err := svc.Update("user-uuid", service.UpdateUserInput{Name: ptr("  Maria Silva ")})

		require.NoError(t, err)
		assert.Equal(t, "Maria Silva", resp.Name)
		assert.Equal(t, "https://exemplo.com/maria.jpg", saved.ImageProfile)
	})

	t.Run("foto vazia remove a foto", func(t *testing.T) {
		var saved *model.User
		svc := service.NewUserService(newRepo(&saved), "test-secret")

		resp, err := svc.Update("user-uuid", service.UpdateUserInput{ImageProfile: ptr("")})

		require.NoError(t, err)
		assert.Empty(t, resp.ImageProfile)
		assert.Equal(t, "Maria", saved.Name)
	})

	t.Run("nome só com espaços é recusado", func(t *testing.T) {
		var saved *model.User
		svc := service.NewUserService(newRepo(&saved), "test-secret")

		_, err := svc.Update("user-uuid", service.UpdateUserInput{Name: ptr("   ")})

		assert.ErrorIs(t, err, service.ErrEmptyName)
		assert.Nil(t, saved, "não deveria gravar")
	})
}
