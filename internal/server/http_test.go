package server_test

import (
	"context"
	"errors"
	"net/http"
	"strings"
	"testing"

	"panda-cooking-go-api/internal/handler"
	"panda-cooking-go-api/internal/model"
	"panda-cooking-go-api/internal/repository"
	"panda-cooking-go-api/internal/server"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gorm.io/gorm"
)

func validRecipe() map[string]any {
	return map[string]any{
		"name":         "Bolo de fubá",
		"description":  "Bolo simples de fubá com erva-doce.",
		"time":         "50 minutos",
		"portions":     8,
		"category_id":  12,
		"images":       []map[string]string{{"url": "https://img.exemplo.com/bolo.jpg"}},
		"ingredients":  []map[string]string{{"name": "Fubá", "amount": "2 xícaras"}},
		"preparations": []map[string]string{{"description": "Misture tudo e asse."}},
	}
}

func storedRecipe() *model.Recipe {
	return &model.Recipe{
		ID: recipeID, Name: "Bolo de fubá", Description: "Bolo simples de fubá com erva-doce.", Time: "50 minutos",
		Portions: 8, UserID: userID, CategoryID: 12,
		User:     model.User{ID: userID, Name: "Maria", Email: "maria@email.com", PasswordHash: "hash-secreto"},
		Category: model.Category{ID: 12, Name: "Pães e Bolos"},
	}
}

func (e *env) withRecipe() {
	e.recipes.FindByIDFn = func(_ context.Context, id string) (*model.Recipe, error) {
		if id != recipeID {
			return nil, gorm.ErrRecordNotFound
		}
		return storedRecipe(), nil
	}
	e.recipes.ExistsFn = func(_ context.Context, id string) (bool, error) { return id == recipeID, nil }
	e.categories.ExistsFn = func(_ context.Context, id uint) (bool, error) { return id == 12, nil }
}

func TestValidation(t *testing.T) {
	e := newEnv(t)

	tests := []struct {
		name       string
		body       any
		wantStatus int
		wantError  string
		wantFields map[string]string
	}{
		{
			name:       "campos obrigatórios e listas vazias",
			body:       map[string]any{"name": "Bo", "description": "curta", "portions": 0, "ingredients": []any{}},
			wantStatus: http.StatusUnprocessableEntity,
			wantError:  "dados inválidos",
			wantFields: map[string]string{
				"name":         "use pelo menos 3 caracteres",
				"description":  "use pelo menos 10 caracteres",
				"time":         "campo obrigatório",
				"portions":     "campo obrigatório",
				"category_id":  "campo obrigatório",
				"ingredients":  "adicione pelo menos um item",
				"preparations": "adicione pelo menos um item",
			},
		},
		{
			name: "erro dentro da lista aponta o item",
			body: func() map[string]any {
				r := validRecipe()
				r["images"] = []map[string]string{{"url": "http://img.exemplo.com/bolo.jpg"}}
				r["ingredients"] = []map[string]string{{"name": "Fubá", "amount": ""}}
				r["portions"] = 101
				return r
			}(),
			wantStatus: http.StatusUnprocessableEntity,
			wantError:  "dados inválidos",
			wantFields: map[string]string{
				"images[0].url":         "use um link que comece com https://",
				"ingredients[0].amount": "campo obrigatório",
				"portions":              "o máximo é 100",
			},
		},
		{
			name:       "só espaços conta como vazio",
			body:       func() map[string]any { r := validRecipe(); r["time"] = "   "; return r }(),
			wantStatus: http.StatusUnprocessableEntity,
			wantError:  "dados inválidos",
			wantFields: map[string]string{"time": "campo obrigatório"},
		},
		{
			name:       "campo desconhecido é recusado",
			body:       func() map[string]any { r := validRecipe(); r["user_id"] = otherID; return r }(),
			wantStatus: http.StatusUnprocessableEntity,
			wantError:  "dados inválidos",
			wantFields: map[string]string{"user_id": "campo não permitido"},
		},
		{
			name:       "tipo errado",
			body:       `{"name": 123}`,
			wantStatus: http.StatusUnprocessableEntity,
			wantError:  "dados inválidos",
			wantFields: map[string]string{"name": "tipo de valor inválido"},
		},
		{name: "JSON malformado", body: `{"name": `, wantStatus: http.StatusBadRequest, wantError: "JSON inválido"},
		{name: "dois JSON seguidos", body: `{} {}`, wantStatus: http.StatusBadRequest, wantError: "JSON inválido"},
		{name: "corpo vazio", body: "", wantStatus: http.StatusBadRequest, wantError: "corpo da requisição vazio"},
		{
			name:       "corpo maior que 1 MB",
			body:       `{"name": "` + strings.Repeat("a", server.MaxBodyBytes) + `"}`,
			wantStatus: http.StatusRequestEntityTooLarge,
			wantError:  "corpo da requisição grande demais",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := e.as(userID, false, http.MethodPost, "/api/recipes", tt.body)

			require.Equal(t, tt.wantStatus, w.Code, w.Body.String())
			body := decode[errorBody](t, w)
			assert.Equal(t, tt.wantError, body.Error)
			if tt.wantFields != nil {
				assert.Equal(t, tt.wantFields, body.Fields)
			}
		})
	}
}

