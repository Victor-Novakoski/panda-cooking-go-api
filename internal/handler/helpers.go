package handler

import (
	"context"
	"encoding/json"
	"errors"
	"io"
	"log/slog"
	"math"
	"net/http"
	"strconv"

	"panda-cooking-go-api/internal/service"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
)

// msgInternal é o que o cliente vê em qualquer erro inesperado.
// O detalhe fica só no log, para não expor banco, SQL ou caminhos internos.
const msgInternal = "erro interno, tente novamente mais tarde"

// statusClientClosed é o 499 do nginx: o cliente fechou a conexão antes da
// resposta. Só aparece no log.
const statusClientClosed = 499

// msgInvalid acompanha o 422: o detalhe de cada campo vai em "fields".
const msgInvalid = "dados inválidos"

var statusByKind = map[service.Kind]int{
	service.KindNotFound:        http.StatusNotFound,
	service.KindForbidden:       http.StatusForbidden,
	service.KindConflict:        http.StatusConflict,
	service.KindUnauthorized:    http.StatusUnauthorized,
	service.KindTooManyRequests: http.StatusTooManyRequests,
	service.KindInvalid:         http.StatusUnprocessableEntity,
}

// ErrorBody é o formato de todo erro da API. Fields vem no 422 (e no 409 de
// e-mail já cadastrado), com a mensagem de cada campo do corpo.
type ErrorBody struct {
	Error  string            `json:"error"`
	Fields map[string]string `json:"fields,omitempty"`
}

// respondError responde um erro do service: erro conhecido sai com o status e
// a mensagem dele; o resto vira 500 com mensagem genérica e vai para o log.
func respondError(c *gin.Context, err error) {
	var appErr *service.Error
	if errors.As(err, &appErr) {
		if status, ok := statusByKind[appErr.Kind]; ok {
			body := ErrorBody{Error: appErr.Message}
			if appErr.Field != "" {
				body.Fields = map[string]string{appErr.Field: appErr.Message}
			}
			if appErr.RetryAfter > 0 {
				c.Header("Retry-After", strconv.Itoa(int(math.Ceil(appErr.RetryAfter.Seconds()))))
			}
			c.AbortWithStatusJSON(status, body)
			return
		}
	}

	ctx := c.Request.Context()
	if errors.Is(err, context.Canceled) {
		// o cliente desistiu (fechou a aba); não há a quem responder
		c.AbortWithStatus(statusClientClosed)
		return
	}
	if errors.Is(err, context.DeadlineExceeded) {
		slog.WarnContext(ctx, "requisição passou do prazo", "path", c.Request.URL.Path, "erro", err)
		c.AbortWithStatusJSON(http.StatusServiceUnavailable, ErrorBody{Error: "a requisição demorou demais, tente novamente"})
		return
	}

	slog.ErrorContext(ctx, "erro interno", "method", c.Request.Method, "route", c.FullPath(), "erro", err)
	c.AbortWithStatusJSON(http.StatusInternalServerError, ErrorBody{Error: msgInternal})
}

// normalizer é implementado pelas entradas que limpam espaços e caixa antes
// da validação (ver service/inputs.go).
type normalizer interface{ Normalize() }

// bindJSON lê o corpo em dst e valida. Responde e devolve false se não der:
// 413 corpo grande demais, 400 JSON malformado, 422 campo inválido ou
// desconhecido (com a mensagem de cada campo).
func bindJSON(c *gin.Context, dst any) bool {
	dec := json.NewDecoder(c.Request.Body)
	dec.DisallowUnknownFields()
	err := dec.Decode(dst)
	if err == nil && dec.More() {
		err = errors.New("conteúdo depois do JSON")
	}
	if err != nil {
		respondDecodeError(c, err)
		return false
	}
	return validate(c, dst)
}

func respondDecodeError(c *gin.Context, err error) {
	var maxErr *http.MaxBytesError
	var typeErr *json.UnmarshalTypeError
	switch {
	case errors.As(err, &maxErr):
		c.AbortWithStatusJSON(http.StatusRequestEntityTooLarge, ErrorBody{Error: "corpo da requisição grande demais"})
	case errors.Is(err, io.EOF):
		c.AbortWithStatusJSON(http.StatusBadRequest, ErrorBody{Error: "corpo da requisição vazio"})
	case errors.As(err, &typeErr) && typeErr.Field != "":
		c.AbortWithStatusJSON(http.StatusUnprocessableEntity, ErrorBody{
			Error:  msgInvalid,
			Fields: map[string]string{typeErr.Field: "tipo de valor inválido"},
		})
	case unknownField(err) != "":
		c.AbortWithStatusJSON(http.StatusUnprocessableEntity, ErrorBody{
			Error:  msgInvalid,
			Fields: map[string]string{unknownField(err): "campo não permitido"},
		})
	default:
		c.AbortWithStatusJSON(http.StatusBadRequest, ErrorBody{Error: "JSON inválido"})
	}
}

// unknownField extrai o nome do campo da mensagem do encoding/json
// (`json: unknown field "x"`), que não tem tipo de erro próprio.
func unknownField(err error) string {
	const prefix = `json: unknown field "`
	msg := err.Error()
	if len(msg) > len(prefix)+1 && msg[:len(prefix)] == prefix {
		return msg[len(prefix) : len(msg)-1]
	}
	return ""
}

// bindQuery lê e valida a query string (filtros e paginação).
func bindQuery(c *gin.Context, dst any) bool {
	if err := binding.MapFormWithTag(dst, c.Request.URL.Query(), "form"); err != nil {
		c.AbortWithStatusJSON(http.StatusUnprocessableEntity, ErrorBody{Error: "parâmetro da busca inválido"})
		return false
	}
	return validate(c, dst)
}

func validate(c *gin.Context, dst any) bool {
	if n, ok := dst.(normalizer); ok {
		n.Normalize()
	}
	if err := binding.Validator.ValidateStruct(dst); err != nil {
		if fields, ok := fieldErrors(err); ok {
			c.AbortWithStatusJSON(http.StatusUnprocessableEntity, ErrorBody{Error: msgInvalid, Fields: fields})
			return false
		}
		c.AbortWithStatusJSON(http.StatusBadRequest, ErrorBody{Error: msgInvalid})
		return false
	}
	return true
}

// idParam lê um id numérico da rota. Id que não é número responde 404,
// como um id que não existe.
func idParam(c *gin.Context, name string, notFound error) (uint, bool) {
	// 63 bits: o maior valor que cabe no BIGINT do Postgres
	v, err := strconv.ParseUint(c.Param(name), 10, 63)
	if err != nil || v == 0 {
		respondError(c, notFound)
		return 0, false
	}
	return uint(v), true
}
