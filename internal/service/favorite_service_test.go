package service_test

import (
	"context"
	"errors"
	"testing"

	"panda-cooking-go-api/internal/model"
	"panda-cooking-go-api/internal/repository"
	"panda-cooking-go-api/internal/service"
	"panda-cooking-go-api/internal/service/mocks"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gorm.io/gorm"
)

func TestFavoriteService_Status(t *testing.T) {
	repo := &mocks.FavoriteRepoMock{ExistsFn: func(_ context.Context, u, r string) (bool, error) {
		return u == userID && r == recipeID, nil
	}}
	svc := service.NewFavoriteService(repo, recipeExists(true))

	status, err := svc.Status(ctx, userID, recipeID)
	require.NoError(t, err)
	assert.True(t, status.Favorite)

	status, err = svc.Status(ctx, otherID, recipeID)
	require.NoError(t, err)
	assert.False(t, status.Favorite)

	_, err = service.NewFavoriteService(repo, recipeExists(false)).Status(ctx, userID, recipeID)
	assert.ErrorIs(t, err, service.ErrRecipeNotFound)
}

func TestFavoriteService_Add(t *testing.T) {
	tests := []struct {
		name      string
		createErr error
		wantErr   error
	}{
		{name: "favorita"},
		{name: "favorito repetido é 409", createErr: gorm.ErrDuplicatedKey, wantErr: service.ErrAlreadyFavorite},
		{name: "receita apagada no meio do caminho é 404", createErr: gorm.ErrForeignKeyViolated, wantErr: service.ErrRecipeNotFound},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := &mocks.FavoriteRepoMock{CreateFn: func(_ context.Context, f *model.FavoriteRecipe) error {
				assert.Equal(t, model.FavoriteRecipe{UserID: userID, RecipeID: recipeID}, *f)
				return tt.createErr
			}}

			status, err := service.NewFavoriteService(repo, recipeExists(true)).Add(ctx, userID, recipeID)

			if tt.wantErr != nil {
				assert.ErrorIs(t, err, tt.wantErr)
				return
			}
			require.NoError(t, err)
			assert.True(t, status.Favorite)
		})
	}

	t.Run("receita que não existe é 404", func(t *testing.T) {
		_, err := service.NewFavoriteService(&mocks.FavoriteRepoMock{}, recipeExists(false)).Add(ctx, userID, recipeID)

		assert.ErrorIs(t, err, service.ErrRecipeNotFound)
	})
}

func TestFavoriteService_Remove(t *testing.T) {
	repo := &mocks.FavoriteRepoMock{DeleteFn: func(_ context.Context, u, r string) (bool, error) { return u == userID, nil }}
	svc := service.NewFavoriteService(repo, recipeExists(true))

	assert.NoError(t, svc.Remove(ctx, userID, recipeID))
	assert.ErrorIs(t, svc.Remove(ctx, otherID, recipeID), service.ErrFavoriteNotFound)

	failing := &mocks.FavoriteRepoMock{DeleteFn: func(context.Context, string, string) (bool, error) { return false, errors.New("banco fora") }}
	err := service.NewFavoriteService(failing, recipeExists(true)).Remove(ctx, userID, recipeID)
	require.Error(t, err)
	assert.NotErrorIs(t, err, service.ErrFavoriteNotFound)
}

func TestFavoriteService_List(t *testing.T) {
	var got repository.Page
	repo := &mocks.FavoriteRepoMock{ListRecipesFn: func(_ context.Context, u string, p repository.Page) ([]model.Recipe, int64, error) {
		assert.Equal(t, userID, u)
		got = p
		return []model.Recipe{*storedRecipe()}, 1, nil
	}}

	page, err := service.NewFavoriteService(repo, recipeExists(true)).List(ctx, userID, service.PageQuery{Page: 2})

	require.NoError(t, err)
	assert.Equal(t, repository.Page{Number: 2, Size: service.DefaultRecipesPerPage}, got)
	assert.Equal(t, int64(1), page.Total)
	assert.Equal(t, "Bolo de fubá", page.Items[0].Name)
}