func TestSignup(t *testing.T) {
	t.Run("cria a conta sem devolver a senha", func(t *testing.T) {
		e := newEnv(t)
		e.users.CreateFn = func(_ context.Context, u *model.User) error {
			assert.Equal(t, "maria@email.com", u.Email)
			u.ID = userID
			return nil
		}

		w := e.do(http.MethodPost, "/api/users", map[string]string{"name": "Maria", "email": "Maria@Email.com", "password": "panelas-de-barro"})

		require.Equal(t, http.StatusCreated, w.Code, w.Body.String())
		assert.NotContains(t, w.Body.String(), "password")
		assert.NotContains(t, w.Body.String(), "panelas-de-barro")
	})

	t.Run("e-mail já cadastrado vem no campo", func(t *testing.T) {
		e := newEnv(t)
		e.users.CreateFn = func(context.Context, *model.User) error { return gorm.ErrDuplicatedKey }

		w := e.do(http.MethodPost, "/api/users", map[string]string{"name": "Maria", "email": "maria@email.com", "password": "panelas-de-barro"})

		require.Equal(t, http.StatusConflict, w.Code)
		assert.Equal(t, errorBody{Error: "e-mail já cadastrado", Fields: map[string]string{"email": "e-mail já cadastrado"}}, decode[errorBody](t, w))
	})

	t.Run("senha fraca vem no campo", func(t *testing.T) {
		w := newEnv(t).do(http.MethodPost, "/api/users", map[string]string{"name": "Maria", "email": "maria@email.com", "password": "123"})

		require.Equal(t, http.StatusUnprocessableEntity, w.Code)
		assert.Equal(t, map[string]string{"password": "use pelo menos 10 caracteres"}, decode[errorBody](t, w).Fields)
	})

	t.Run("não dá para se cadastrar como admin", func(t *testing.T) {
		w := newEnv(t).do(http.MethodPost, "/api/users", map[string]any{"name": "Maria", "email": "maria@email.com", "password": "panelas-de-barro", "is_adm": true})

		require.Equal(t, http.StatusUnprocessableEntity, w.Code)
		assert.Equal(t, map[string]string{"is_adm": "campo não permitido"}, decode[errorBody](t, w).Fields)
	})

	t.Run("limite de cadastros por IP", func(t *testing.T) {
		e := newEnv(t)
		// o limite conta toda tentativa, até a que nem passa na validação
		for range server.SignupLimit {
			require.Equal(t, http.StatusUnprocessableEntity, e.do(http.MethodPost, "/api/users", map[string]string{}).Code)
		}

		assert.Equal(t, http.StatusTooManyRequests, e.do(http.MethodPost, "/api/users", map[string]string{}).Code)
	})
}

