package config

import (
	"errors"
	"fmt"
	"os"
)

// MinSecretKeyLen é o tamanho mínimo da SECRET_KEY: 32 bytes = 256 bits,
// o mesmo tamanho da saída do HS256.
const MinSecretKeyLen = 32

// ExampleSecretKey é o valor do .env.example. Produção recusa essa chave,
// porque qualquer um que leia o repositório consegue assinar tokens com ela.
const ExampleSecretKey = "troque-por-uma-chave-aleatoria-de-32-caracteres-ou-mais" //nolint:gosec // valor público do .env.example, não é segredo

type Config struct {
	Env       string
	Port      string
	SecretKey string
	DB        DBConfig
}

type DBConfig struct {
	Host     string
	Port     string
	User     string
	Password string
	Name     string
}

func (d DBConfig) DSN() string {
	return fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		d.Host, d.Port, d.User, d.Password, d.Name,
	)
}

// IsProduction diz se a API está rodando com APP_ENV=production.
func (c Config) IsProduction() bool { return c.Env == "production" }

// Load lê a configuração das variáveis de ambiente e recusa uma SECRET_KEY
// ausente, curta ou, em produção, igual à de exemplo.
func Load() (Config, error) {
	cfg := Config{
		Env:       env("APP_ENV", "development"),
		Port:      env("PORT", "8080"),
		SecretKey: os.Getenv("SECRET_KEY"),
		DB: DBConfig{
			Host:     env("DB_HOST", "localhost"),
			Port:     env("DB_PORT", "5432"),
			User:     env("DB_USER", "postgres"),
			Password: env("DB_PASSWORD", "postgres"),
			Name:     env("DB_NAME", "panda_cooking"),
		},
	}

	switch {
	case cfg.SecretKey == "":
		return Config{}, errors.New("SECRET_KEY não definida")
	case len(cfg.SecretKey) < MinSecretKeyLen:
		return Config{}, fmt.Errorf("SECRET_KEY precisa ter pelo menos %d caracteres", MinSecretKeyLen)
	case cfg.IsProduction() && cfg.SecretKey == ExampleSecretKey:
		return Config{}, errors.New("SECRET_KEY de exemplo não pode ser usada em produção")
	}

	return cfg, nil
}

func env(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}
