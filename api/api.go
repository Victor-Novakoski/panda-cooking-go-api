// Package api guarda a especificação OpenAPI da API (openapi.yaml), embutida
// no binário e servida em /api/openapi.yaml.
package api

import _ "embed"

//go:embed openapi.yaml
var Spec []byte
