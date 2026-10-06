package testhelper

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"

	"github.com/gin-gonic/gin"
)

func init() {
	gin.SetMode(gin.TestMode)
}

// NewRequest cria uma request HTTP pronta para uso nos testes.
func NewRequest(method, url string, body any) *http.Request {
	var b bytes.Buffer
	if body != nil {
		json.NewEncoder(&b).Encode(body)
	}
	req, _ := http.NewRequest(method, url, &b)
	req.Header.Set("Content-Type", "application/json")
	return req
}

// Execute roda a request no router e retorna o response.
func Execute(router *gin.Engine, req *http.Request) *httptest.ResponseRecorder {
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	return w
}

// Decode desserializa o body do response em v.
func Decode(w *httptest.ResponseRecorder, v any) {
	json.NewDecoder(w.Body).Decode(v)
}