func TestProfile(t *testing.T) {
	t.Run("não devolve o hash da senha", func(t *testing.T) {
		e := newEnv(t)
		e.users.FindByIDFn = func(context.Context, string) (*model.User, error) {
			return &model.User{ID: userID, Name: "Maria", Email: "maria@email.com", PasswordHash: "hash-secreto"}, nil
		}

		w := e.as(userID, false, http.MethodGet, "/api/users/profile", nil)

		require.Equal(t, http.StatusOK, w.Code)
		assert.NotContains(t, w.Body.String(), "hash-secreto")
		assert.Equal(t, "maria@email.com", decode[map[string]any](t, w)["email"])
	})

	t.Run("foto precisa ser https", func(t *testing.T) {
		w := newEnv(t).as(userID, false, http.MethodPatch, "/api/users/profile", map[string]string{"image_profile": "javascript:alert(1)"})

		require.Equal(t, http.StatusUnprocessableEntity, w.Code)
		assert.Equal(t, map[string]string{"image_profile": "use um link que comece com https://"}, decode[errorBody](t, w).Fields)
	})

	t.Run("apagar a conta apaga os cookies", func(t *testing.T) {
		e := newEnv(t)
		e.users.FindByIDFn = func(context.Context, string) (*model.User, error) { return &model.User{ID: userID}, nil }
		e.users.DeleteFn = func(context.Context, string) error { return nil }

		w := e.as(userID, false, http.MethodDelete, "/api/users/profile", nil)

		require.Equal(t, http.StatusNoContent, w.Code)
		assert.Negative(t, cookie(w, handler.RefreshCookie).MaxAge)
	})
}

func TestRecipes(t *testing.T) {
	t.Run("listagem pública com filtros", func(t *testing.T) {
		e := newEnv(t)
		e.recipes.ListFn = func(_ context.Context, f repository.RecipeFilter, p repository.Page) ([]model.Recipe, int64, error) {
			assert.Equal(t, repository.RecipeFilter{Search: "pão de queijo", CategoryID: 12, UserID: userID}, f)
			assert.Equal(t, repository.Page{Number: 2, Size: 6}, p)
			return []model.Recipe{*storedRecipe()}, 7, nil
		}

		w := e.do(http.MethodGet, "/api/recipes?search=+p%C3%A3o++de+queijo&category_id=12&user_id="+userID+"&page=2&per_page=6", nil)

		require.Equal(t, http.StatusOK, w.Code, w.Body.String())
		page := decode[map[string]any](t, w)
		assert.InDelta(t, 2, page["total_pages"], 0)
		assert.NotContains(t, w.Body.String(), "maria@email.com", "o e-mail do autor não é público")
	})

	t.Run("parâmetros inválidos da listagem", func(t *testing.T) {
		w := newEnv(t).do(http.MethodGet, "/api/recipes?per_page=500&user_id=abc&search="+strings.Repeat("a", 101), nil)

		require.Equal(t, http.StatusUnprocessableEntity, w.Code)
		assert.Equal(t, map[string]string{
			"per_page": "o máximo é 50",
			"user_id":  "id inválido",
			"search":   "use no máximo 100 caracteres",
		}, decode[errorBody](t, w).Fields)
	})

	t.Run("número inválido na query", func(t *testing.T) {
		w := newEnv(t).do(http.MethodGet, "/api/recipes?page=abc", nil)

		assert.Equal(t, http.StatusUnprocessableEntity, w.Code)
	})

	t.Run("cria logado", func(t *testing.T) {
		e := newEnv(t)
		e.withRecipe()
		e.recipes.CreateFn = func(_ context.Context, r *model.Recipe) error {
			assert.Equal(t, userID, r.UserID)
			r.ID = recipeID
			return nil
		}

		w := e.as(userID, false, http.MethodPost, "/api/recipes", validRecipe())

		require.Equal(t, http.StatusCreated, w.Code, w.Body.String())
		assert.Equal(t, recipeID, decode[map[string]any](t, w)["id"])
	})

	t.Run("criar sem login é 401", func(t *testing.T) {
		assert.Equal(t, http.StatusUnauthorized, newEnv(t).do(http.MethodPost, "/api/recipes", validRecipe()).Code)
	})

	t.Run("categoria inexistente vem no campo", func(t *testing.T) {
		e := newEnv(t)
		e.withRecipe()
		r := validRecipe()
		r["category_id"] = 99

		w := e.as(userID, false, http.MethodPost, "/api/recipes", r)

		require.Equal(t, http.StatusUnprocessableEntity, w.Code)
		assert.Equal(t, map[string]string{"category_id": "categoria não encontrada"}, decode[errorBody](t, w).Fields)
	})

	t.Run("editar receita de outra pessoa é 403", func(t *testing.T) {
		e := newEnv(t)
		e.withRecipe()

		w := e.as(otherID, false, http.MethodPut, "/api/recipes/"+recipeID, validRecipe())

		assert.Equal(t, http.StatusForbidden, w.Code)
		assert.Equal(t, "sem permissão para editar esta receita", decode[errorBody](t, w).Error)
	})

	t.Run("admin também não edita receita alheia", func(t *testing.T) {
		e := newEnv(t)
		e.withRecipe()

		assert.Equal(t, http.StatusForbidden, e.as(otherID, true, http.MethodDelete, "/api/recipes/"+recipeID, nil).Code)
	})

	t.Run("receita que não existe é 404", func(t *testing.T) {
		e := newEnv(t)
		e.withRecipe()

		w := e.do(http.MethodGet, "/api/recipes/nao-existe", nil)

		assert.Equal(t, http.StatusNotFound, w.Code)
		assert.Equal(t, "receita não encontrada", decode[errorBody](t, w).Error)
	})

	t.Run("id de item que não é número é 404", func(t *testing.T) {
		e := newEnv(t)
		e.withRecipe()

		w := e.as(userID, false, http.MethodDelete, "/api/recipes/"+recipeID+"/images/abc", nil)

		assert.Equal(t, http.StatusNotFound, w.Code)
		assert.Equal(t, "imagem não encontrada", decode[errorBody](t, w).Error)
	})
}

