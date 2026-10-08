//go:build integration

// Testes de integração: API, services e repositórios contra um Postgres de
// verdade, com as migrations aplicadas. Rodam com make test-integration.
//
// O banco é um container novo (testcontainers), ou o de TEST_DATABASE_URL
// quando definido. Os testes apagam todos os dados: por isso o nome do banco
// em TEST_DATABASE_URL precisa ter "test".
package integration_test

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net/url"
	"os"
	"strings"
	"testing"

	"panda-cooking-go-api/internal/database"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/require"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/postgres"
	"gorm.io/gorm"
)

var (
	db    *gorm.DB
	dbURL string
	ctx   = context.Background()
)

func TestMain(m *testing.M) {
	gin.SetMode(gin.TestMode)
	os.Exit(run(m))
}

func run(m *testing.M) int {
	var stop func()
	var err error
	dbURL, stop, err = startPostgres()
	if err != nil {
		fmt.Fprintln(os.Stderr, "subir o Postgres:", err)
		return 1
	}
	defer stop()

	db, err = database.Connect(ctx, dbURL, slog.New(slog.DiscardHandler))
	if err != nil {
		fmt.Fprintln(os.Stderr, "conectar:", err)
		return 1
	}
	if err := database.Migrate(ctx, db, dbURL); err != nil {
		fmt.Fprintln(os.Stderr, "migrations:", err)
		return 1
	}
	return m.Run()
}

func startPostgres() (string, func(), error) {
	if raw := os.Getenv("TEST_DATABASE_URL"); raw != "" {
		u, err := url.Parse(raw)
		if err != nil {
			return "", nil, err
		}
		if !strings.Contains(u.Path, "test") {
			return "", nil, errors.New(`TEST_DATABASE_URL precisa apontar para um banco com "test" no nome: os testes apagam tudo`)
		}
		return raw, func() {}, nil
	}

	ctr, err := postgres.Run(ctx, "postgres:16-alpine",
		postgres.WithDatabase("panda_cooking_test"),
		postgres.WithUsername("panda"),
		postgres.WithPassword("panda"),
		postgres.BasicWaitStrategies(),
	)
	if err != nil {
		return "", nil, err
	}
	stop := func() { _ = testcontainers.TerminateContainer(ctr) }
	connURL, err := ctr.ConnectionString(ctx, "sslmode=disable")
	if err != nil {
		stop()
		return "", nil, err
	}
	return connURL, stop, nil
}

// reset apaga os dados entre os testes. As categorias ficam: vêm da migration.
func reset(t *testing.T) {
	t.Helper()
	require.NoError(t, db.Exec("TRUNCATE users, ingredients RESTART IDENTITY CASCADE").Error)
}
