package server_test

import (
	"bytes"
	"context"
	"encoding/json"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"panda-cooking-go-api/api"
	"panda-cooking-go-api/internal/config"
	"panda-cooking-go-api/internal/server"
	"panda-cooking-go-api/internal/service/mocks"
	"panda-cooking-go-api/pkg/token"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/require"
)

const (
	secretKey = "chave-de-teste-com-mais-de-32-caracteres"
	userID    = "6f1c2b8e-1d7a-4c55-9a0e-2f4b7c9d1e30"
	otherID   = "0b9e7c41-5a3d-4e21-8f6a-7c2d9b1e4f58"
	recipeID  = "3d6f9a2c-8b1e-4f7a-9c0d-5e2b8a1f6c47"
	frontURL  = "http://localhost:3000"
)

func init() { gin.SetMode(gin.TestMode) }

// env é a API montada com repositórios falsos. Cada teste preenche as
// funções dos mocks que a rota testada usa.
type env struct {
	t          *testing.T
	router     *gin.Engine
	users      *mocks.UserRepoMock
	recipes    *mocks.RecipeRepoMock
	categories *mocks.CategoryRepoMock
	comments   *mocks.CommentRepoMock
	favorites  *mocks.FavoriteRepoMock
	sessions   *mocks.SessionRepoMock
	pingErr    error
	logs       *bytes.Buffer
}

func defaultConfig() config.Config {
	return config.Config{Env: "development", SecretKey: secretKey, CORSOrigins: []string{frontURL}}
}

func newEnv(t *testing.T) *env {
	return newEnvWith(t, defaultConfig())
}

func newEnvWith(t *testing.T, cfg config.Config) *env {
	t.Helper()
	e := &env{
		t:          t,
		users:      &mocks.UserRepoMock{},
		recipes:    &mocks.RecipeRepoMock{},
		categories: &mocks.CategoryRepoMock{},
		comments:   &mocks.CommentRepoMock{},
		favorites:  &mocks.FavoriteRepoMock{},
		sessions:   &mocks.SessionRepoMock{},
		logs:       &bytes.Buffer{},
	}
	log := slog.New(slog.NewTextHandler(e.logs, nil))
	// os handlers registram erro no logger padrão, como no main
	previous := slog.Default()
	slog.SetDefault(log)
	t.Cleanup(func() { slog.SetDefault(previous) })

	router, err := server.New(cfg, server.Deps{
		Users:      e.users,
		Recipes:    e.recipes,
		Categories: e.categories,
		Comments:   e.comments,
		Favorites:  e.favorites,
		Sessions:   e.sessions,
		Ping:       func(context.Context) error { return e.pingErr },
		Log:        log,
		Spec:       api.Spec,
	})
	require.NoError(t, err)
	e.router = router
	return e
}

// request monta uma requisição; body pode ser string (enviada como está) ou
// qualquer valor (vira JSON).
func request(method, path string, body any) *http.Request {
	var reader *strings.Reader
	switch b := body.(type) {
	case nil:
		reader = strings.NewReader("")
	case string:
		reader = strings.NewReader(b)
	default:
		raw, err := json.Marshal(b)
		if err != nil {
			panic(err)
		}
		reader = strings.NewReader(string(raw))
	}
	req := httptest.NewRequest(method, path, reader)
	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}
	return req
}

func (e *env) serve(req *http.Request) *httptest.ResponseRecorder {
	w := httptest.NewRecorder()
	e.router.ServeHTTP(w, req)
	return w
}

func (e *env) do(method, path string, body any) *httptest.ResponseRecorder {
	return e.serve(request(method, path, body))
}

// as faz a requisição logado como o usuário.
func (e *env) as(user string, isAdm bool, method, path string, body any) *httptest.ResponseRecorder {
	req := request(method, path, body)
	req.Header.Set("Authorization", "Bearer "+accessToken(e.t, user, isAdm))
	return e.serve(req)
}

func accessToken(t *testing.T, user string, isAdm bool) string {
	t.Helper()
	tok, err := token.Generate(user, isAdm, secretKey, time.Now(), 15*time.Minute)
	require.NoError(t, err)
	return tok
}

type errorBody struct {
	Error  string            `json:"error"`
	Fields map[string]string `json:"fields"`
}

func decode[T any](t *testing.T, w *httptest.ResponseRecorder) T {
	t.Helper()
	var v T
	require.NoError(t, json.Unmarshal(w.Body.Bytes(), &v), "corpo: %s", w.Body.String())
	return v
}

func cookie(w *httptest.ResponseRecorder, name string) *http.Cookie {
	for _, c := range w.Result().Cookies() {
		if c.Name == name {
			return c
		}
	}
	return nil
}