func TestComments(t *testing.T) {
	t.Run("lista pública e paginada", func(t *testing.T) {
		e := newEnv(t)
		e.withRecipe()
		e.comments.ListByRecipeFn = func(_ context.Context, id string, p repository.Page) ([]model.Comment, int64, error) {
			assert.Equal(t, repository.Page{Number: 1, Size: 10}, p)
			return []model.Comment{{ID: 1, Description: "Ótima!", RecipeID: id, User: model.User{ID: otherID, Name: "João", Email: "joao@email.com"}}}, 1, nil
		}

		w := e.do(http.MethodGet, "/api/recipes/"+recipeID+"/comments", nil)

		require.Equal(t, http.StatusOK, w.Code)
		assert.Contains(t, w.Body.String(), `"name":"João"`)
		assert.NotContains(t, w.Body.String(), "joao@email.com")
	})

	t.Run("comentário vazio é 422", func(t *testing.T) {
		e := newEnv(t)
		e.withRecipe()

		w := e.as(userID, false, http.MethodPost, "/api/recipes/"+recipeID+"/comments", map[string]string{"description": "   "})

		require.Equal(t, http.StatusUnprocessableEntity, w.Code)
		assert.Equal(t, map[string]string{"description": "campo obrigatório"}, decode[errorBody](t, w).Fields)
	})

	t.Run("admin apaga comentário de qualquer um", func(t *testing.T) {
		e := newEnv(t)
		e.comments.FindByIDFn = func(context.Context, uint) (*model.Comment, error) {
			return &model.Comment{ID: 3, UserID: userID, RecipeID: recipeID}, nil
		}
		e.comments.DeleteFn = func(context.Context, uint) error { return nil }

		assert.Equal(t, http.StatusForbidden, e.as(otherID, false, http.MethodDelete, "/api/comments/3", nil).Code)
		assert.Equal(t, http.StatusNoContent, e.as(otherID, true, http.MethodDelete, "/api/comments/3", nil).Code)
	})
}

func TestFavorites(t *testing.T) {
	e := newEnv(t)
	e.withRecipe()
	e.favorites.CreateFn = func(context.Context, *model.FavoriteRecipe) error { return nil }
	e.favorites.ExistsFn = func(context.Context, string, string) (bool, error) { return true, nil }
	e.favorites.DeleteFn = func(context.Context, string, string) (bool, error) { return false, nil }

	w := e.as(userID, false, http.MethodPost, "/api/favorites/"+recipeID, nil)
	require.Equal(t, http.StatusCreated, w.Code)
	assert.JSONEq(t, `{"favorite": true}`, w.Body.String())

	w = e.as(userID, false, http.MethodGet, "/api/favorites/"+recipeID, nil)
	require.Equal(t, http.StatusOK, w.Code)
	assert.JSONEq(t, `{"favorite": true}`, w.Body.String())

	w = e.as(userID, false, http.MethodDelete, "/api/favorites/"+recipeID, nil)
	assert.Equal(t, http.StatusNotFound, w.Code)

	w = e.as(userID, false, http.MethodPost, "/api/favorites/nao-existe", nil)
	assert.Equal(t, http.StatusNotFound, w.Code)
}

