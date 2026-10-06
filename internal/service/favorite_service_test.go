package service_test

import (
	"panda-cooking-go-api/internal/model"
	"panda-cooking-go-api/internal/service"
	"panda-cooking-go-api/internal/service/mocks"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gorm.io/gorm"
)

func setupFavoriteService(favRepo *mocks.FavoriteRepoMock, recipeRepo *mocks.RecipeRepoMock) *service.FavoriteService {
	return service.NewFavoriteService(favRepo, recipeRepo)
}

func TestFavoriteService_Add(t *testing.T) {
	r := baseRecipe()

	t.Run("adiciona receita aos favoritos", func(t *testing.T) {
		favRepo := &mocks.FavoriteRepoMock{
			FindFn: func(userID, recipeID string) (*model.FavoriteRecipe, error) {
				return nil, gorm.ErrRecordNotFound // ainda não é favorito
			},
			CreateFn: func(fav *model.FavoriteRecipe) error {
				fav.ID = 1
				return nil
			},
		}
		recipeRepo := &mocks.RecipeRepoMock{
			FindByIDFn: func(id string) (*model.Recipe, error) { return &r, nil },
		}
		svc := setupFavoriteService(favRepo, recipeRepo)

		resp, err := svc.Add("user-uuid", "recipe-uuid")

		require.NoError(t, err)
		assert.Equal(t, uint(1), resp.ID)
		assert.Equal(t, "Panqueca", resp.Recipe.Name)
	})

	t.Run("receita já favoritada retorna conflito", func(t *testing.T) {
		favRepo := &mocks.FavoriteRepoMock{
			FindFn: func(userID, recipeID string) (*model.FavoriteRecipe, error) {
				return &model.FavoriteRecipe{ID: 1}, nil // já existe
			},
		}
		recipeRepo := &mocks.RecipeRepoMock{
			FindByIDFn: func(id string) (*model.Recipe, error) { return &r, nil },
		}
		svc := setupFavoriteService(favRepo, recipeRepo)

		_, err := svc.Add("user-uuid", "recipe-uuid")

		require.Error(t, err)
		assert.Equal(t, "receita já está nos favoritos", err.Error())
	})

	t.Run("receita inexistente retorna erro", func(t *testing.T) {
		favRepo := &mocks.FavoriteRepoMock{}
		recipeRepo := &mocks.RecipeRepoMock{
			FindByIDFn: func(id string) (*model.Recipe, error) {
				return nil, gorm.ErrRecordNotFound
			},
		}
		svc := setupFavoriteService(favRepo, recipeRepo)

		_, err := svc.Add("user-uuid", "id-invalido")

		require.Error(t, err)
		assert.Equal(t, "receita não encontrada", err.Error())
	})
}

func TestFavoriteService_Remove(t *testing.T) {
	t.Run("remove receita dos favoritos", func(t *testing.T) {
		favRepo := &mocks.FavoriteRepoMock{
			FindFn: func(userID, recipeID string) (*model.FavoriteRecipe, error) {
				return &model.FavoriteRecipe{ID: 1}, nil
			},
			DeleteFn: func(id uint) error { return nil },
		}
		svc := setupFavoriteService(favRepo, nil)

		err := svc.Remove("user-uuid", "recipe-uuid")
		require.NoError(t, err)
	})

	t.Run("receita não estava nos favoritos retorna erro", func(t *testing.T) {
		favRepo := &mocks.FavoriteRepoMock{
			FindFn: func(userID, recipeID string) (*model.FavoriteRecipe, error) {
				return nil, gorm.ErrRecordNotFound
			},
		}
		svc := setupFavoriteService(favRepo, nil)

		err := svc.Remove("user-uuid", "recipe-uuid")

		require.Error(t, err)
		assert.Equal(t, "receita não está nos favoritos", err.Error())
	})
}
