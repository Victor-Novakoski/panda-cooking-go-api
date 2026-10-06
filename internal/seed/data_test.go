package seed

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// Os dados de demonstração se referem a autor e categoria por posição e nome;
// um erro de digitação aqui só apareceria ao rodar o seed contra o banco.
func TestDemoRecipes_ReferenciasValidas(t *testing.T) {
	categories := make(map[string]bool, len(categoryNames))
	for _, name := range categoryNames {
		categories[name] = true
	}

	for _, r := range demoRecipes {
		assert.Truef(t, categories[r.category], "%q: categoria %q não existe", r.recipe.Name, r.category)
		assert.Truef(t, r.author >= 0 && r.author < len(demoUsers), "%q: autor %d não existe", r.recipe.Name, r.author)
		assert.NotEmptyf(t, r.ingredients, "%q sem ingredientes", r.recipe.Name)
		assert.NotEmptyf(t, r.preparations, "%q sem modo de preparo", r.recipe.Name)
	}
}

// A senha de demonstração precisa passar na validação do cadastro e do login.
func TestDemoPassword_TamanhoAceito(t *testing.T) {
	assert.GreaterOrEqual(t, len(DemoPassword), 10)
	assert.LessOrEqual(t, len(DemoPassword), 72)
}
