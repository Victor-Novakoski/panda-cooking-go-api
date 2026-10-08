package service

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCheckPassword(t *testing.T) {
	tests := []struct {
		name     string
		password string
		email    string
		wantErr  string
	}{
		{"senha boa", "panelas-de-barro-42", "maria@email.com", ""},
		{"senha comum", "1234567890", "maria@email.com", "senha muito comum"},
		{"senha comum com maiúscula", "QwertyUiop", "maria@email.com", "senha muito comum"},
		{"senha da demonstração", "panda-cooking-demo", "maria@email.com", "senha muito comum"},
		{"caractere repetido", "zzzzzzzzzzzz", "maria@email.com", "senha muito comum"},
		{"igual ao e-mail", "Maria@Email.com", "maria@email.com", "não pode ser o seu e-mail"},
		{"igual ao começo do e-mail", "mariasilva1", "mariasilva1@email.com", "não pode ser o seu e-mail"},
		{"passa de 72 bytes", "ção-ção-ção-ção-ção-ção-ção-ção-ção-ção-ção-ção-ção-ção-ção", "maria@email.com", "72 bytes"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := checkPassword(tt.password, tt.email)
			if tt.wantErr == "" {
				assert.NoError(t, err)
				return
			}
			var appErr *Error
			if assert.ErrorAs(t, err, &appErr) {
				assert.Equal(t, "password", appErr.Field)
				assert.Contains(t, appErr.Message, tt.wantErr)
			}
		})
	}
}
