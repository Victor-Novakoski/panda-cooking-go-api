package config_test

import (
	"strings"
	"testing"

	"panda-cooking-go-api/internal/config"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestLoad_SecretKey(t *testing.T) {
	strong := strings.Repeat("k", config.MinSecretKeyLen)

	tests := []struct {
		name    string
		env     string
		key     string
		wantErr string
	}{
		{name: "sem SECRET_KEY recusa", key: "", wantErr: "SECRET_KEY não definida"},
		{name: "chave curta recusa", key: "changeme", wantErr: "pelo menos 32"},
		{name: "chave de exemplo em produção recusa", env: "production", key: config.ExampleSecretKey, wantErr: "exemplo"},
		{name: "chave de exemplo em desenvolvimento aceita", key: config.ExampleSecretKey},
		{name: "chave forte em produção aceita", env: "production", key: strong},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Setenv("APP_ENV", tt.env)
			t.Setenv("SECRET_KEY", tt.key)

			cfg, err := config.Load()

			if tt.wantErr != "" {
				require.Error(t, err)
				assert.Contains(t, err.Error(), tt.wantErr)
				return
			}
			require.NoError(t, err)
			assert.Equal(t, tt.key, cfg.SecretKey)
		})
	}
}
