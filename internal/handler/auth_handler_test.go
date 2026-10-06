package handler_test

import (
	"net/http"
	"panda-cooking-go-api/internal/handler"
	"panda-cooking-go-api/internal/handler/testhelper"
	"panda-cooking-go-api/internal/service"
	"panda-cooking-go-api/internal/service/mocks"
	"panda-cooking-go-api/internal/model"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

func setupAuthRouter() (*gin.Engine, *mocks.UserRepoMock) {
	repo := &mocks.UserRepoMock{}
	svc := service.NewUserService(repo, "test-secret")
	h := handler.NewAuthHandler(svc)

	r := gin.New()
	h.RegisterRoutes(r.Group("/auth"))
	return r, repo
}

func TestAuthHandler_Login(t *testing.T) {
	hash, _ := bcrypt.GenerateFromPassword([]byte("123456"), bcrypt.DefaultCost)

	t.Run("login válido retorna 200 e token", func(t *testing.T) {
		r, repo := setupAuthRouter()
		repo.FindByEmailFn = func(email string) (*model.User, error) {
			return &model.User{ID: "abc", Email: email, Password: string(hash)}, nil
		}

		req := testhelper.NewRequest(http.MethodPost, "/auth", gin.H{
			"email":    "victor@email.com",
			"password": "123456",
		})
		w := testhelper.Execute(r, req)

		assert.Equal(t, http.StatusOK, w.Code)

		var resp map[string]string
		testhelper.Decode(w, &resp)
		assert.NotEmpty(t, resp["token"])
	})

	t.Run("senha errada retorna 401", func(t *testing.T) {
		r, repo := setupAuthRouter()
		repo.FindByEmailFn = func(email string) (*model.User, error) {
			return &model.User{Password: string(hash)}, nil
		}

		req := testhelper.NewRequest(http.MethodPost, "/auth", gin.H{
			"email":    "victor@email.com",
			"password": "senha-errada",
		})
		w := testhelper.Execute(r, req)

		assert.Equal(t, http.StatusUnauthorized, w.Code)
	})

	t.Run("email não encontrado retorna 401", func(t *testing.T) {
		r, repo := setupAuthRouter()
		repo.FindByEmailFn = func(email string) (*model.User, error) {
			return nil, gorm.ErrRecordNotFound
		}

		req := testhelper.NewRequest(http.MethodPost, "/auth", gin.H{
			"email":    "naoexiste@email.com",
			"password": "123456",
		})
		w := testhelper.Execute(r, req)

		assert.Equal(t, http.StatusUnauthorized, w.Code)
	})

	t.Run("body inválido retorna 400", func(t *testing.T) {
		r, _ := setupAuthRouter()

		req := testhelper.NewRequest(http.MethodPost, "/auth", gin.H{
			"email": "nao-é-um-email",
		})
		w := testhelper.Execute(r, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})
}

func TestAuthHandler_Register(t *testing.T) {
	r := gin.New()
	repo := &mocks.UserRepoMock{}
	svc := service.NewUserService(repo, "test-secret")
	h := handler.NewUserHandler(svc)
	h.RegisterRoutes(r.Group("/users"), func(c *gin.Context) { c.Next() })

	t.Run("registro válido retorna 201", func(t *testing.T) {
		repo.CreateFn = func(user *model.User) error {
			user.ID = "novo-uuid"
			return nil
		}

		req := testhelper.NewRequest(http.MethodPost, "/users", gin.H{
			"name":     "Victor",
			"email":    "victor@email.com",
			"password": "123456",
		})
		w := testhelper.Execute(r, req)

		require.Equal(t, http.StatusCreated, w.Code)

		var resp map[string]any
		testhelper.Decode(w, &resp)
		assert.Equal(t, "Victor", resp["name"])
		assert.Equal(t, "victor@email.com", resp["email"])
		assert.Nil(t, resp["password"]) // senha nunca volta na resposta
	})

	t.Run("email inválido retorna 400", func(t *testing.T) {
		req := testhelper.NewRequest(http.MethodPost, "/users", gin.H{
			"name":     "Victor",
			"email":    "isso-nao-e-email",
			"password": "123456",
		})
		w := testhelper.Execute(r, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("senha curta retorna 400", func(t *testing.T) {
		req := testhelper.NewRequest(http.MethodPost, "/users", gin.H{
			"name":     "Victor",
			"email":    "v@v.com",
			"password": "123",
		})
		w := testhelper.Execute(r, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})
}
