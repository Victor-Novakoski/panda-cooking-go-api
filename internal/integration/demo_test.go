//go:build integration

package integration_test

import (
	"fmt"
	"net/http"
	"os"
	"testing"

	"panda-cooking-go-api/internal/database"
	"panda-cooking-go-api/internal/seed"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/pgx/v5"
	"github.com/golang-migrate/migrate/v4/source/iofs"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// Quem entra com a senha pública de demonstração usa o site, mas não apaga
// nem muda o que o seed criou.
func TestDemoDataIsProtected(t *testing.T) {
	reset(t)
	created, err := seed.Demo(ctx, db, true)
	require.NoError(t, err)
	require.True(t, created)

	srv := newAPI(t)
	maria := newClient(t, srv)
	maria.login("maria@pandacooking.com", seed.DemoPassword)

	var me struct {
		ID     string `json:"id"`
		IsAdm  bool   `json:"is_adm"`
		IsDemo bool   `json:"is_demo"`
	}
	maria.must(http.MethodGet, "/api/users/profile", nil, http.StatusOK, &me)
	require.True(t, me.IsAdm)
	require.True(t, me.IsDemo)

	var mine page[recipe]
	maria.must(http.MethodGet, "/api/recipes?user_id="+me.ID, nil, http.StatusOK, &mine)
	require.NotEmpty(t, mine.Items)
	var r struct {
		recipe
		IsDemo bool `json:"is_demo"`
	}
	maria.must(http.MethodGet, "/api/recipes/"+mine.Items[0].ID, nil, http.StatusOK, &r)
	require.True(t, r.IsDemo)
	require.NotEmpty(t, r.Images)
	require.NotEmpty(t, r.Ingredients)
	require.NotEmpty(t, r.Preparations)

	base := "/api/recipes/" + r.ID
	for _, tc := range []struct {
		method, path string
		body         any
	}{
		{http.MethodPut, base, recipeBody("Outro nome", "Outra descrição qualquer.", catPaes)},
		{http.MethodPatch, base, map[string]string{"name": "Outro nome"}},
		{http.MethodPost, base + "/images", map[string]string{"url": "https://images.unsplash.com/photo-1"}},
		{http.MethodDelete, fmt.Sprintf("%s/images/%d", base, r.Images[0].ID), nil},
		{http.MethodDelete, fmt.Sprintf("%s/ingredients/%d", base, r.Ingredients[0].ID), nil},
		{http.MethodPatch, fmt.Sprintf("%s/preparations/%d", base, r.Preparations[0].ID), map[string]string{"description": "Outro passo."}},
		{http.MethodDelete, fmt.Sprintf("%s/preparations/%d", base, r.Preparations[0].ID), nil},
		{http.MethodDelete, base, nil},
		{http.MethodPatch, "/api/users/profile", map[string]string{"name": "Outro nome"}},
		{http.MethodDelete, "/api/users/profile", nil},
	} {
		status, body := maria.call(tc.method, tc.path, tc.body)
		assert.Equal(t, http.StatusForbidden, status, "%s %s: %s", tc.method, tc.path, body)
	}

	var after struct {
		recipe
		IsDemo bool `json:"is_demo"`
	}
	maria.must(http.MethodGet, base, nil, http.StatusOK, &after)
	assert.Equal(t, r.Name, after.Name)
	assert.Len(t, after.Images, len(r.Images))
	assert.Len(t, after.Ingredients, len(r.Ingredients))
	assert.Len(t, after.Preparations, len(r.Preparations))

	// o resto continua liberado: comentar, favoritar, publicar e apagar o
	// que a própria visita publicou, e a admin apagar comentário
	var c comment
	maria.must(http.MethodPost, base+"/comments", map[string]string{"description": "Testando o site."}, http.StatusCreated, &c)
	maria.must(http.MethodDelete, fmt.Sprintf("/api/comments/%d", c.ID), nil, http.StatusNoContent, nil)
	joao := newClient(t, srv)
	joao.login("joao@pandacooking.com", seed.DemoPassword)
	joao.must(http.MethodPost, base+"/comments", map[string]string{"description": "Comentário de visita."}, http.StatusCreated, &c)
	maria.must(http.MethodDelete, fmt.Sprintf("/api/comments/%d", c.ID), nil, http.StatusNoContent, nil)

	var mineNew struct {
		ID     string `json:"id"`
		IsDemo bool   `json:"is_demo"`
	}
	maria.must(http.MethodPost, "/api/recipes", recipeBody("Receita da visita", "Publicada por quem entrou na demonstração.", catPaes), http.StatusCreated, &mineNew)
	assert.False(t, mineNew.IsDemo)
	joao.must(http.MethodPost, "/api/favorites/"+mineNew.ID, nil, http.StatusCreated, nil)
	maria.must(http.MethodPatch, "/api/recipes/"+mineNew.ID, map[string]string{"name": "Receita da visita editada"}, http.StatusOK, nil)
	maria.must(http.MethodDelete, "/api/recipes/"+mineNew.ID, nil, http.StatusNoContent, nil)
}

// A migration 000006 marca o seed que já estava no banco antes dela, e só ele.
func TestDemoMigrationMarksExistingSeed(t *testing.T) {
	reset(t)
	_, err := seed.Demo(ctx, db, true)
	require.NoError(t, err)
	visitor := newClient(t, newAPI(t))
	visitor.signup("Visita", "visita@email.com", password)

	src, err := iofs.New(os.DirFS("../database/migrations"), ".")
	require.NoError(t, err)
	m, err := migrate.NewWithSourceInstance("iofs", src, "pgx5://"+dbURL[len("postgres://"):])
	require.NoError(t, err)
	defer func() { _, _ = m.Close() }()
	t.Cleanup(func() { require.NoError(t, database.Migrate(ctx, db, dbURL)) })

	// volta para antes da 000006, como o banco de produção, e sobe de novo
	require.NoError(t, m.Migrate(5))
	require.NoError(t, m.Up())

	var demoUsers, demoRecipes, allRecipes int64
	require.NoError(t, db.Table("users").Where("is_demo").Count(&demoUsers).Error)
	require.NoError(t, db.Table("recipes").Where("is_demo").Count(&demoRecipes).Error)
	require.NoError(t, db.Table("recipes").Count(&allRecipes).Error)
	assert.Equal(t, int64(3), demoUsers, "só os três usuários do seed")
	assert.Positive(t, demoRecipes)
	assert.Equal(t, allRecipes, demoRecipes)
}
