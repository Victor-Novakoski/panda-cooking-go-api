package handler_test

import (
	"fmt"
	"net/http"
	"panda-cooking-go-api/internal/handler"
	"panda-cooking-go-api/internal/handler/testhelper"
	"panda-cooking-go-api/internal/middleware"
	"panda-cooking-go-api/internal/model"
	"panda-cooking-go-api/internal/service"
	"panda-cooking-go-api/internal/service/mocks"
	"panda-cooking-go-api/pkg/token"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gorm.io/gorm"
)

const recipeTestSecret = "test-secret"

// fakeAuth injeta userID e isAdm no contexto sem validar token real.
// Usamos isso nos testes de handler para não depender do pacote token.
func fakeAuth(userID string, isAdm bool) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Set(middleware.UserIDKey, userID)
		c.Set(middleware.IsAdmKey, isAdm)
		c.Next()
	}
}

// realAuth usa o middleware real com token JWT válido.
func realAuth() gin.HandlerFunc {
	return middleware.Auth(recipeTestSecret)
}

func validToken(userID string) string {
	t, _ := token.Generate(userID, false, recipeTestSecret)
	return "Bearer " + t
}

func setupRecipeRouter(repo *mocks.RecipeRepoMock) *gin.Engine {
	svc := service.NewRecipeService(repo)
	h := handler.NewRecipeHandler(svc)
	r := gin.New()
	h.RegisterRoutes(r.Group("/recipes"), fakeAuth("user-uuid", false))
	return r
}

func baseRecipeModel() model.Recipe {
	return model.Recipe{
		ID:       "recipe-uuid",
		Name:     "Panqueca",
		UserID:   "user-uuid",
		Category: model.Category{ID: 1, Name: "Café da manhã"},
	}
}

func TestRecipeHandler_GetAll(t *testing.T) {
	t.Run("retorna 200 com lista de receitas", func(t *testing.T) {
		r := baseRecipeModel()
		repo := &mocks.RecipeRepoMock{
			FindAllFn: func() ([]model.Recipe, error) { return []model.Recipe{r, r}, nil },
		}
		router := setupRecipeRouter(repo)

		w := testhelper.Execute(router, testhelper.NewRequest(http.MethodGet, "/recipes", nil))

		assert.Equal(t, http.StatusOK, w.Code)

		var resp []map[string]any
		testhelper.Decode(w, &resp)
		assert.Len(t, resp, 2)
	})

	t.Run("erro no banco retorna 500", func(t *testing.T) {
		repo := &mocks.RecipeRepoMock{
			FindAllFn: func() ([]model.Recipe, error) { return nil, fmt.Errorf("db down") },
		}
		router := setupRecipeRouter(repo)

		w := testhelper.Execute(router, testhelper.NewRequest(http.MethodGet, "/recipes", nil))

		assert.Equal(t, http.StatusInternalServerError, w.Code)
	})
}

func TestRecipeHandler_GetByID(t *testing.T) {
	t.Run("receita encontrada retorna 200", func(t *testing.T) {
		r := baseRecipeModel()
		repo := &mocks.RecipeRepoMock{
			FindByIDFn: func(id string) (*model.Recipe, error) { return &r, nil },
		}
		router := setupRecipeRouter(repo)

		w := testhelper.Execute(router, testhelper.NewRequest(http.MethodGet, "/recipes/recipe-uuid", nil))

		assert.Equal(t, http.StatusOK, w.Code)

		var resp map[string]any
		testhelper.Decode(w, &resp)
		assert.Equal(t, "Panqueca", resp["name"])
	})

	t.Run("receita não encontrada retorna 404", func(t *testing.T) {
		repo := &mocks.RecipeRepoMock{
			FindByIDFn: func(id string) (*model.Recipe, error) {
				return nil, gorm.ErrRecordNotFound
			},
		}
		router := setupRecipeRouter(repo)

		w := testhelper.Execute(router, testhelper.NewRequest(http.MethodGet, "/recipes/id-invalido", nil))

		assert.Equal(t, http.StatusNotFound, w.Code)
	})
}