func TestErrorsDoNotLeakInternals(t *testing.T) {
	e := newEnv(t)
	e.recipes.ListFn = func(context.Context, repository.RecipeFilter, repository.Page) ([]model.Recipe, int64, error) {
		return nil, 0, errors.New(`pq: relation "recipes" does not exist at 10.0.0.5:5432`)
	}

	w := e.do(http.MethodGet, "/api/recipes", nil)

	require.Equal(t, http.StatusInternalServerError, w.Code)
	assert.Equal(t, "erro interno, tente novamente mais tarde", decode[errorBody](t, w).Error)
	assert.NotContains(t, w.Body.String(), "10.0.0.5")
	assert.Contains(t, e.logs.String(), "10.0.0.5", "o detalhe fica no log")
}

func TestPanicBecomesJSON500(t *testing.T) {
	e := newEnv(t)
	e.recipes.ListFn = func(context.Context, repository.RecipeFilter, repository.Page) ([]model.Recipe, int64, error) {
		panic("bug")
	}

	w := e.do(http.MethodGet, "/api/recipes", nil)

	require.Equal(t, http.StatusInternalServerError, w.Code)
	assert.Equal(t, "erro interno, tente novamente mais tarde", decode[errorBody](t, w).Error)
	assert.Contains(t, e.logs.String(), "panic na requisição")
}

func TestUnknownRoutes(t *testing.T) {
	e := newEnv(t)

	w := e.do(http.MethodGet, "/api/nao-existe", nil)
	require.Equal(t, http.StatusNotFound, w.Code)
	assert.Equal(t, "rota não encontrada", decode[errorBody](t, w).Error)

	w = e.do(http.MethodPut, "/api/categories", map[string]string{})
	require.Equal(t, http.StatusMethodNotAllowed, w.Code)
	assert.Equal(t, "método não permitido nesta rota", decode[errorBody](t, w).Error)
}

func TestSecurityHeaders(t *testing.T) {
	e := newEnv(t)
	e.categories.FindAllFn = func(context.Context) ([]model.Category, error) {
		return []model.Category{{ID: 1, Name: "Bebidas"}}, nil
	}

	w := e.do(http.MethodGet, "/api/categories", nil)

	require.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, "nosniff", w.Header().Get("X-Content-Type-Options"))
	assert.Equal(t, "DENY", w.Header().Get("X-Frame-Options"))
	assert.Equal(t, "no-referrer", w.Header().Get("Referrer-Policy"))
	assert.Equal(t, "default-src 'none'; frame-ancestors 'none'", w.Header().Get("Content-Security-Policy"))
	assert.Equal(t, "no-store", w.Header().Get("Cache-Control"))
	assert.Empty(t, w.Header().Get("Strict-Transport-Security"), "HSTS só em produção")
	assert.Empty(t, w.Header().Get("Server"))
}

func TestRequestID(t *testing.T) {
	e := newEnv(t)

	generated := e.do(http.MethodGet, "/api/nao-existe", nil).Header().Get("X-Request-ID")
	assert.Regexp(t, `^[0-9a-f]{24}$`, generated)

	req := request(http.MethodGet, "/api/nao-existe", nil)
	req.Header.Set("X-Request-ID", "front-123")
	assert.Equal(t, "front-123", e.serve(req).Header().Get("X-Request-ID"))

	req = request(http.MethodGet, "/api/nao-existe", nil)
	req.Header.Set("X-Request-ID", "quebra\nde linha")
	assert.NotEqual(t, "quebra\nde linha", e.serve(req).Header().Get("X-Request-ID"))
}

