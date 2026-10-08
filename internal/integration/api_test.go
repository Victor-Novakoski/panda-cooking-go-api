//go:build integration

package integration_test

import (
	"fmt"
	"net/http"
	"net/url"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const (
	catSalgados = 2
	catPaes     = 12
	password    = "panelas-de-barro"
)

func TestAuthFlow(t *testing.T) {
	reset(t)
	c := newClient(t, newAPI(t))

	created := c.signup("Maria", " Maria@Email.com ", password)
	assert.Equal(t, "maria@email.com", created.Email)
	status, _ := c.call(http.MethodPost, "/api/users", map[string]string{"name": "Outra", "email": "MARIA@email.com", "password": password})
	assert.Equal(t, http.StatusConflict, status, "o e-mail é único sem diferenciar maiúscula")

	c.login("MARIA@email.com", password)
	first := c.cookie("/api/auth/refresh", "panda_refresh")
	require.NotEmpty(t, first)
	assert.Empty(t, c.cookie("/api/recipes", "panda_refresh"), "o refresh token só vai para /api/auth")
	assert.Equal(t, "1", c.cookie("/", "panda_session"))

	var me user
	c.must(http.MethodGet, "/api/users/profile", nil, http.StatusOK, &me)
	assert.Equal(t, created.ID, me.ID)

	// renovar troca o cookie
	require.Equal(t, http.StatusOK, c.refresh())
	second := c.cookie("/api/auth/refresh", "panda_refresh")
	assert.NotEqual(t, first, second)
	c.must(http.MethodGet, "/api/users/profile", nil, http.StatusOK, nil)

	// o token anterior ainda vale por alguns segundos (duas abas renovando juntas)
	c.setCookie("/api/auth", "panda_refresh", first)
	require.Equal(t, http.StatusOK, c.refresh())
	current := c.cookie("/api/auth/refresh", "panda_refresh")

	// sair encerra a sessão: nenhum token dela renova mais
	c.must(http.MethodPost, "/api/auth/logout", nil, http.StatusNoContent, nil)
	assert.Empty(t, c.cookie("/api/auth/refresh", "panda_refresh"))
	for _, token := range []string{current, second, first} {
		c.setCookie("/api/auth", "panda_refresh", token)
		assert.Equal(t, http.StatusUnauthorized, c.refresh())
	}

	status, _ = c.call(http.MethodPost, "/api/auth/login", map[string]string{"email": "maria@email.com", "password": "senha-errada"})
	assert.Equal(t, http.StatusUnauthorized, status)
}

func TestRecipeSearchAndPagination(t *testing.T) {
	reset(t)
	srv := newAPI(t)
	maria := newClient(t, srv)
	mariaUser := maria.signup("Maria", "maria@email.com", password)
	maria.login("maria@email.com", password)
	joao := newClient(t, srv)
	joaoUser := joao.signup("João", "joao@email.com", password)

	for i := 1; i <= 13; i++ {
		maria.must(http.MethodPost, "/api/recipes", recipeBody(fmt.Sprintf("Receita %02d", i), "Uma receita qualquer de teste.", catSalgados), http.StatusCreated, nil)
	}
	maria.must(http.MethodPost, "/api/recipes", recipeBody("Pão de Queijo Mineiro", "Clássico de Minas Gerais, crocante por fora.", catPaes), http.StatusCreated, nil)
	maria.must(http.MethodPost, "/api/recipes", recipeBody("Bolo 100% integral", "Bolo fofinho com farinha_integral.", catPaes), http.StatusCreated, nil)

	list := func(query string) page[recipe] {
		var p page[recipe]
		maria.must(http.MethodGet, "/api/recipes?"+query, nil, http.StatusOK, &p)
		return p
	}

	first := list("")
	assert.Equal(t, int64(15), first.Total)
	assert.Equal(t, 2, first.TotalPages)
	require.Len(t, first.Items, 12)
	assert.Equal(t, "Bolo 100% integral", first.Items[0].Name, "as mais novas primeiro")
	assert.Equal(t, "Pão de Queijo Mineiro", first.Items[1].Name)
	assert.Equal(t, "Maria", first.Items[0].Author.Name)
	assert.Equal(t, "Pães e Bolos", first.Items[0].Category.Name)
	assert.Equal(t, "https://img.exemplo.com/bolo-100-integral.jpg", first.Items[0].ImageURL)

	second := list("page=2")
	assert.Equal(t, []string{"Receita 03", "Receita 02", "Receita 01"}, names(second.Items))
	assert.Empty(t, list("page=3").Items)

	search := func(term string) []string { return names(list("search=" + url.QueryEscape(term)).Items) }
	assert.Equal(t, []string{"Pão de Queijo Mineiro"}, search("pao de queijo"), "sem acento")
	assert.Equal(t, []string{"Pão de Queijo Mineiro"}, search("PÃO DE QUEIJO"), "sem diferenciar maiúscula")
	assert.Equal(t, []string{"Pão de Queijo Mineiro"}, search("minas gerais"), "busca também na descrição")
	assert.Equal(t, []string{"Bolo 100% integral"}, search("%"), "% é texto, não curinga")
	assert.Equal(t, []string{"Bolo 100% integral"}, search("_"), "_ é texto, não curinga")
	assert.Empty(t, search("lasanha"))

	assert.Equal(t, int64(2), list(fmt.Sprintf("category_id=%d", catPaes)).Total)
	assert.Equal(t, int64(15), list("user_id="+mariaUser.ID).Total)
	assert.Equal(t, int64(0), list("user_id="+joaoUser.ID).Total)
	assert.Equal(t, []string{"Bolo 100% integral"}, names(list(fmt.Sprintf("category_id=%d&search=bolo&user_id=%s", catPaes, mariaUser.ID)).Items))
}

func TestRecipeLifecycle(t *testing.T) {
	reset(t)
	srv := newAPI(t)
	maria := newClient(t, srv)
	maria.signup("Maria", "maria@email.com", password)
	maria.login("maria@email.com", password)
	joao := newClient(t, srv)
	joao.signup("João", "joao@email.com", password)
	joao.login("joao@email.com", password)

	var bolo recipe
	maria.must(http.MethodPost, "/api/recipes", recipeBody("Bolo de cenoura", "Bolo fofinho de cenoura.", catPaes, "Farinha de Trigo:2 xícaras", "  OVOS :3 unidades"), http.StatusCreated, &bolo)
	require.Len(t, bolo.Ingredients, 2)
	assert.Equal(t, "farinha de trigo", bolo.Ingredients[0].Name)
	assert.Equal(t, "ovos", bolo.Ingredients[1].Name)
	assert.Equal(t, "3 unidades", bolo.Ingredients[1].Amount)

	// o mesmo ingrediente escrito de outro jeito é o mesmo registro
	maria.must(http.MethodPost, "/api/recipes", recipeBody("Pão caseiro", "Pão de fermentação natural.", catPaes, "farinha  de TRIGO:1 kg"), http.StatusCreated, nil)
	var count int64
	require.NoError(t, db.Raw("SELECT count(*) FROM ingredients WHERE name = 'farinha de trigo'").Scan(&count).Error)
	assert.Equal(t, int64(1), count)

	// PUT troca tudo, na ordem enviada
	body := recipeBody("Bolo de cenoura com chocolate", "Bolo de cenoura com cobertura.", catPaes, "Cenoura:3", "Açúcar:2 xícaras", "Chocolate:1 barra")
	body["images"] = []map[string]string{}
	var replaced recipe
	maria.must(http.MethodPut, "/api/recipes/"+bolo.ID, body, http.StatusOK, &replaced)
	assert.Equal(t, "Bolo de cenoura com chocolate", replaced.Name)
	assert.Empty(t, replaced.Images)
	assert.Equal(t, []string{"cenoura", "açúcar", "chocolate"}, []string{replaced.Ingredients[0].Name, replaced.Ingredients[1].Name, replaced.Ingredients[2].Name})

	// PATCH muda só o que veio
	var patched recipe
	maria.must(http.MethodPatch, "/api/recipes/"+bolo.ID, map[string]any{"portions": 10}, http.StatusOK, &patched)
	assert.Equal(t, 10, patched.Portions)
	assert.Equal(t, "Bolo de cenoura com chocolate", patched.Name)
	assert.Len(t, patched.Ingredients, 3)
	status, _ := maria.call(http.MethodPatch, "/api/recipes/"+bolo.ID, map[string]any{"category_id": 999})
	assert.Equal(t, http.StatusUnprocessableEntity, status)

	// itens um a um
	var img struct {
		ID  uint   `json:"id"`
		URL string `json:"url"`
	}
	maria.must(http.MethodPost, "/api/recipes/"+bolo.ID+"/images", map[string]string{"url": "https://img.exemplo.com/a.jpg"}, http.StatusCreated, &img)
	maria.must(http.MethodPatch, fmt.Sprintf("/api/recipes/%s/images/%d", bolo.ID, img.ID), map[string]string{"url": "https://img.exemplo.com/b.jpg"}, http.StatusOK, &img)
	assert.Equal(t, "https://img.exemplo.com/b.jpg", img.URL)
	var ing struct {
		ID   uint   `json:"id"`
		Name string `json:"name"`
	}
	maria.must(http.MethodPost, "/api/recipes/"+bolo.ID+"/ingredients", map[string]string{"name": "Fermento  em Pó", "amount": "1 colher"}, http.StatusCreated, &ing)
	assert.Equal(t, "fermento em pó", ing.Name)
	var prep struct {
		ID uint `json:"id"`
	}
	maria.must(http.MethodPost, "/api/recipes/"+bolo.ID+"/preparations", map[string]string{"description": "Asse por 40 minutos."}, http.StatusCreated, &prep)
	maria.must(http.MethodPatch, fmt.Sprintf("/api/recipes/%s/preparations/%d", bolo.ID, prep.ID), map[string]string{"description": "Asse por 45 minutos."}, http.StatusOK, nil)

	// outra pessoa não mexe
	status, _ = joao.call(http.MethodPut, "/api/recipes/"+bolo.ID, body)
	assert.Equal(t, http.StatusForbidden, status)
	status, _ = joao.call(http.MethodDelete, fmt.Sprintf("/api/recipes/%s/images/%d", bolo.ID, img.ID), nil)
	assert.Equal(t, http.StatusForbidden, status)

	// item de outra receita pela URL desta é 404 (IDOR)
	var pao recipe
	maria.must(http.MethodPost, "/api/recipes", recipeBody("Pão de milho", "Pão de milho caseiro.", catPaes), http.StatusCreated, &pao)
	status, _ = maria.call(http.MethodDelete, fmt.Sprintf("/api/recipes/%s/preparations/%d", pao.ID, prep.ID), nil)
	assert.Equal(t, http.StatusNotFound, status)

	maria.must(http.MethodDelete, fmt.Sprintf("/api/recipes/%s/images/%d", bolo.ID, img.ID), nil, http.StatusNoContent, nil)
	maria.must(http.MethodDelete, fmt.Sprintf("/api/recipes/%s/ingredients/%d", bolo.ID, ing.ID), nil, http.StatusNoContent, nil)
	maria.must(http.MethodDelete, fmt.Sprintf("/api/recipes/%s/preparations/%d", bolo.ID, prep.ID), nil, http.StatusNoContent, nil)

	var final recipe
	maria.must(http.MethodGet, "/api/recipes/"+bolo.ID, nil, http.StatusOK, &final)
	assert.Empty(t, final.Images)
	assert.Len(t, final.Ingredients, 3)
	assert.Len(t, final.Preparations, 2)

	// apagar leva junto comentários e favoritos
	joao.must(http.MethodPost, "/api/recipes/"+bolo.ID+"/comments", map[string]string{"description": "Que delícia!"}, http.StatusCreated, nil)
	joao.must(http.MethodPost, "/api/favorites/"+bolo.ID, nil, http.StatusCreated, nil)
	maria.must(http.MethodDelete, "/api/recipes/"+bolo.ID, nil, http.StatusNoContent, nil)
	status, _ = maria.call(http.MethodGet, "/api/recipes/"+bolo.ID, nil)
	assert.Equal(t, http.StatusNotFound, status)
	for _, table := range []string{"comments", "favorite_recipes", "image_recipes", "ingredient_recipes", "preparations"} {
		require.NoError(t, db.Raw("SELECT count(*) FROM "+table+" WHERE recipe_id = ?", bolo.ID).Scan(&count).Error)
		assert.Zerof(t, count, "%s da receita apagada", table)
	}
}

func TestCommentsAndFavorites(t *testing.T) {
	reset(t)
	srv := newAPI(t)
	maria := newClient(t, srv)
	maria.signup("Maria", "maria@email.com", password)
	maria.login("maria@email.com", password)
	joao := newClient(t, srv)
	joao.signup("João", "joao@email.com", password)
	joao.login("joao@email.com", password)

	var r1, r2 recipe
	maria.must(http.MethodPost, "/api/recipes", recipeBody("Coxinha", "Coxinha de frango cremosa.", catSalgados), http.StatusCreated, &r1)
	maria.must(http.MethodPost, "/api/recipes", recipeBody("Empada", "Empada de palmito.", catSalgados), http.StatusCreated, &r2)

	var c1, c2 comment
	joao.must(http.MethodPost, "/api/recipes/"+r1.ID+"/comments", map[string]string{"description": "Primeiro!"}, http.StatusCreated, &c1)
	joao.must(http.MethodPost, "/api/recipes/"+r1.ID+"/comments", map[string]string{"description": "Fiz de novo."}, http.StatusCreated, &c2)
	assert.Equal(t, "João", c2.User.Name)

	var comments page[comment]
	maria.must(http.MethodGet, "/api/recipes/"+r1.ID+"/comments?per_page=1", nil, http.StatusOK, &comments)
	assert.Equal(t, int64(2), comments.Total)
	assert.Equal(t, 2, comments.TotalPages)
	assert.Equal(t, "Fiz de novo.", comments.Items[0].Description, "os mais novos primeiro")

	joao.must(http.MethodPatch, fmt.Sprintf("/api/comments/%d", c1.ID), map[string]string{"description": "Primeiro! Ficou ótima."}, http.StatusOK, nil)
	status, _ := maria.call(http.MethodPatch, fmt.Sprintf("/api/comments/%d", c1.ID), map[string]string{"description": "Editado"})
	assert.Equal(t, http.StatusForbidden, status)
	status, _ = maria.call(http.MethodDelete, fmt.Sprintf("/api/comments/%d", c1.ID), nil)
	assert.Equal(t, http.StatusForbidden, status)

	// admin apaga comentário de qualquer um; o papel vai no token, então entra de novo
	require.NoError(t, db.Exec("UPDATE users SET is_adm = true WHERE email = 'maria@email.com'").Error)
	maria.login("maria@email.com", password)
	maria.must(http.MethodDelete, fmt.Sprintf("/api/comments/%d", c1.ID), nil, http.StatusNoContent, nil)

	// favoritos
	joao.must(http.MethodPost, "/api/favorites/"+r1.ID, nil, http.StatusCreated, nil)
	status, _ = joao.call(http.MethodPost, "/api/favorites/"+r1.ID, nil)
	assert.Equal(t, http.StatusConflict, status)
	joao.must(http.MethodPost, "/api/favorites/"+r2.ID, nil, http.StatusCreated, nil)

	var fav struct {
		Favorite bool `json:"favorite"`
	}
	joao.must(http.MethodGet, "/api/favorites/"+r1.ID, nil, http.StatusOK, &fav)
	assert.True(t, fav.Favorite)
	maria.must(http.MethodGet, "/api/favorites/"+r1.ID, nil, http.StatusOK, &fav)
	assert.False(t, fav.Favorite)

	var favorites page[recipe]
	joao.must(http.MethodGet, "/api/users/profile/favorites", nil, http.StatusOK, &favorites)
	assert.Equal(t, []string{"Empada", "Coxinha"}, names(favorites.Items), "o favoritado por último primeiro")

	joao.must(http.MethodDelete, "/api/favorites/"+r1.ID, nil, http.StatusNoContent, nil)
	status, _ = joao.call(http.MethodDelete, "/api/favorites/"+r1.ID, nil)
	assert.Equal(t, http.StatusNotFound, status)
	status, _ = joao.call(http.MethodPost, "/api/favorites/00000000-0000-4000-8000-000000000000", nil)
	assert.Equal(t, http.StatusNotFound, status)
}

func TestDeleteAccount(t *testing.T) {
	reset(t)
	srv := newAPI(t)
	maria := newClient(t, srv)
	maria.signup("Maria", "maria@email.com", password)
	maria.login("maria@email.com", password)
	joao := newClient(t, srv)
	joao.signup("João", "joao@email.com", password)
	joao.login("joao@email.com", password)

	var r recipe
	maria.must(http.MethodPost, "/api/recipes", recipeBody("Pudim", "Pudim de leite condensado.", 1), http.StatusCreated, &r)
	joao.must(http.MethodPost, "/api/recipes/"+r.ID+"/comments", map[string]string{"description": "Sem furinhos!"}, http.StatusCreated, nil)
	joao.must(http.MethodPost, "/api/favorites/"+r.ID, nil, http.StatusCreated, nil)

	maria.must(http.MethodDelete, "/api/users/profile", nil, http.StatusNoContent, nil)

	status, _ := joao.call(http.MethodGet, "/api/recipes/"+r.ID, nil)
	assert.Equal(t, http.StatusNotFound, status, "as receitas vão junto com a conta")
	var favorites page[recipe]
	joao.must(http.MethodGet, "/api/users/profile/favorites", nil, http.StatusOK, &favorites)
	assert.Empty(t, favorites.Items)

	assert.Equal(t, http.StatusUnauthorized, maria.refresh(), "as sessões vão junto")
	status, _ = maria.call(http.MethodPost, "/api/auth/login", map[string]string{"email": "maria@email.com", "password": password})
	assert.Equal(t, http.StatusUnauthorized, status)
	// o access token que ainda não venceu não acha mais a conta
	status, _ = maria.call(http.MethodGet, "/api/users/profile", nil)
	assert.Equal(t, http.StatusNotFound, status)
}
