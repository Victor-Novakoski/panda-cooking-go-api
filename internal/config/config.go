package config

import (
	"errors"
	"fmt"
	"net"
	"net/url"
	"os"
	"strings"
)

// MinSecretKeyLen é o tamanho mínimo da SECRET_KEY: 32 bytes = 256 bits,
// o mesmo tamanho da saída do HS256.
const MinSecretKeyLen = 32

// ExampleSecretKey é o valor do .env.example. Produção recusa essa chave,
// porque qualquer um que leia o repositório consegue assinar tokens com ela.
const ExampleSecretKey = "troque-por-uma-chave-aleatoria-de-32-caracteres-ou-mais" //nolint:gosec // valor público do .env.example, não é segredo

// DefaultCORSOrigin é o front rodando local (next dev e docker compose).
const DefaultCORSOrigin = "http://localhost:3000"

type Config struct {
	Env       string
	Port      string
	SecretKey string
	DB        DBConfig
	// SeedDemo cria os usuários e receitas de demonstração ao subir, se o
	// banco não tiver usuário nenhum (SEED_DEMO=true).
	SeedDemo bool
	// CORSOrigins são as origens do front que podem chamar a API pelo
	// navegador e usar o cookie de sessão nas rotas de /api/auth.
	CORSOrigins []string
	// TrustedProxies são os IPs ou faixas (CIDR) dos proxies na frente da
	// API. Só deles o X-Forwarded-For é aceito para descobrir o IP do cliente.
	TrustedProxies []string
}

type DBConfig struct {
	Host     string
	Port     string
	User     string
	Password string
	Name     string
	SSLMode  string
}

// URL monta a conexão no formato postgres://, com usuário e senha escapados
// (uma senha com espaço ou aspas quebraria o formato chave=valor).
func (d DBConfig) URL() string {
	u := url.URL{
		Scheme:   "postgres",
		User:     url.UserPassword(d.User, d.Password),
		Host:     net.JoinHostPort(d.Host, d.Port),
		Path:     "/" + d.Name,
		RawQuery: url.Values{"sslmode": {d.SSLMode}}.Encode(),
	}
	return u.String()
}

// IsProduction diz se a API está rodando com APP_ENV=production.
func (c Config) IsProduction() bool { return c.Env == "production" }

// Load lê a configuração das variáveis de ambiente e recusa o que deixaria a
// API insegura: SECRET_KEY ausente, curta ou, em produção, igual à de
// exemplo; origem de CORS sem https em produção; proxy confiável inválido.
func Load() (Config, error) {
	cfg := Config{
		Env:            env("APP_ENV", "development"),
		Port:           env("PORT", "8080"),
		SecretKey:      os.Getenv("SECRET_KEY"),
		DB:             LoadDB(),
		SeedDemo:       os.Getenv("SEED_DEMO") == "true",
		CORSOrigins:    list(env("CORS_ORIGINS", DefaultCORSOrigin)),
		TrustedProxies: list(os.Getenv("TRUSTED_PROXIES")),
	}

	switch {
	case cfg.SecretKey == "":
		return Config{}, errors.New("SECRET_KEY não definida")
	case len(cfg.SecretKey) < MinSecretKeyLen:
		return Config{}, fmt.Errorf("SECRET_KEY precisa ter pelo menos %d caracteres", MinSecretKeyLen)
	case cfg.IsProduction() && cfg.SecretKey == ExampleSecretKey:
		return Config{}, errors.New("SECRET_KEY de exemplo não pode ser usada em produção")
	}

	for _, origin := range cfg.CORSOrigins {
		if err := checkOrigin(origin, cfg.IsProduction()); err != nil {
			return Config{}, err
		}
	}
	for _, proxy := range cfg.TrustedProxies {
		if net.ParseIP(proxy) == nil {
			if _, _, err := net.ParseCIDR(proxy); err != nil {
				return Config{}, fmt.Errorf("TRUSTED_PROXIES: %q não é IP nem faixa CIDR", proxy)
			}
		}
	}

	return cfg, nil
}

// LoadDB lê só a conexão com o banco; o seed usa sem precisar da SECRET_KEY.
func LoadDB() DBConfig {
	return DBConfig{
		Host:     env("DB_HOST", "localhost"),
		Port:     env("DB_PORT", "5432"),
		User:     env("DB_USER", "postgres"),
		Password: env("DB_PASSWORD", "postgres"),
		Name:     env("DB_NAME", "panda_cooking"),
		SSLMode:  env("DB_SSLMODE", "disable"),
	}
}

// checkOrigin aceita só origem completa (esquema e host, sem caminho). Em
// produção o front precisa estar em https, senão o cookie Secure não volta.
func checkOrigin(origin string, production bool) error {
	u, err := url.Parse(origin)
	if err != nil || u.Host == "" || (u.Path != "" && u.Path != "/") || u.RawQuery != "" {
		return fmt.Errorf("CORS_ORIGINS: %q não é uma origem (ex.: https://panda.exemplo.com)", origin)
	}
	if u.Scheme != "https" && (production || u.Scheme != "http") {
		return fmt.Errorf("CORS_ORIGINS: %q precisa usar https em produção", origin)
	}
	return nil
}

// list separa uma variável por vírgula, ignorando espaços e itens vazios.
func list(value string) []string {
	var items []string
	for _, item := range strings.Split(value, ",") {
		if item = strings.TrimSpace(item); item != "" {
			items = append(items, strings.TrimSuffix(item, "/"))
		}
	}
	return items
}

func env(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}
