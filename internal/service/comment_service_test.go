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

func recipeExists(exists bool) *mocks.RecipeRepoMock {
	return &mocks.RecipeRepoMock{ExistsFn: func(context.Context, string) (bool, error) { return exists, nil }}
}

func storedComment() *model.Comment {
	return &model.Comment{
		ID: 5, Description: "Ficou ótima!", UserID: userID, RecipeID: recipeID,
		User: model.User{ID: userID, Name: "Maria", Email: "maria@email.com"},
	}
}

func TestCommentService_Create(t *testing.T) {
	t.Run("grava e devolve com o autor", func(t *testing.T) {
		var saved *model.Comment
		repo := &mocks.CommentRepoMock{
			CreateFn: func(_ context.Context, c *model.Comment) error {
				saved = c
				c.ID = 5
				return nil
			},
			FindByIDFn: func(context.Context, uint) (*model.Comment, error) { return storedComment(), nil },
		}

		resp, err := service.NewCommentService(repo, recipeExists(true)).Create(ctx, userID, recipeID, service.CommentInput{Description: "Ficou ótima!"})

		require.NoError(t, err)
		assert.Equal(t, userID, saved.UserID)
		assert.Equal(t, recipeID, saved.RecipeID)
		assert.Equal(t, "Maria", resp.User.Name)
		assert.Equal(t, recipeID, resp.RecipeID)
	})

	t.Run("receita que não existe é 404", func(t *testing.T) {
		_, err := service.NewCommentService(&mocks.CommentRepoMock{}, recipeExists(false)).Create(ctx, userID, recipeID, service.CommentInput{Description: "Oi"})

		assert.ErrorIs(t, err, service.ErrRecipeNotFound)
	})

	t.Run("receita apagada no meio do caminho é 404", func(t *testing.T) {
		repo := &mocks.CommentRepoMock{CreateFn: func(context.Context, *model.Comment) error { return gorm.ErrForeignKeyViolated }}

		_, err := service.NewCommentService(repo, recipeExists(true)).Create(ctx, userID, recipeID, service.CommentInput{Description: "Oi"})

		assert.ErrorIs(t, err, service.ErrRecipeNotFound)
	})
}

func TestCommentService_ListByRecipe(t *testing.T) {
	t.Run("página padrão de 10", func(t *testing.T) {
		var got repository.Page
		repo := &mocks.CommentRepoMock{ListByRecipeFn: func(_ context.Context, id string, p repository.Page) ([]model.Comment, int64, error) {
			got = p
			return []model.Comment{*storedComment()}, 1, nil
		}}

		page, err := service.NewCommentService(repo, recipeExists(true)).ListByRecipe(ctx, recipeID, service.PageQuery{})

		require.NoError(t, err)
		assert.Equal(t, repository.Page{Number: 1, Size: service.DefaultCommentsPerPage}, got)
		require.Len(t, page.Items, 1)
		assert.Equal(t, "Maria", page.Items[0].User.Name)
	})

	t.Run("receita que não existe é 404", func(t *testing.T) {
		_, err := service.NewCommentService(&mocks.CommentRepoMock{}, recipeExists(false)).ListByRecipe(ctx, recipeID, service.PageQuery{})

		assert.ErrorIs(t, err, service.ErrRecipeNotFound)
	})
}

func TestCommentService_Update(t *testing.T) {
	t.Run("autor edita", func(t *testing.T) {
		repo := &mocks.CommentRepoMock{
			FindByIDFn: func(context.Context, uint) (*model.Comment, error) { return storedComment(), nil },
			UpdateFn: func(_ context.Context, c *model.Comment) error {
				assert.Equal(t, "Ficou melhor ainda.", c.Description)
				return nil
			},
		}

		_, err := service.NewCommentService(repo, recipeExists(true)).Update(ctx, 5, userID, service.CommentInput{Description: "Ficou melhor ainda."})

		require.NoError(t, err)
	})

	t.Run("outra pessoa, mesmo admin, não edita", func(t *testing.T) {
		repo := &mocks.CommentRepoMock{FindByIDFn: func(context.Context, uint) (*model.Comment, error) { return storedComment(), nil }}

		_, err := service.NewCommentService(repo, recipeExists(true)).Update(ctx, 5, otherID, service.CommentInput{Description: "Editado"})

		assert.ErrorIs(t, err, service.ErrForbiddenEditComment)
	})

	t.Run("comentário que não existe é 404", func(t *testing.T) {
		repo := &mocks.CommentRepoMock{FindByIDFn: func(context.Context, uint) (*model.Comment, error) { return nil, gorm.ErrRecordNotFound }}

		_, err := service.NewCommentService(repo, recipeExists(true)).Update(ctx, 5, userID, service.CommentInput{Description: "Editado"})

		assert.ErrorIs(t, err, service.ErrCommentNotFound)
	})
}

func TestCommentService_Delete(t *testing.T) {
	tests := []struct {
		name    string
		userID  string
		isAdm   bool
		wantErr error
	}{
		{name: "autor apaga", userID: userID},
		{name: "admin apaga o de qualquer um", userID: otherID, isAdm: true},
		{name: "outra pessoa recebe 403", userID: otherID, wantErr: service.ErrForbiddenDelComment},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			deleted := false
			repo := &mocks.CommentRepoMock{
				FindByIDFn: func(context.Context, uint) (*model.Comment, error) { return storedComment(), nil },
				DeleteFn: func(context.Context, uint) error {
					deleted = true
					return nil
				},
			}

			err := service.NewCommentService(repo, recipeExists(true)).Delete(ctx, 5, tt.userID, tt.isAdm)

			if tt.wantErr != nil {
				assert.ErrorIs(t, err, tt.wantErr)
				assert.False(t, deleted)
				return
			}
			require.NoError(t, err)
			assert.True(t, deleted)
		})
	}
}
