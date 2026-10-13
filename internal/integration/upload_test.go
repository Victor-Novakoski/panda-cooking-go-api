//go:build integration

package integration_test

import (
	"bytes"
	"encoding/json"
	"io"
	"mime/multipart"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// Começo de um PNG: basta para o formato ser reconhecido.
var pngBytes = []byte("\x89PNG\r\n\x1a\n\x00\x00\x00\rIHDRfoto-de-teste")

func TestUploads(t *testing.T) {
	reset(t)
	c := newClient(t, newAPI(t))
	c.signup("Maria", "maria@email.com", password)

	status, _ := c.upload(pngBytes)
	assert.Equal(t, http.StatusUnauthorized, status, "enviar foto exige login")

	c.login("maria@email.com", password)
	status, body := c.upload([]byte("<svg onload=alert(1)>"))
	assert.Equal(t, http.StatusUnprocessableEntity, status, "%s", body)

	var sent struct {
		URL string `json:"url"`
	}
	status, body = c.upload(pngBytes)
	require.Equal(t, http.StatusCreated, status, "%s", body)
	require.NoError(t, json.Unmarshal(body, &sent))
	assert.Regexp(t, `^/api/uploads/[0-9a-f-]{36}$`, sent.URL)

	// a foto é pública e vem com o tipo conferido pela API
	resp, err := http.Get(c.base + sent.URL)
	require.NoError(t, err)
	defer func() { _ = resp.Body.Close() }()
	data, err := io.ReadAll(resp.Body)
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	assert.Equal(t, "image/png", resp.Header.Get("Content-Type"))
	assert.Contains(t, resp.Header.Get("Content-Security-Policy"), "sandbox")
	assert.Equal(t, pngBytes, data)

	// o endereço serve de foto da receita e do perfil
	body2 := recipeBody("Bolo de foto enviada", "Bolo com a foto enviada do computador.", catPaes)
	body2["images"] = []map[string]string{{"url": sent.URL}, {"url": "https://img.exemplo.com/2.jpg"}}
	var created recipe
	c.must(http.MethodPost, "/api/recipes", body2, http.StatusCreated, &created)
	assert.Equal(t, sent.URL, created.Images[0].URL)
	c.must(http.MethodPatch, "/api/users/profile", map[string]string{"image_profile": sent.URL}, http.StatusOK, nil)

	// caminho parecido que não é de foto enviada continua recusado
	body2["images"] = []map[string]string{{"url": "/api/uploads/../users"}}
	status, _ = c.call(http.MethodPost, "/api/recipes", body2)
	assert.Equal(t, http.StatusUnprocessableEntity, status)

	status, _ = c.call(http.MethodGet, "/api/uploads/00000000-0000-0000-0000-000000000000", nil)
	assert.Equal(t, http.StatusNotFound, status)
	status, _ = c.call(http.MethodGet, "/api/uploads/nao-e-uuid", nil)
	assert.Equal(t, http.StatusNotFound, status)
}

// upload manda a foto como o front manda: multipart, campo "file".
func (c *client) upload(data []byte) (int, []byte) {
	c.t.Helper()
	var buf bytes.Buffer
	w := multipart.NewWriter(&buf)
	part, err := w.CreateFormFile("file", "foto.png")
	require.NoError(c.t, err)
	_, err = part.Write(data)
	require.NoError(c.t, err)
	require.NoError(c.t, w.Close())

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, c.base+"/api/uploads", &buf)
	require.NoError(c.t, err)
	req.Header.Set("Content-Type", w.FormDataContentType())
	if c.token != "" {
		req.Header.Set("Authorization", "Bearer "+c.token)
	}
	resp, err := c.http.Do(req)
	require.NoError(c.t, err)
	defer func() { _ = resp.Body.Close() }()
	out, err := io.ReadAll(resp.Body)
	require.NoError(c.t, err)
	return resp.StatusCode, out
}
