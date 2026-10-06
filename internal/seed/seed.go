// Package seed cria os dados de demonstração: categorias, ingredientes,
// três usuários que conseguem logar e dez receitas brasileiras.
package seed

import (
	"fmt"

	"panda-cooking-go-api/internal/model"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

// DemoPassword é a senha dos usuários de demonstração. É pública de
// propósito (está no README) para quem visita o projeto conseguir entrar.
const DemoPassword = "panda-cooking-demo" //nolint:gosec // senha pública dos usuários de demonstração

const truncateAll = "TRUNCATE TABLE favorite_recipes, comments, preparations, ingredient_recipes, " +
	"image_recipes, recipes, ingredients, categories, users RESTART IDENTITY CASCADE"

// Demo popula o banco com os dados de demonstração. Com reset=false só age se
// ainda não houver usuário nenhum, para não apagar o que já foi cadastrado;
// com reset=true apaga tudo e recria. Diz se populou.
func Demo(db *gorm.DB, reset bool) (bool, error) {
	if !reset {
		var users int64
		if err := db.Model(&model.User{}).Count(&users).Error; err != nil {
			return false, fmt.Errorf("contar usuários: %w", err)
		}
		if users > 0 {
			return false, nil
		}
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(DemoPassword), bcrypt.DefaultCost)
	if err != nil {
		return false, fmt.Errorf("gerar hash da senha: %w", err)
	}

	err = db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Exec(truncateAll).Error; err != nil {
			return fmt.Errorf("limpar tabelas: %w", err)
		}

		categories := make(map[string]uint, len(categoryNames))
		for _, name := range categoryNames {
			c := model.Category{Name: name}
			if err := tx.Create(&c).Error; err != nil {
				return fmt.Errorf("criar categoria %q: %w", name, err)
			}
			categories[name] = c.ID
		}

		ingredients := make(map[string]uint, len(ingredientNames))
		for _, name := range ingredientNames {
			if _, err := ingredientID(tx, ingredients, name); err != nil {
				return err
			}
		}

		users := make([]string, len(demoUsers))
		for i, u := range demoUsers {
			user := model.User{
				Name:         u.name,
				Email:        u.email,
				Password:     string(hash),
				ImageProfile: u.image,
				IsAdm:        u.isAdm,
			}
			if err := tx.Create(&user).Error; err != nil {
				return fmt.Errorf("criar usuário %q: %w", u.email, err)
			}
			users[i] = user.ID
		}

		for _, r := range demoRecipes {
			if err := createRecipe(tx, r, users[r.author], categories[r.category], ingredients); err != nil {
				return err
			}
		}
		return nil
	})
	if err != nil {
		return false, err
	}
	return true, nil
}

func createRecipe(tx *gorm.DB, r demoRecipe, userID string, categoryID uint, ingredients map[string]uint) error {
	recipe := r.recipe
	recipe.UserID = userID
	recipe.CategoryID = categoryID
	if err := tx.Create(&recipe).Error; err != nil {
		return fmt.Errorf("criar receita %q: %w", recipe.Name, err)
	}

	for _, url := range r.images {
		if err := tx.Create(&model.ImageRecipe{URL: url, RecipeID: recipe.ID}).Error; err != nil {
			return fmt.Errorf("criar foto de %q: %w", recipe.Name, err)
		}
	}

	for name, amount := range r.ingredients {
		id, err := ingredientID(tx, ingredients, name)
		if err != nil {
			return err
		}
		item := model.IngredientRecipe{Amount: amount, RecipeID: recipe.ID, IngredientID: id}
		if err := tx.Create(&item).Error; err != nil {
			return fmt.Errorf("criar ingrediente de %q: %w", recipe.Name, err)
		}
	}

	for _, step := range r.preparations {
		if err := tx.Create(&model.Preparation{Description: step, RecipeID: recipe.ID}).Error; err != nil {
			return fmt.Errorf("criar passo de %q: %w", recipe.Name, err)
		}
	}
	return nil
}

// ingredientID devolve o id do ingrediente, criando se ainda não existir.
func ingredientID(tx *gorm.DB, ingredients map[string]uint, name string) (uint, error) {
	if id, ok := ingredients[name]; ok {
		return id, nil
	}
	ing := model.Ingredient{Name: name}
	if err := tx.Create(&ing).Error; err != nil {
		return 0, fmt.Errorf("criar ingrediente %q: %w", name, err)
	}
	ingredients[name] = ing.ID
	return ing.ID, nil
}