func TestCORS(t *testing.T) {
	e := newEnv(t)

	req := request(http.MethodOptions, "/api/recipes", nil)
	req.Header.Set("Origin", frontURL)
	req.Header.Set("Access-Control-Request-Method", http.MethodPost)
	req.Header.Set("Access-Control-Request-Headers", "authorization,content-type")
	w := e.serve(req)

	assert.Equal(t, http.StatusNoContent, w.Code)
	assert.Equal(t, frontURL, w.Header().Get("Access-Control-Allow-Origin"))
	assert.Equal(t, "true", w.Header().Get("Access-Control-Allow-Credentials"))

	req = request(http.MethodOptions, "/api/recipes", nil)
	req.Header.Set("Origin", "https://outro-site.exemplo")
	req.Header.Set("Access-Control-Request-Method", http.MethodPost)
	w = e.serve(req)

	assert.Empty(t, w.Header().Get("Access-Control-Allow-Origin"))
}

func TestHealth(t *testing.T) {
	e := newEnv(t)

	w := e.do(http.MethodGet, "/health", nil)
	assert.Equal(t, http.StatusOK, w.Code)
	assert.JSONEq(t, `{"status":"ok"}`, w.Body.String())
	assert.NotContains(t, e.logs.String(), "/health", "healthcheck que passa não enche o log")

	e.pingErr = errors.New("banco fora")
	w = e.do(http.MethodGet, "/health", nil)
	assert.Equal(t, http.StatusServiceUnavailable, w.Code)
	assert.JSONEq(t, `{"status":"unavailable"}`, w.Body.String())
}

func TestGlobalRateLimit(t *testing.T) {
	e := newEnv(t)
	e.categories.FindAllFn = func(context.Context) ([]model.Category, error) { return nil, nil }

	for range server.GlobalLimit {
		require.Equal(t, http.StatusOK, e.do(http.MethodGet, "/api/categories", nil).Code)
	}
	w := e.do(http.MethodGet, "/api/categories", nil)

	assert.Equal(t, http.StatusTooManyRequests, w.Code)
	assert.Equal(t, "muitas requisições, tente novamente mais tarde", decode[errorBody](t, w).Error)
	// o healthcheck não entra no limite
	assert.Equal(t, http.StatusOK, e.do(http.MethodGet, "/health", nil).Code)
}

// O limite por IP usa o IP da conexão. O X-Forwarded-For só vale quando a
// conexão vem de um proxy listado em TRUSTED_PROXIES.
func TestClientIP(t *testing.T) {
	// corpo vazio: responde 422 sem gastar tempo com bcrypt, mas conta no limite
	signup := func(e *env, remote, forwarded string) int {
		req := request(http.MethodPost, "/api/users", map[string]string{})
		req.RemoteAddr = remote + ":5000"
		if forwarded != "" {
			req.Header.Set("X-Forwarded-For", forwarded)
		}
		return e.serve(req).Code
	}
	exhaust := func(e *env, remote, forwarded string) {
		for range server.SignupLimit {
			require.Equal(t, http.StatusUnprocessableEntity, signup(e, remote, forwarded))
		}
	}

	t.Run("X-Forwarded-For falso não burla o limite", func(t *testing.T) {
		e := newEnv(t)
		exhaust(e, "203.0.113.7", "")

		assert.Equal(t, http.StatusTooManyRequests, signup(e, "203.0.113.7", "198.51.100.1"))
		assert.Equal(t, http.StatusUnprocessableEntity, signup(e, "203.0.113.8", ""), "outro IP não é afetado")
	})

	t.Run("atrás de proxy confiável vale o IP que ele informa", func(t *testing.T) {
		cfg := defaultConfig()
		cfg.TrustedProxies = []string{"10.0.0.0/8"}
		e := newEnvWith(t, cfg)
		exhaust(e, "10.0.0.2", "198.51.100.1")

		assert.Equal(t, http.StatusTooManyRequests, signup(e, "10.0.0.2", "198.51.100.1"))
		assert.Equal(t, http.StatusUnprocessableEntity, signup(e, "10.0.0.2", "198.51.100.2"), "outro cliente atrás do mesmo proxy")
	})
}

func TestNewRejectsInvalidTrustedOrigin(t *testing.T) {
	cfg := defaultConfig()
	cfg.CORSOrigins = []string{"não é origem"}

	_, err := server.New(cfg, server.Deps{})

	assert.Error(t, err)
}
