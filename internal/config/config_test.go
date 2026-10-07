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
			t.Setenv("CORS_ORIGINS", "https://panda.exemplo.com")

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

func TestLoad_CORSOrigins(t *testing.T) {
	tests := []struct {
		name    string
		env     string
		origins string
		want    []string
		wantErr string
	}{
		{name: "sem valor usa o front local", origins: "", want: []string{config.DefaultCORSOrigin}},
		{name: "lista separada por vírgula", origins: " http://localhost:3000, https://panda.exemplo.com/ ", want: []string{"http://localhost:3000", "https://panda.exemplo.com"}},
		{name: "produção exige https", env: "production", origins: "http://panda.exemplo.com", wantErr: "https"},
		{name: "produção com https aceita", env: "production", origins: "https://panda.exemplo.com", want: []string{"https://panda.exemplo.com"}},
		{name: "caminho não é origem", origins: "https://panda.exemplo.com/app", wantErr: "não é uma origem"},
		{name: "sem esquema recusa", origins: "panda.exemplo.com", wantErr: "não é uma origem"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Setenv("APP_ENV", tt.env)
			t.Setenv("SECRET_KEY", strings.Repeat("k", config.MinSecretKeyLen))
			t.Setenv("CORS_ORIGINS", tt.origins)

			cfg, err := config.Load()

			if tt.wantErr != "" {
				require.Error(t, err)
				assert.Contains(t, err.Error(), tt.wantErr)
				return
			}
			require.NoError(t, err)
			assert.Equal(t, tt.want, cfg.CORSOrigins)
		})
	}
}

func TestLoad_TrustedProxies(t *testing.T) {
	t.Setenv("SECRET_KEY", strings.Repeat("k", config.MinSecretKeyLen))

	t.Setenv("TRUSTED_PROXIES", "172.16.0.0/12, 10.0.0.5")
	cfg, err := config.Load()
	require.NoError(t, err)
	assert.Equal(t, []string{"172.16.0.0/12", "10.0.0.5"}, cfg.TrustedProxies)

	t.Setenv("TRUSTED_PROXIES", "proxy.local")
	_, err = config.Load()
	assert.ErrorContains(t, err, "TRUSTED_PROXIES")
}

func TestDBConfig_URL(t *testing.T) {
	db := config.DBConfig{Host: "postgres", Port: "5432", User: "panda", Password: "p@ss word'", Name: "panda_cooking", SSLMode: "require"}

	assert.Equal(t, "postgres://panda:p%40ss%20word%27@postgres:5432/panda_cooking?sslmode=require", db.URL())
}
