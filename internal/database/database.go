// Package database abre a conexão com o Postgres e aplica as migrations.
package database

import (
	"context"
	"database/sql"
	"embed"
	"errors"
	"fmt"
	"log/slog"
	"strings"
	"time"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/pgx/v5" // driver pgx5:// das migrations
	"github.com/golang-migrate/migrate/v4/source/iofs"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

//go:embed migrations/*.sql
var migrations embed.FS

// ErrLegacySchema indica um banco criado pela versão antiga da API, que usava
// o AutoMigrate do GORM e não tem o controle de versão das migrations.
var ErrLegacySchema = errors.New("o banco foi criado pela versão antiga da API (AutoMigrate) e não tem migrations; " +
	"em desenvolvimento, apague o volume e suba de novo: docker compose down -v && docker compose up --build")

// Connect abre o pool de conexões e espera o banco responder, tentando de
// novo por até 30 segundos (o Postgres pode subir depois da API).
func Connect(ctx context.Context, url string, log *slog.Logger) (*gorm.DB, error) {
	db, err := gorm.Open(postgres.Open(url), &gorm.Config{
		// converte erros do Postgres (ex.: chave duplicada) nos erros do GORM
		TranslateError: true,
		Logger: quietLogger{logger.NewSlogLogger(log, logger.Config{
			SlowThreshold:             200 * time.Millisecond,
			LogLevel:                  logger.Warn,
			IgnoreRecordNotFoundError: true,
			// o log mostra a consulta com $1, $2..., nunca os valores (e-mail, hash)
			ParameterizedQueries: true,
		})},
	})
	if err != nil {
		return nil, fmt.Errorf("abrir conexão: %w", err)
	}

	sqlDB, err := db.DB()
	if err != nil {
		return nil, err
	}
	sqlDB.SetMaxOpenConns(20)
	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetConnMaxIdleTime(5 * time.Minute)
	sqlDB.SetConnMaxLifetime(30 * time.Minute)

	if err := waitForDB(ctx, sqlDB, 30*time.Second); err != nil {
		return nil, err
	}
	return db, nil
}

func waitForDB(ctx context.Context, db *sql.DB, limit time.Duration) error {
	ctx, cancel := context.WithTimeout(ctx, limit)
	defer cancel()
	for {
		err := db.PingContext(ctx)
		if err == nil {
			return nil
		}
		select {
		case <-ctx.Done():
			return fmt.Errorf("banco não respondeu: %w", err)
		case <-time.After(time.Second):
		}
	}
}

// Migrate aplica as migrations que faltam. Recusa um banco da versão antiga
// em vez de tentar criar tabelas por cima das que já existem.
func Migrate(ctx context.Context, db *gorm.DB, url string) error {
	if err := checkLegacy(ctx, db); err != nil {
		return err
	}

	src, err := iofs.New(migrations, "migrations")
	if err != nil {
		return err
	}
	m, err := migrate.NewWithSourceInstance("iofs", src, "pgx5://"+strings.TrimPrefix(url, "postgres://"))
	if err != nil {
		return fmt.Errorf("preparar migrations: %w", err)
	}
	defer func() { _, _ = m.Close() }()

	if err := m.Up(); err != nil && !errors.Is(err, migrate.ErrNoChange) {
		return fmt.Errorf("aplicar migrations: %w", err)
	}
	return nil
}

// checkLegacy detecta tabela de usuários sem a tabela de controle das
// migrations, que é o que o AutoMigrate deixava.
func checkLegacy(ctx context.Context, db *gorm.DB) error {
	var legacy bool
	err := db.WithContext(ctx).Raw(
		"SELECT to_regclass('public.users') IS NOT NULL AND to_regclass('public.schema_migrations') IS NULL",
	).Scan(&legacy).Error
	if err != nil {
		return fmt.Errorf("conferir versão do banco: %w", err)
	}
	if legacy {
		return ErrLegacySchema
	}
	return nil
}

// Ping confere se o banco responde; usado pelo /health.
func Ping(ctx context.Context, db *gorm.DB) error {
	sqlDB, err := db.DB()
	if err != nil {
		return err
	}
	return sqlDB.PingContext(ctx)
}

// quietLogger tira do log os erros que a API já trata como resposta
// esperada: e-mail ou favorito duplicado (409) e requisição cancelada pelo
// cliente. O resto (erro de verdade e consulta lenta) continua no log.
type quietLogger struct{ logger.Interface }

func (l quietLogger) LogMode(level logger.LogLevel) logger.Interface {
	return quietLogger{l.Interface.LogMode(level)}
}

func (l quietLogger) Trace(ctx context.Context, begin time.Time, fc func() (string, int64), err error) {
	if errors.Is(err, gorm.ErrDuplicatedKey) || errors.Is(err, context.Canceled) {
		err = nil
	}
	l.Interface.Trace(ctx, begin, fc, err)
}

// ParamsFilter repassa o filtro do logger de dentro; sem ele o GORM voltaria
// a colocar os valores das consultas no log.
func (l quietLogger) ParamsFilter(ctx context.Context, sql string, params ...any) (string, []any) {
	if f, ok := l.Interface.(gorm.ParamsFilter); ok {
		return f.ParamsFilter(ctx, sql, params...)
	}
	return sql, params
}
