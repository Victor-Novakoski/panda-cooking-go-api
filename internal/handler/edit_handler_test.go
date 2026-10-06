package handler_test

import (
	"net/http"
	"testing"

	"panda-cooking-go-api/internal/handler"
	"panda-cooking-go-api/internal/handler/testhelper"
	"panda-cooking-go-api/internal/model"
	"panda-cooking-go-api/internal/service"
	"panda-cooking-go-api/internal/service/mocks"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func validRecipeBody() gin.H {
	return gin.H{
		"name": "Panqueca", "description": "Fofinha", "time": "20 minutos",
		"portions": 2, "category_id": 1,
		"images":       []gin.H{{"url": "https://exemplo.com/p.jpg"}},
		"ingredients":  []gin.H{{"name": "farinha", "amount": "1 xícara"}},
		"preparations": []gin.H{{"description": "Misture"}},
	}
}

func TestRecipeHandler_Replace(t *testing.T) {
	t.Run("dono recebe 200 com a receita nova", func(t *testing.T) {
		r := baseRecipeModel()
		repo := &mocks.RecipeRepoMock{
			FindByIDFn: func(string) (*model.Recipe, error) { return &r, nil },
			ReplaceFn:  func(*model.Recipe) error { return nil },
		}

		w := testhelper.Execute(setupRecipeRouter(repo),
			testhelper.NewRequest(http.MethodPut, "/recipes/recipe-uuid", validRecipeBody()))

		assert.Equal(t, http.StatusOK, w.Code)
	})

	t.Run("outro usuário recebe 403", func(t *testing.T) {
		r := baseRecipeModel()
		r.UserID = "outro-uuid"
		repo := &mocks.RecipeRepoMock{FindByIDFn: func(string) (*model.Recipe, error) { return &r, nil }}

		w := testhelper.Execute(setupRecipeRouter(repo),
			testhelper.NewRequest(http.MethodPut, "/recipes/recipe-uuid", validRecipeBody()))

		assert.Equal(t, http.StatusForbidden, w.Code)
	})

	t.Run("campo obrigatório faltando é 400", func(t *testing.T) {
		body := validRecipeBody()
		delete(body, "name")

		w := testhelper.Execute(setupRecipeRouter(&mocks.RecipeRepoMock{}),
			testhelper.NewRequest(http.MethodPut, "/recipes/recipe-uuid", body))

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("item da lista também é validado", func(t *testing.T) {
		body := validRecipeBody()
		body["images"] = []gin.H{{"url": "não é url"}}

		w := testhelper.Execute(setupRecipeRouter(&mocks.RecipeRepoMock{}),
			testhelper.NewRequest(http.MethodPut, "/recipes/recipe-uuid", body))

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("sem login é 401", func(t *testing.T) {
		h := handler.NewRecipeHandler(service.NewRecipeService(&mocks.RecipeRepoMock{}))
		router := gin.New()
		h.RegisterRoutes(router.Group("/recipes"), realAuth())

		w := testhelper.Execute(router,
			testhelper.NewRequest(http.MethodPut, "/recipes/recipe-uuid", validRecipeBody()))

		assert.Equal(t, http.StatusUnauthorized, w.Code)
	})
}

func TestCommentHandler_GetByRecipe(t *testing.T) {
	r := baseRecipeModel()
	commentRepo := &mocks.CommentRepoMock{
		FindByRecipeFn: func(string) ([]model.Comment, error) {
			return []model.Comment{{
				ID: 1, Description: "Muito bom", RecipeID: "recipe-uuid", UserID: "user-uuid",
				User: model.User{ID: "user-uuid", Name: "Maria", Email: "maria@email.com"},
			}}, nil
		},
	}
	recipeRepo := &mocks.RecipeRepoMock{FindByIDFn: func(string) (*model.Recipe, error) { return &r, nil }}
	h := handler.NewCommentHandler(service.NewCommentService(commentRepo, recipeRepo))
	router := gin.New()
	h.RegisterRecipeRoutes(router.Group("/recipes"))

	w := testhelper.Execute(router, testhelper.NewRequest(http.MethodGet, "/recipes/recipe-uuid/comments", nil))

	require.Equal(t, http.StatusOK, w.Code, "lista é pública, sem login")
	assert.Contains(t, w.Body.String(), `"name":"Maria"`)
	assert.NotContains(t, w.Body.String(), "maria@email.com", "e-mail do autor não sai")
}

func TestUserHandler_UpdateProfile(t *testing.T) {
	setup := func() *gin.Engine {
		repo := &mocks.UserRepoMock{
			FindByIDFn: func(id string) (*model.User, error) {
				return &model.User{ID: id, Name: "Maria", ImageProfile: "https://exemplo.com/maria.jpg"}, nil
			},
			UpdateFn: func(*model.User) error { return nil },
		}
		h := handler.NewUserHandler(service.NewUserService(repo, "test-secret"))
		router := gin.New()
		h.RegisterRoutes(router.Group("/users"), fakeAuth("user-uuid", false))
		return router
	}

	tests := []struct {
		name   string
		body   gin.H
		status int
	}{
		{"troca nome e foto", gin.H{"name": "Maria Silva", "image_profile": "https://exemplo.com/nova.jpg"}, http.StatusOK},
		{"foto vazia remove a foto", gin.H{"image_profile": ""}, http.StatusOK},
		{"foto que não é URL é 400", gin.H{"image_profile": "foto.jpg"}, http.StatusBadRequest},
		{"nome vazio é 400", gin.H{"name": ""}, http.StatusBadRequest},
		{"nome só com espaços é 400", gin.H{"name": "   "}, http.StatusBadRequest},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := testhelper.Execute(setup(), testhelper.NewRequest(http.MethodPatch, "/users/profile", tt.body))
			assert.Equal(t, tt.status, w.Code, w.Body.String())
		})
	}
}
