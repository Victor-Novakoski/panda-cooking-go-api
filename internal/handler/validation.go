package handler

import (
	"errors"
	"fmt"
	"net/url"
	"reflect"
	"strings"
	"sync"

	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
)

var setupValidator sync.Once

// SetupValidator configura o validador do Gin: os erros citam o campo pelo
// nome do JSON (ou da query string), não pelo nome da struct em Go, e a
// tag httpsurl confere links de imagem.
func SetupValidator() {
	setupValidator.Do(func() {
		v, ok := binding.Validator.Engine().(*validator.Validate)
		if !ok {
			panic("validador do Gin não é o go-playground/validator")
		}
		v.RegisterTagNameFunc(func(f reflect.StructField) string {
			for _, tag := range []string{"json", "form"} {
				name, _, _ := strings.Cut(f.Tag.Get(tag), ",")
				if name == "-" {
					return ""
				}
				if name != "" {
					return name
				}
			}
			return ""
		})
		if err := v.RegisterValidation("httpsurl", isHTTPSURL); err != nil {
			panic(err)
		}
	})
}

// isHTTPSURL aceita só link https com host e sem usuário ou senha embutidos.
// Vazio passa: "obrigatório" é conferido pelo required. Só https porque o
// site roda em https e imagem http seria bloqueada como conteúdo misto.
func isHTTPSURL(fl validator.FieldLevel) bool {
	s := fl.Field().String()
	if s == "" {
		return true
	}
	if strings.ContainsAny(s, " \t\r\n") {
		return false
	}
	u, err := url.Parse(s)
	return err == nil && u.Scheme == "https" && u.Hostname() != "" && u.User == nil
}

// fieldErrors traduz os erros do validador para {campo: mensagem}.
func fieldErrors(err error) (map[string]string, bool) {
	var verrs validator.ValidationErrors
	if !errors.As(err, &verrs) {
		return nil, false
	}
	fields := make(map[string]string, len(verrs))
	for _, fe := range verrs {
		name := fieldPath(fe.Namespace())
		if _, ok := fields[name]; !ok {
			fields[name] = fieldMessage(fe)
		}
	}
	return fields, true
}

// fieldPath tira o nome da struct do começo: "RecipeInput.ingredients[0].name"
// vira "ingredients[0].name".
func fieldPath(namespace string) string {
	_, path, found := strings.Cut(namespace, ".")
	if !found {
		return namespace
	}
	return path
}

func fieldMessage(fe validator.FieldError) string {
	isText := fe.Kind() == reflect.String
	isList := fe.Kind() == reflect.Slice
	switch fe.Tag() {
	case "required":
		if isList {
			return "adicione pelo menos um item"
		}
		return "campo obrigatório"
	case "email":
		return "e-mail inválido"
	case "httpsurl":
		return "use um link que comece com https://"
	case "uuid":
		return "id inválido"
	case "min":
		switch {
		case isText:
			return fmt.Sprintf("use pelo menos %s caracteres", fe.Param())
		case isList && fe.Param() == "1":
			return "adicione pelo menos um item"
		case isList:
			return fmt.Sprintf("adicione pelo menos %s itens", fe.Param())
		default:
			return fmt.Sprintf("o mínimo é %s", fe.Param())
		}
	case "max":
		switch {
		case isText:
			return fmt.Sprintf("use no máximo %s caracteres", fe.Param())
		case isList:
			return fmt.Sprintf("no máximo %s itens", fe.Param())
		default:
			return fmt.Sprintf("o máximo é %s", fe.Param())
		}
	}
	return "valor inválido"
}
