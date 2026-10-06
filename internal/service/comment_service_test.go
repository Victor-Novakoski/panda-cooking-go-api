package service_test

import (
	"testing"

	"panda-cooking-go-api/internal/model"
	"panda-cooking-go-api/internal/service"
	"panda-cooking-go-api/internal/service/mocks"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gorm.io/gorm"
)

func setupCommentService(commentRepo *mocks.CommentRepoMock, recipeRepo *mocks.RecipeRepoMock) *service.CommentService {
	return service.NewCommentService(commentRepo, recipeRepo)
}

func TestCommentService_Create(t *testing.T) {
	r := baseRecipe()
	comment := model.Comment{ID: 1, Description: "Delicioso!", RecipeID: "recipe-uuid", UserID: "user-uuid"}

	t.Run("cria comentário em receita existente", func(t *testing.T) {
		commentRepo := &mocks.CommentRepoMock{
			CreateFn:   func(c *model.Comment) error { c.ID = 1; return nil },
			FindByIDFn: func(id uint) (*model.Comment, error) { return &comment, nil },
		}
		recipeRepo := &mocks.RecipeRepoMock{
			FindByIDFn: func(id string) (*model.Recipe, error) { return &r, nil },
		}
		svc := setupCommentService(commentRepo, recipeRepo)

		resp, err := svc.Create("user-uuid", service.CreateCommentInput{
			Description: "Delicioso!",
			RecipeID:    "recipe-uuid",
		})

		require.NoError(t, err)
		assert.Equal(t, "Delicioso!", resp.Description)
	})

	t.Run("receita inexistente retorna erro", func(t *testing.T) {
		commentRepo := &mocks.CommentRepoMock{}
		recipeRepo := &mocks.RecipeRepoMock{
			FindByIDFn: func(id string) (*model.Recipe, error) {
				return nil, gorm.ErrRecordNotFound
			},
		}
		svc := setupCommentService(commentRepo, recipeRepo)

		_, err := svc.Create("user-uuid", service.CreateCommentInput{
			Description: "Oi",
			RecipeID:    "id-invalido",
		})

		require.Error(t, err)
		assert.Equal(t, "receita não encontrada", err.Error())
	})
}

func TestCommentService_Delete(t *testing.T) {
	t.Run("dono do comentário pode deletar", func(t *testing.T) {
		c := model.Comment{ID: 1, UserID: "user-uuid"}
		commentRepo := &mocks.CommentRepoMock{
			FindByIDFn: func(id uint) (*model.Comment, error) { return &c, nil },
			DeleteFn:   func(id uint) error { return nil },
		}
		svc := setupCommentService(commentRepo, nil)

		err := svc.Delete(1, "user-uuid", false)
		require.NoError(t, err)
	})

	t.Run("admin pode deletar comentário de qualquer usuário", func(t *testing.T) {
		c := model.Comment{ID: 1, UserID: "outro-user"}
		commentRepo := &mocks.CommentRepoMock{
			FindByIDFn: func(id uint) (*model.Comment, error) { return &c, nil },
			DeleteFn:   func(id uint) error { return nil },
		}
		svc := setupCommentService(commentRepo, nil)

		err := svc.Delete(1, "admin-uuid", true) // isAdm = true
		require.NoError(t, err)
	})

	t.Run("usuário comum não pode deletar comentário alheio", func(t *testing.T) {
		c := model.Comment{ID: 1, UserID: "dono-uuid"}
		commentRepo := &mocks.CommentRepoMock{
			FindByIDFn: func(id uint) (*model.Comment, error) { return &c, nil },
		}
		svc := setupCommentService(commentRepo, nil)

		err := svc.Delete(1, "outro-user", false)
		require.Error(t, err)
		assert.Equal(t, "sem permissão para deletar este comentário", err.Error())
	})

	t.Run("comentário inexistente retorna erro", func(t *testing.T) {
		commentRepo := &mocks.CommentRepoMock{
			FindByIDFn: func(id uint) (*model.Comment, error) {
				return nil, gorm.ErrRecordNotFound
			},
		}
		svc := setupCommentService(commentRepo, nil)

		err := svc.Delete(999, "user-uuid", false)
		require.Error(t, err)
		assert.Equal(t, "comentário não encontrado", err.Error())
	})
}

func TestCommentService_Update(t *testing.T) {
	t.Run("dono pode editar o comentário", func(t *testing.T) {
		c := model.Comment{ID: 1, Description: "Antes", UserID: "user-uuid"}
		updated := model.Comment{ID: 1, Description: "Depois", UserID: "user-uuid"}

		commentRepo := &mocks.CommentRepoMock{
			FindByIDFn: func(id uint) (*model.Comment, error) {
				if id == 1 && c.Description == "Antes" {
					return &c, nil
				}
				return &updated, nil
			},
			UpdateFn: func(comment *model.Comment) error {
				c.Description = comment.Description // simula a persistência
				return nil
			},
		}
		svc := setupCommentService(commentRepo, nil)

		resp, err := svc.Update(1, "user-uuid", service.UpdateCommentInput{Description: "Depois"})

		require.NoError(t, err)
		assert.Equal(t, "Depois", resp.Description)
	})

	t.Run("outro usuário não pode editar", func(t *testing.T) {
		c := model.Comment{ID: 1, UserID: "dono-uuid"}
		commentRepo := &mocks.CommentRepoMock{
			FindByIDFn: func(id uint) (*model.Comment, error) { return &c, nil },
		}
		svc := setupCommentService(commentRepo, nil)

		_, err := svc.Update(1, "invasor", service.UpdateCommentInput{Description: "Hackeado"})
		require.Error(t, err)
		assert.Equal(t, "sem permissão para editar este comentário", err.Error())
	})
}
