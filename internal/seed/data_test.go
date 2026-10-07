package seed

import (
	"os"
	"testing"

	"panda-cooking-go-api/internal/handler"
	"panda-cooking-go-api/internal/service"

	"github.com/gin-gonic/gin/binding"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// Os dados de demonstração passam pela mesma validação de uma receita
// enviada pela API: o que o seed cria, um usuário também conseguiria criar
// e editar.
func TestDemoRecipes_PassamNaValidacaoDaAPI(t *testing.T) {
	handler.SetupValidator()
	migration, err := os.ReadFile("../database/migrations/000001_init.up.sql")
	require.NoError(t, err)

	for _, r := range demoRecipes {
		t.Run(r.recipe.Name, func(t *testing.T) {
			input := service.RecipeInput{
				Name:        r.recipe.Name,
				Description: r.recipe.Description,
				Time:        r.recipe.Time,
				Portions:    r.recipe.Portions,
				CategoryID:  1,
			}
			for _, url := range r.images {
				input.Images = append(input.Images, service.ImageRecipeInput{URL: url})
			}
			for _, ing := range r.ingredients {
				input.Ingredients = append(input.Ingredients, service.IngredientRecipeInput{Name: ing.name, Amount: ing.amount})
			}
			for _, step := range r.preparations {
				input.Preparations = append(input.Preparations, service.PreparationInput{Description: step})
			}

			assert.NoError(t, binding.Validator.ValidateStruct(&input))
			assert.Containsf(t, string(migration), "('"+r.category+"')", "categoria %q não está na migration", r.category)
			assert.Truef(t, r.author >= 0 && r.author < len(demoUsers), "autor %d não existe", r.author)
			for _, c := range r.comments {
				assert.True(t, c.author >= 0 && c.author < len(demoUsers))
				assert.NoError(t, binding.Validator.ValidateStruct(&service.CommentInput{Description: c.text}))
			}
		})
	}
}

func TestDemoFavorites_ReceitasExistem(t *testing.T) {
	names := map[string]bool{}
	for _, r := range demoRecipes {
		names[r.recipe.Name] = true
	}
	for _, f := range demoFavorites {
		assert.Truef(t, names[f.recipe], "favorito de receita que não existe: %q", f.recipe)
		assert.True(t, f.user >= 0 && f.user < len(demoUsers))
	}
}

// A senha de demonstração precisa caber no bcrypt; ela é recusada no
// cadastro de propósito (está no README), mas os usuários do seed entram.
func TestDemoPassword_TamanhoAceito(t *testing.T) {
	assert.GreaterOrEqual(t, len(DemoPassword), 10)
	assert.LessOrEqual(t, len(DemoPassword), service.MaxPasswordBytes)
}
