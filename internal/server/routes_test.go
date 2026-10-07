package server_test

import (
	"net/http"
	"regexp"
	"strings"
	"testing"

	"panda-cooking-go-api/api"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.yaml.in/yaml/v3"
)

// A especificação OpenAPI e as rotas registradas precisam bater: rota nova
// sem documentação (ou documentação de rota que não existe) quebra o teste.
func TestOpenAPI_DocumentaTodasAsRotas(t *testing.T) {
	var spec struct {
		Paths map[string]map[string]any `yaml:"paths"`
	}
	require.NoError(t, yaml.Unmarshal(api.Spec, &spec))

	documented := map[string]bool{}
	for path, item := range spec.Paths {
		for method := range item {
			if method != "parameters" {
				documented[strings.ToUpper(method)+" "+path] = true
			}
		}
	}

	registered := map[string]bool{}
	for _, r := range newEnv(t).router.Routes() {
		registered[r.Method+" "+ginParam.ReplaceAllString(r.Path, "{$1}")] = true
	}

	assert.Equal(t, documented, registered)
}

var ginParam = regexp.MustCompile(`:(\w+)`)

func TestOpenAPI_Servida(t *testing.T) {
	w := newEnv(t).do(http.MethodGet, "/api/openapi.yaml", nil)

	require.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, "application/yaml; charset=utf-8", w.Header().Get("Content-Type"))
	assert.Equal(t, api.Spec, w.Body.Bytes())
}
