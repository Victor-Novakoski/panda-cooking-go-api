//go:build integration

package integration_test

import (
	"bytes"
	"encoding/json"
	"io"
	"log/slog"
	"net/http"
	"net/http/cookiejar"
	"net/http/httptest"
	"strings"
	"testing"

	"panda-cooking-go-api/api"
	"panda-cooking-go-api/internal/config"
	"panda-cooking-go-api/internal/server"

	"github.com/stretchr/testify/require"
)

const testSecret = "chave-de-teste-com-mais-de-32-caracteres"

// newAPI sobe a API inteira ligada ao banco de teste.
func newAPI(t *testing.T) *httptest.Server {
	t.Helper()
	cfg := config.Config{Env: "development", SecretKey: testSecret, CORSOrigins: []string{config.DefaultCORSOrigin}}
	router, err := server.New(cfg, server.FromDB(db, slog.New(slog.DiscardHandler), api.Spec))
	require.NoError(t, err)
	srv := httptest.NewServer(router)
	t.Cleanup(srv.Close)
	return srv
}

// client é um navegador simplificado: guarda os cookies e o access token.
type client struct {
	t     *testing.T
	base  string
	http  *http.Client
	token string
}

func newClient(t *testing.T, srv *httptest.Server) *client {
	t.Helper()
	jar, err := cookiejar.New(nil)
	require.NoError(t, err)
	return &client{t: t, base: srv.URL, http: &http.Client{Jar: jar}}
}

// call faz a requisição e devolve status e corpo.
func (c *client) call(method, path string, body any) (int, []byte) {
	c.t.Helper()
	var reader io.Reader
	if body != nil {
		raw, err := json.Marshal(body)
		require.NoError(c.t, err)
		reader = bytes.NewReader(raw)
	}
	req, err := http.NewRequestWithContext(ctx, method, c.base+path, reader)
	require.NoError(c.t, err)
	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}
	if c.token != "" {
		req.Header.Set("Authorization", "Bearer "+c.token)
	}
	resp, err := c.http.Do(req)
	require.NoError(c.t, err)
	defer func() { _ = resp.Body.Close() }()
	data, err := io.ReadAll(resp.Body)
	require.NoError(c.t, err)
	return resp.StatusCode, data
}

// must faz a requisição, confere o status e decodifica o corpo em out (se não for nil).
func (c *client) must(method, path string, body any, want int, out any) {
	c.t.Helper()
	status, data := c.call(method, path, body)
	require.Equalf(c.t, want, status, "%s %s: %s", method, path, data)
	if out != nil {
		require.NoError(c.t, json.Unmarshal(data, out), string(data))
	}
}

func (c *client) signup(name, email, password string) user {
	c.t.Helper()
	var u user
	c.must(http.MethodPost, "/api/users", map[string]string{"name": name, "email": email, "password": password}, http.StatusCreated, &u)
	return u
}

func (c *client) login(email, password string) {
	c.t.Helper()
	var res struct {
		AccessToken string `json:"access_token"`
	}
	c.must(http.MethodPost, "/api/auth/login", map[string]string{"email": email, "password": password}, http.StatusOK, &res)
	c.token = res.AccessToken
}

func (c *client) refresh() int {
	c.t.Helper()
	status, data := c.call(http.MethodPost, "/api/auth/refresh", nil)
	if status == http.StatusOK {
		var res struct {
			AccessToken string `json:"access_token"`
		}
		require.NoError(c.t, json.Unmarshal(data, &res))
		c.token = res.AccessToken
	}
	return status
}

// cookie devolve o valor do cookie que o "navegador" mandaria para path.
func (c *client) cookie(path, name string) string {
	c.t.Helper()
	req, err := http.NewRequest(http.MethodGet, c.base+path, nil)
	require.NoError(c.t, err)
	for _, ck := range c.http.Jar.Cookies(req.URL) {
		if ck.Name == name {
			return ck.Value
		}
	}
	return ""
}

func (c *client) setCookie(path, name, value string) {
	c.t.Helper()
	req, err := http.NewRequest(http.MethodGet, c.base+path, nil)
	require.NoError(c.t, err)
	c.http.Jar.SetCookies(req.URL, []*http.Cookie{{Name: name, Value: value, Path: path}})
}

type user struct {
	ID    string `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
	IsAdm bool   `json:"is_adm"`
}

type recipe struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Time        string `json:"time"`
	Portions    int    `json:"portions"`
	ImageURL    string `json:"image_url"`
	Category    struct {
		ID   uint   `json:"id"`
		Name string `json:"name"`
	} `json:"category"`
	Author struct {
		ID   string `json:"id"`
		Name string `json:"name"`
	} `json:"author"`
	Images []struct {
		ID  uint   `json:"id"`
		URL string `json:"url"`
	} `json:"images"`
	Ingredients []struct {
		ID     uint   `json:"id"`
		Name   string `json:"name"`
		Amount string `json:"amount"`
	} `json:"ingredients"`
	Preparations []struct {
		ID          uint   `json:"id"`
		Description string `json:"description"`
	} `json:"preparations"`
}

type page[T any] struct {
	Items      []T   `json:"items"`
	Page       int   `json:"page"`
	PerPage    int   `json:"per_page"`
	Total      int64 `json:"total"`
	TotalPages int   `json:"total_pages"`
}

type comment struct {
	ID          uint   `json:"id"`
	Description string `json:"description"`
	User        struct {
		Name string `json:"name"`
	} `json:"user"`
}

func names(recipes []recipe) []string {
	out := make([]string, len(recipes))
	for i, r := range recipes {
		out[i] = r.Name
	}
	return out
}

func recipeBody(name, description string, categoryID uint, ingredients ...string) map[string]any {
	items := make([]map[string]string, 0, len(ingredients))
	for _, ing := range ingredients {
		n, amount, _ := strings.Cut(ing, ":")
		items = append(items, map[string]string{"name": n, "amount": amount})
	}
	if len(items) == 0 {
		items = append(items, map[string]string{"name": "Sal", "amount": "a gosto"})
	}
	return map[string]any{
		"name":         name,
		"description":  description,
		"time":         "30 minutos",
		"portions":     4,
		"category_id":  categoryID,
		"images":       []map[string]string{{"url": "https://img.exemplo.com/" + slug(name) + ".jpg"}},
		"ingredients":  items,
		"preparations": []map[string]string{{"description": "Misture tudo."}, {"description": "Sirva."}},
	}
}

// slug deixa só letras sem acento, números e hífen ("Bolo 100% integral" vira "bolo-100-integral").
func slug(s string) string {
	var b strings.Builder
	for _, r := range strings.ToLower(s) {
		switch {
		case r >= 'a' && r <= 'z', r >= '0' && r <= '9':
			b.WriteRune(r)
		case r == ' ':
			b.WriteRune('-')
		}
	}
	return strings.ReplaceAll(b.String(), "--", "-")
}