func TestRecipeHandler_Create(t *testing.T) {
	t.Run("cria receita válida retorna 201", func(t *testing.T) {
		r := baseRecipeModel()
		repo := &mocks.RecipeRepoMock{
			CreateFn:   func(recipe *model.Recipe) error { recipe.ID = "recipe-uuid"; return nil },
			FindByIDFn: func(id string) (*model.Recipe, error) { return &r, nil },
			FindOrCreateIngredientFn: func(name string) (*model.Ingredient, error) {
				return &model.Ingredient{ID: 1, Name: name}, nil
			},
			AddIngredientFn: func(ir *model.IngredientRecipe) error { return nil },
		}
		router := setupRecipeRouter(repo)

		w := testhelper.Execute(router, testhelper.NewRequest(http.MethodPost, "/recipes", gin.H{
			"name":        "Panqueca",
			"description": "Gostosa",
			"time":        "20min",
			"portions":    4,
			"category_id": 1,
		}))

		assert.Equal(t, http.StatusCreated, w.Code)
	})

	t.Run("body inválido retorna 400", func(t *testing.T) {
		router := setupRecipeRouter(&mocks.RecipeRepoMock{})

		w := testhelper.Execute(router, testhelper.NewRequest(http.MethodPost, "/recipes", gin.H{
			"name": "Sem os campos obrigatórios",
		}))

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})
}

func TestRecipeHandler_Delete(t *testing.T) {
	t.Run("dono deleta e recebe 204", func(t *testing.T) {
		r := baseRecipeModel()
		repo := &mocks.RecipeRepoMock{
			FindByIDFn: func(id string) (*model.Recipe, error) { return &r, nil },
			DeleteFn:   func(id string) error { return nil },
		}
		router := setupRecipeRouter(repo)

		w := testhelper.Execute(router, testhelper.NewRequest(http.MethodDelete, "/recipes/recipe-uuid", nil))

		assert.Equal(t, http.StatusNoContent, w.Code)
	})

	t.Run("não dono recebe 403", func(t *testing.T) {
		r := baseRecipeModel() // UserID = "user-uuid"

		// router com outro usuário autenticado
		svc := service.NewRecipeService(&mocks.RecipeRepoMock{
			FindByIDFn: func(id string) (*model.Recipe, error) { return &r, nil },
		})
		h := handler.NewRecipeHandler(svc)
		router := gin.New()
		h.RegisterRoutes(router.Group("/recipes"), fakeAuth("outro-user", false))

		w := testhelper.Execute(router, testhelper.NewRequest(http.MethodDelete, "/recipes/recipe-uuid", nil))

		assert.Equal(t, http.StatusForbidden, w.Code)
	})
}

func TestRecipeHandler_Auth_Middleware(t *testing.T) {
	t.Run("criar receita sem token retorna 401", func(t *testing.T) {
		// aqui usamos o middleware real, não o fake
		svc := service.NewRecipeService(&mocks.RecipeRepoMock{})
		h := handler.NewRecipeHandler(svc)
		router := gin.New()
		h.RegisterRoutes(router.Group("/recipes"), realAuth())

		w := testhelper.Execute(router, testhelper.NewRequest(http.MethodPost, "/recipes", gin.H{
			"name": "Qualquer",
		}))

		assert.Equal(t, http.StatusUnauthorized, w.Code)
	})

	t.Run("token válido passa pelo middleware", func(t *testing.T) {
		r := baseRecipeModel()
		repo := &mocks.RecipeRepoMock{
			CreateFn:   func(recipe *model.Recipe) error { recipe.ID = "recipe-uuid"; return nil },
			FindByIDFn: func(id string) (*model.Recipe, error) { return &r, nil },
			FindOrCreateIngredientFn: func(name string) (*model.Ingredient, error) {
				return &model.Ingredient{ID: 1, Name: name}, nil
			},
			AddIngredientFn: func(ir *model.IngredientRecipe) error { return nil },
		}
		svc := service.NewRecipeService(repo)
		h := handler.NewRecipeHandler(svc)
		router := gin.New()
		h.RegisterRoutes(router.Group("/recipes"), realAuth())

		req := testhelper.NewRequest(http.MethodPost, "/recipes", gin.H{
			"name": "Panqueca", "description": "Boa", "time": "20min", "portions": 4, "category_id": 1,
		})
		req.Header.Set("Authorization", validToken("user-uuid"))

		w := testhelper.Execute(router, req)

		require.Equal(t, http.StatusCreated, w.Code)
	})
}
