//go:build integration

package integration_test

import (
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"net/url"
	"os"
	"sync"
	"testing"
	"time"

	"panda-cooking-go-api/internal/database"
	"panda-cooking-go-api/internal/model"
	"panda-cooking-go-api/internal/repository"
	"panda-cooking-go-api/internal/seed"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/pgx/v5"
	"github.com/golang-migrate/migrate/v4/source/iofs"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func createUser(t *testing.T, email string) model.User {
	t.Helper()
	u := model.User{Name: "Teste", Email: email, PasswordHash: "hash"}
	require.NoError(t, db.Create(&u).Error)
	return u
}

func TestSessionRotate(t *testing.T) {
	reset(t)
	repo := repository.NewSessionRepository(db)
	user := createUser(t, "maria@email.com")
	now := time.Now().UTC().Truncate(time.Microsecond)
	grace := 20 * time.Second

	newSession := func(t *testing.T, hash string, sessionTTL, tokenTTL time.Duration) model.Session {
		t.Helper()
		s := model.Session{UserID: user.ID, ExpiresAt: now.Add(sessionTTL)}
		require.NoError(t, repo.Create(ctx, &s, &model.RefreshToken{TokenHash: hash, ExpiresAt: now.Add(tokenTTL)}))
		return s
	}
	rotate := func(t *testing.T, hash, next string, at time.Time) repository.RotateResult {
		t.Helper()
		res, err := repo.Rotate(ctx, hash, &model.RefreshToken{TokenHash: next, ExpiresAt: at.Add(7 * 24 * time.Hour)}, at, grace)
		require.NoError(t, err)
		return res
	}

	t.Run("token usado de novo depois da tolerância derruba a sessão", func(t *testing.T) {
		session := newSession(t, "t1", 30*24*time.Hour, 7*24*time.Hour)

		res := rotate(t, "t1", "t2", now)
		assert.Equal(t, repository.Rotated, res.Outcome)
		assert.Equal(t, session.ID, res.Session.ID)

		// dentro da tolerância: outra aba renovando junto
		assert.Equal(t, repository.Rotated, rotate(t, "t1", "t3", now.Add(10*time.Second)).Outcome)

		// fora dela: alguém tem uma cópia do token
		assert.Equal(t, repository.Reused, rotate(t, "t1", "t4", now.Add(time.Minute)).Outcome)

		// e a sessão inteira acabou, inclusive para o token mais novo
		assert.Equal(t, repository.SessionEnded, rotate(t, "t2", "t5", now.Add(2*time.Minute)).Outcome)
		var revoked model.Session
		require.NoError(t, db.First(&revoked, "id = ?", session.ID).Error)
		assert.NotNil(t, revoked.RevokedAt)
	})

	t.Run("token que não existe", func(t *testing.T) {
		assert.Equal(t, repository.TokenNotFound, rotate(t, "inventado", "x1", now).Outcome)
	})

	t.Run("token vencido", func(t *testing.T) {
		newSession(t, "vencido", 30*24*time.Hour, time.Hour)
		assert.Equal(t, repository.SessionEnded, rotate(t, "vencido", "x2", now.Add(2*time.Hour)).Outcome)
	})

	t.Run("o novo token não passa do limite da sessão", func(t *testing.T) {
		session := newSession(t, "curta", 24*time.Hour, 7*24*time.Hour)
		require.Equal(t, repository.Rotated, rotate(t, "curta", "curta-2", now).Outcome)

		var next model.RefreshToken
		require.NoError(t, db.First(&next, "token_hash = ?", "curta-2").Error)
		assert.WithinDuration(t, session.ExpiresAt, next.ExpiresAt, time.Millisecond)
		assert.Equal(t, repository.SessionEnded, rotate(t, "curta-2", "curta-3", now.Add(25*time.Hour)).Outcome)
	})

	t.Run("renovações simultâneas do mesmo token esperam a vez", func(t *testing.T) {
		newSession(t, "junto", 30*24*time.Hour, 7*24*time.Hour)

		var wg sync.WaitGroup
		outcomes := make([]repository.RotateOutcome, 5)
		errs := make([]error, 5)
		for i := range 5 {
			wg.Go(func() {
				res, err := repo.Rotate(ctx, "junto", &model.RefreshToken{TokenHash: fmt.Sprintf("junto-%d", i), ExpiresAt: now.Add(time.Hour)}, now, grace)
				outcomes[i], errs[i] = res.Outcome, err
			})
		}
		wg.Wait()

		for i := range 5 {
			require.NoError(t, errs[i])
			assert.Equal(t, repository.Rotated, outcomes[i])
		}
	})

	t.Run("logout e limpeza", func(t *testing.T) {
		active := newSession(t, "ativa", 30*24*time.Hour, 7*24*time.Hour)
		ended := newSession(t, "encerrada", 30*24*time.Hour, 7*24*time.Hour)
		require.NoError(t, repo.RevokeByToken(ctx, "encerrada", now))
		expired := newSession(t, "velha", time.Hour, time.Hour)

		require.NoError(t, repo.DeleteExpired(ctx, user.ID, now.Add(2*time.Hour)))

		var ids []string
		require.NoError(t, db.Model(&model.Session{}).Where("user_id = ?", user.ID).Pluck("id", &ids).Error)
		assert.Contains(t, ids, active.ID)
		assert.NotContains(t, ids, ended.ID)
		assert.NotContains(t, ids, expired.ID)
	})
}

func TestSeed(t *testing.T) {
	reset(t)

	created, err := seed.Demo(ctx, db, false)
	require.NoError(t, err)
	assert.True(t, created)

	counts := func() map[string]int64 {
		out := map[string]int64{}
		for _, table := range []string{"users", "recipes", "comments", "favorite_recipes", "image_recipes", "preparations"} {
			var n int64
			require.NoError(t, db.Raw("SELECT count(*) FROM "+table).Scan(&n).Error)
			out[table] = n
		}
		return out
	}
	first := counts()
	assert.Equal(t, int64(3), first["users"])
	assert.Equal(t, int64(10), first["recipes"])
	assert.Equal(t, int64(14), first["comments"])
	assert.Equal(t, int64(7), first["favorite_recipes"])

	var notNormalized int64
	require.NoError(t, db.Raw("SELECT count(*) FROM ingredients WHERE name <> lower(name)").Scan(&notNormalized).Error)
	assert.Zero(t, notNormalized, "ingredientes do seed seguem a regra da API")

	// com usuário no banco, não mexe em nada
	created, err = seed.Demo(ctx, db, false)
	require.NoError(t, err)
	assert.False(t, created)

	// reset recria do zero
	created, err = seed.Demo(ctx, db, true)
	require.NoError(t, err)
	assert.True(t, created)
	assert.Equal(t, first, counts())

	c := newClient(t, newAPI(t))
	c.login("maria@pandacooking.com", seed.DemoPassword)
	var recipes page[recipe]
	c.must(http.MethodGet, "/api/recipes", nil, http.StatusOK, &recipes)
	assert.Equal(t, "Brigadeiro Tradicional", recipes.Items[0].Name)
}

func TestMigrations(t *testing.T) {
	t.Run("descem e sobem de novo", func(t *testing.T) {
		src, err := iofs.New(os.DirFS("../database/migrations"), ".")
		require.NoError(t, err)
		m, err := migrate.NewWithSourceInstance("iofs", src, "pgx5://"+dbURL[len("postgres://"):])
		require.NoError(t, err)
		defer func() { _, _ = m.Close() }()

		require.NoError(t, m.Down())
		require.NoError(t, m.Up())

		var categories int64
		require.NoError(t, db.Raw("SELECT count(*) FROM categories").Scan(&categories).Error)
		assert.Equal(t, int64(12), categories)
		// subir de novo sem nada pendente não é erro
		assert.NoError(t, database.Migrate(ctx, db, dbURL))
	})

	t.Run("recusa banco criado pela versão antiga", func(t *testing.T) {
		if err := db.Exec("CREATE DATABASE panda_legacy_test").Error; err != nil {
			t.Skip("sem permissão para criar banco:", err)
		}
		t.Cleanup(func() { db.Exec("DROP DATABASE IF EXISTS panda_legacy_test WITH (FORCE)") })

		u, err := url.Parse(dbURL)
		require.NoError(t, err)
		u.Path = "/panda_legacy_test"
		legacyURL := u.String()
		legacy, err := database.Connect(ctx, legacyURL, slog.New(slog.DiscardHandler))
		require.NoError(t, err)
		sqlDB, err := legacy.DB()
		require.NoError(t, err)
		defer func() { _ = sqlDB.Close() }()
		require.NoError(t, legacy.Exec("CREATE TABLE users (id TEXT PRIMARY KEY)").Error)

		err = database.Migrate(ctx, legacy, legacyURL)

		assert.True(t, errors.Is(err, database.ErrLegacySchema), "erro: %v", err)
	})
}
