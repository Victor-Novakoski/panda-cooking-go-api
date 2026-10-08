// Package seed cria os dados de demonstração: três usuários que conseguem
// entrar, dez receitas brasileiras, comentários e favoritos.
package seed

import (
	"context"
	"fmt"
	"time"

	"panda-cooking-go-api/internal/model"
	"panda-cooking-go-api/internal/service"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

// DemoPassword é a senha dos usuários de demonstração. É pública de
// propósito (está no README) para quem visita o projeto conseguir entrar.
const DemoPassword = "panda-cooking-demo" //nolint:gosec // senha pública dos usuários de demonstração

// As categorias ficam de fora: são fixas e vêm da migration.
const truncateAll = "TRUNCATE TABLE refresh_tokens, sessions, favorite_recipes, comments, preparations, " +
	"ingredient_recipes, image_recipes, recipes, ingredients, users RESTART IDENTITY CASCADE"

// Demo popula o banco com os dados de demonstração. Com reset=false só age se
// ainda não houver usuário nenhum, para não apagar o que já foi cadastrado;
// com reset=true apaga tudo e recria. Diz se populou.
func Demo(ctx context.Context, db *gorm.DB, reset bool) (bool, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(DemoPassword), bcrypt.DefaultCost)
	if err != nil {
		return false, fmt.Errorf("gerar hash da senha: %w", err)
	}

	created := false
	err = db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		// duas instâncias subindo juntas não populam o banco duas vezes
		if err := tx.Exec("SELECT pg_advisory_xact_lock(hashtext('panda-cooking-seed'))").Error; err != nil {
			return fmt.Errorf("travar o seed: %w", err)
		}
		if !reset {
			var users int64
			if err := tx.Model(&model.User{}).Count(&users).Error; err != nil {
				return fmt.Errorf("contar usuários: %w", err)
			}
			if users > 0 {
				return nil
			}
		}
		if err := tx.Exec(truncateAll).Error; err != nil {
			return fmt.Errorf("limpar tabelas: %w", err)
		}
		if err := populate(tx, string(hash), time.Now()); err != nil {
			return err
		}
		created = true
		return nil
	})
	if err != nil {
		return false, err
	}
	return created, nil
}

func populate(tx *gorm.DB, passwordHash string, now time.Time) error {
	var categories []model.Category
	if err := tx.Find(&categories).Error; err != nil {
		return fmt.Errorf("ler categorias: %w", err)
	}
	categoryIDs := make(map[string]uint, len(categories))
	for _, c := range categories {
		categoryIDs[c.Name] = c.ID
	}

	// tudo "aconteceu" nos últimos dias, a receita mais nova primeiro
	start := now.Add(-time.Duration(len(demoRecipes)+1) * 24 * time.Hour)

	users := make([]string, len(demoUsers))
	for i, u := range demoUsers {
		user := model.User{
			Name:         u.name,
			Email:        u.email,
			PasswordHash: passwordHash,
			ImageProfile: u.image,
			IsAdm:        u.isAdm,
			CreatedAt:    start,
			UpdatedAt:    start,
		}
		if err := tx.Create(&user).Error; err != nil {
			return fmt.Errorf("criar usuário %q: %w", u.email, err)
		}
		users[i] = user.ID
	}

	ingredients := map[string]uint{}
	recipes := make(map[string]string, len(demoRecipes))
	for i, r := range demoRecipes {
		categoryID, ok := categoryIDs[r.category]
		if !ok {
			return fmt.Errorf("receita %q: categoria %q não existe", r.recipe.Name, r.category)
		}
		createdAt := now.Add(-time.Duration(i+1) * 22 * time.Hour)
		id, err := createRecipe(tx, r, users, categoryID, createdAt, ingredients)
		if err != nil {
			return err
		}
		recipes[r.recipe.Name] = id
	}

	for i, f := range demoFavorites {
		fav := model.FavoriteRecipe{
			UserID:    users[f.user],
			RecipeID:  recipes[f.recipe],
			CreatedAt: now.Add(-time.Duration(len(demoFavorites)-i) * time.Hour),
		}
		if err := tx.Create(&fav).Error; err != nil {
			return fmt.Errorf("criar favorito de %q: %w", f.recipe, err)
		}
	}
	return nil
}

func createRecipe(tx *gorm.DB, r demoRecipe, users []string, categoryID uint, createdAt time.Time, ingredients map[string]uint) (string, error) {
	recipe := r.recipe
	recipe.UserID = users[r.author]
	recipe.CategoryID = categoryID
	recipe.CreatedAt = createdAt
	recipe.UpdatedAt = createdAt
	if err := tx.Omit("User", "Category", "Images", "Ingredients", "Preparations").Create(&recipe).Error; err != nil {
		return "", fmt.Errorf("criar receita %q: %w", recipe.Name, err)
	}

	for _, url := range r.images {
		if err := tx.Create(&model.ImageRecipe{URL: url, RecipeID: recipe.ID}).Error; err != nil {
			return "", fmt.Errorf("criar foto de %q: %w", recipe.Name, err)
		}
	}

	for _, ing := range r.ingredients {
		id, err := ingredientID(tx, ingredients, service.NormalizeIngredient(ing.name))
		if err != nil {
			return "", err
		}
		item := model.IngredientRecipe{Amount: ing.amount, RecipeID: recipe.ID, IngredientID: id}
		if err := tx.Omit("Ingredient").Create(&item).Error; err != nil {
			return "", fmt.Errorf("criar ingrediente de %q: %w", recipe.Name, err)
		}
	}

	for _, step := range r.preparations {
		if err := tx.Create(&model.Preparation{Description: step, RecipeID: recipe.ID}).Error; err != nil {
			return "", fmt.Errorf("criar passo de %q: %w", recipe.Name, err)
		}
	}

	for i, c := range r.comments {
		at := createdAt.Add(time.Duration(i+1) * 3 * time.Hour)
		comment := model.Comment{Description: c.text, UserID: users[c.author], RecipeID: recipe.ID, CreatedAt: at, UpdatedAt: at}
		if err := tx.Omit("User").Create(&comment).Error; err != nil {
			return "", fmt.Errorf("criar comentário de %q: %w", recipe.Name, err)
		}
	}
	return recipe.ID, nil
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
