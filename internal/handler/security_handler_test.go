package handler_test

import (
	"errors"
	"net/http"
	"testing"
	"time"

	"panda-cooking-go-api/internal/handler"
	"panda-cooking-go-api/internal/handler/testhelper"
	"panda-cooking-go-api/internal/middleware"
	"panda-cooking-go-api/internal/model"
	"panda-cooking-go-api/internal/ratelimit"
	"panda-cooking-go-api/internal/service"
	"panda-cooking-go-api/internal/service/mocks"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

func TestErrorResponses_DoNotLeakInternals(t *testing.T) {
	t.Run("500 responde mensagem genérica, sem o erro do banco", func(t *testing.T) {
		repo := &mocks.RecipeRepoMock{
			FindAllFn: func() ([]model.Recipe, error) {
				return nil, errors.New(`pq: relation "recipes" does not exist at 10.0.0.5:5432`)
			},
		}
		router := setupRecipeRouter(repo)

		w := testhelper.Execute(router, testhelper.NewRequest(http.MethodGet, "/recipes", nil))

		require.Equal(t, http.StatusInternalServerError, w.Code)
		var resp map[string]string
		testhelper.Decode(w, &resp)
		assert.Equal(t, "erro interno, tente novamente mais tarde", resp["error"])
		assert.NotContains(t, resp["error"], "pq:")
	})

	t.Run("erro do banco no login é 500, não 401", func(t *testing.T) {
		r, repo := setupAuthRouter()
		repo.FindByEmailFn = func(string) (*model.User, error) { return nil, errors.New("connection refused") }

		w := testhelper.Execute(r, testhelper.NewRequest(http.MethodPost, "/auth", gin.H{
			"email": "victor@email.com", "password": "123456",
		}))

		assert.Equal(t, http.StatusInternalServerError, w.Code)
		assert.NotContains(t, w.Body.String(), "connection refused")
	})

	t.Run("comentário em receita inexistente é 404", func(t *testing.T) {
		recipeRepo := &mocks.RecipeRepoMock{
			FindByIDFn: func(string) (*model.Recipe, error) { return nil, gorm.ErrRecordNotFound },
		}
		h := handler.NewCommentHandler(service.NewCommentService(&mocks.CommentRepoMock{}, recipeRepo))
		r := gin.New()
		h.RegisterRoutes(r.Group("/comments"), fakeAuth("user-uuid", false))

		w := testhelper.Execute(r, testhelper.NewRequest(http.MethodPost, "/comments", gin.H{
			"description": "Ficou ótima", "recipe_id": "nao-existe",
		}))

		assert.Equal(t, http.StatusNotFound, w.Code)
	})
}

func TestRecipeHandler_ItemOfAnotherRecipe(t *testing.T) {
	r := baseRecipeModel()
	repo := &mocks.RecipeRepoMock{
		FindByIDFn: func(string) (*model.Recipe, error) { return &r, nil },
		FindImageByIDFn: func(id uint) (*model.ImageRecipe, error) {
			return &model.ImageRecipe{ID: id, RecipeID: "receita-de-outra-pessoa"}, nil
		},
		DeleteImageFn: func(uint) error {
			t.Fatal("não deveria apagar a foto de outra receita")
			return nil
		},
	}
	router := setupRecipeRouter(repo)

	w := testhelper.Execute(router, testhelper.NewRequest(http.MethodDelete, "/recipes/recipe-uuid/images/7", nil))

	assert.Equal(t, http.StatusNotFound, w.Code)
}

func TestUserHandler_Register_EmailTaken(t *testing.T) {
	repo := &mocks.UserRepoMock{
		CreateFn: func(*model.User) error { return gorm.ErrDuplicatedKey },
	}
	h := handler.NewUserHandler(service.NewUserService(repo, "test-secret"))
	r := gin.New()
	h.RegisterRoutes(r.Group("/users"), func(c *gin.Context) { c.Next() })

	w := testhelper.Execute(r, testhelper.NewRequest(http.MethodPost, "/users", gin.H{
		"name": "Victor", "email": "victor@email.com", "password": "123456",
	}))

	assert.Equal(t, http.StatusConflict, w.Code)
	var resp map[string]string
	testhelper.Decode(w, &resp)
	assert.Equal(t, "e-mail já cadastrado", resp["error"])
}

func TestAuthHandler_LoginLockout(t *testing.T) {
	hash, _ := bcrypt.GenerateFromPassword([]byte("senha-correta"), bcrypt.MinCost)
	r, repo := setupAuthRouter()
	repo.FindByEmailFn = func(email string) (*model.User, error) {
		return &model.User{ID: "abc", Email: email, Password: string(hash)}, nil
	}
	login := func(password string) int {
		return testhelper.Execute(r, testhelper.NewRequest(http.MethodPost, "/auth", gin.H{
			"email": "victor@email.com", "password": password,
		})).Code
	}

	for range service.MaxLoginFailures {
		require.Equal(t, http.StatusUnauthorized, login("senha-errada"))
	}

	assert.Equal(t, http.StatusTooManyRequests, login("senha-correta"))
}

func TestRateLimitByIP(t *testing.T) {
	newRouter := func() *gin.Engine {
		r := gin.New()
		// igual ao main: sem proxy confiável o X-Forwarded-For é ignorado
		require.NoError(t, r.SetTrustedProxies(nil))
		r.POST("/auth", middleware.RateLimitByIP(ratelimit.New(2, time.Minute)), func(c *gin.Context) {
			c.Status(http.StatusOK)
		})
		return r
	}
	request := func(r *gin.Engine, remoteAddr, forwardedFor string) *http.Request {
		req := testhelper.NewRequest(http.MethodPost, "/auth", nil)
		req.RemoteAddr = remoteAddr
		if forwardedFor != "" {
			req.Header.Set("X-Forwarded-For", forwardedFor)
		}
		return req
	}

	t.Run("passa do limite e recebe 429 com Retry-After", func(t *testing.T) {
		r := newRouter()
		for range 2 {
			require.Equal(t, http.StatusOK, testhelper.Execute(r, request(r, "1.2.3.4:5000", "")).Code)
		}

		w := testhelper.Execute(r, request(r, "1.2.3.4:5000", ""))

		assert.Equal(t, http.StatusTooManyRequests, w.Code)
		assert.Equal(t, "60", w.Header().Get("Retry-After"))
	})

	t.Run("X-Forwarded-For falso não burla o limite", func(t *testing.T) {
		r := newRouter()
		for range 2 {
			testhelper.Execute(r, request(r, "1.2.3.4:5000", ""))
		}

		w := testhelper.Execute(r, request(r, "1.2.3.4:5000", "9.9.9.9"))

		assert.Equal(t, http.StatusTooManyRequests, w.Code)
	})

	t.Run("outro IP não é afetado", func(t *testing.T) {
		r := newRouter()
		for range 3 {
			testhelper.Execute(r, request(r, "1.2.3.4:5000", ""))
		}

		w := testhelper.Execute(r, request(r, "5.6.7.8:5000", ""))

		assert.Equal(t, http.StatusOK, w.Code)
	})
}
