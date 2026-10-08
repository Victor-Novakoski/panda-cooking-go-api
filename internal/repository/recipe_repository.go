package repository

import (
	"context"
	"strings"

	"panda-cooking-go-api/internal/model"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type RecipeRepository struct {
	db *gorm.DB
}

func NewRecipeRepository(db *gorm.DB) *RecipeRepository {
	return &RecipeRepository{db: db}
}

// byID mantém fotos, ingredientes e passos na ordem em que foram criados;
// sem ORDER BY o Postgres não garante ordem nenhuma.
func byID(db *gorm.DB) *gorm.DB { return db.Order("id") }

// preloadSummary carrega o que o card da receita mostra.
func preloadSummary(db *gorm.DB) *gorm.DB {
	return db.Preload("User", selectAuthor).Preload("Category").Preload("Images", byID)
}

func (r *RecipeRepository) Create(ctx context.Context, recipe *model.Recipe) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := tx.Omit(clause.Associations).Create(recipe).Error; err != nil {
			return err
		}
		return createItems(tx, recipe)
	})
}

// List devolve uma página de receitas, das mais novas para as mais antigas,
// e o total que bate com os filtros.
func (r *RecipeRepository) List(ctx context.Context, filter RecipeFilter, page Page) ([]model.Recipe, int64, error) {
	db := r.db.WithContext(ctx).Model(&model.Recipe{})
	if filter.Search != "" {
		// mesma expressão do índice recipes_search_idx (migration 000003)
		db = db.Where("search_normalize(name || ' ' || description) LIKE '%' || search_normalize(?) || '%'",
			escapeLike(filter.Search))
	}
	if filter.CategoryID != 0 {
		db = db.Where("category_id = ?", filter.CategoryID)
	}
	if filter.UserID != "" {
		if !IsUUID(filter.UserID) {
			return nil, 0, nil
		}
		db = db.Where("user_id = ?", filter.UserID)
	}
	db = db.Session(&gorm.Session{})

	var total int64
	if err := db.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	var recipes []model.Recipe
	err := preloadSummary(db).
		Order("created_at DESC, id DESC").
		Limit(page.Size).Offset(page.offset()).
		Find(&recipes).Error
	return recipes, total, err
}

func (r *RecipeRepository) FindByID(ctx context.Context, id string) (*model.Recipe, error) {
	if !IsUUID(id) {
		return nil, gorm.ErrRecordNotFound
	}
	var recipe model.Recipe
	err := preloadSummary(r.db.WithContext(ctx)).
		Preload("Ingredients", byID).
		Preload("Ingredients.Ingredient").
		Preload("Preparations", byID).
		First(&recipe, "id = ?", id).Error
	if err != nil {
		return nil, err
	}
	return &recipe, nil
}

func (r *RecipeRepository) Exists(ctx context.Context, id string) (bool, error) {
	if !IsUUID(id) {
		return false, nil
	}
	var count int64
	err := r.db.WithContext(ctx).Model(&model.Recipe{}).Where("id = ?", id).Count(&count).Error
	return count > 0, err
}

// Update grava só as colunas da receita. Sem o Select, o GORM também
// gravava a Category carregada pelo Preload e voltava o category_id.
func (r *RecipeRepository) Update(ctx context.Context, recipe *model.Recipe) error {
	return r.db.WithContext(ctx).Model(&model.Recipe{ID: recipe.ID}).
		Select("Name", "Description", "Time", "Portions", "CategoryID", "UpdatedAt").
		Updates(recipe).Error
}

// Replace troca dados, fotos, ingredientes e passos da receita numa transação.
// Os itens antigos são apagados e os novos criados na ordem recebida.
func (r *RecipeRepository) Replace(ctx context.Context, recipe *model.Recipe) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := r.withTx(tx).Update(ctx, recipe); err != nil {
			return err
		}
		for _, item := range []any{&model.ImageRecipe{}, &model.IngredientRecipe{}, &model.Preparation{}} {
			if err := tx.Where("recipe_id = ?", recipe.ID).Delete(item).Error; err != nil {
				return err
			}
		}
		return createItems(tx, recipe)
	})
}

func (r *RecipeRepository) withTx(tx *gorm.DB) *RecipeRepository { return &RecipeRepository{db: tx} }

// createItems grava fotos, passos e ingredientes da receita, na ordem.
// O ingrediente é achado pelo nome ou criado.
func createItems(tx *gorm.DB, recipe *model.Recipe) error {
	for i := range recipe.Images {
		recipe.Images[i].RecipeID = recipe.ID
	}
	if len(recipe.Images) > 0 {
		if err := tx.Create(&recipe.Images).Error; err != nil {
			return err
		}
	}

	for i := range recipe.Preparations {
		recipe.Preparations[i].RecipeID = recipe.ID
	}
	if len(recipe.Preparations) > 0 {
		if err := tx.Create(&recipe.Preparations).Error; err != nil {
			return err
		}
	}

	for i := range recipe.Ingredients {
		ir := &recipe.Ingredients[i]
		if err := findOrCreateIngredient(tx, &ir.Ingredient); err != nil {
			return err
		}
		ir.RecipeID = recipe.ID
		ir.IngredientID = ir.Ingredient.ID
	}
	if len(recipe.Ingredients) > 0 {
		if err := tx.Omit(clause.Associations).Create(&recipe.Ingredients).Error; err != nil {
			return err
		}
	}
	return nil
}

// findOrCreateIngredient usa ON CONFLICT para duas receitas novas com o mesmo
// ingrediente ao mesmo tempo não esbarrarem no índice único do nome.
func findOrCreateIngredient(tx *gorm.DB, ing *model.Ingredient) error {
	err := tx.Clauses(clause.OnConflict{Columns: []clause.Column{{Name: "name"}}, DoNothing: true}).
		Create(ing).Error
	if err != nil {
		return err
	}
	if ing.ID == 0 { // já existia: o DO NOTHING não devolve o id
		return tx.Where("name = ?", ing.Name).First(ing).Error
	}
	return nil
}

// Delete apaga a receita; fotos, itens, comentários e favoritos vão junto
// (ON DELETE CASCADE).
func (r *RecipeRepository) Delete(ctx context.Context, id string) error {
	return r.db.WithContext(ctx).Delete(&model.Recipe{}, "id = ?", id).Error
}

// --- Imagens ---

func (r *RecipeRepository) AddImage(ctx context.Context, image *model.ImageRecipe) error {
	return r.db.WithContext(ctx).Create(image).Error
}

func (r *RecipeRepository) FindImageByID(ctx context.Context, id uint) (*model.ImageRecipe, error) {
	var image model.ImageRecipe
	err := r.db.WithContext(ctx).First(&image, id).Error
	if err != nil {
		return nil, err
	}
	return &image, nil
}

func (r *RecipeRepository) UpdateImage(ctx context.Context, image *model.ImageRecipe) error {
	return r.db.WithContext(ctx).Model(image).Update("url", image.URL).Error
}

func (r *RecipeRepository) DeleteImage(ctx context.Context, id uint) error {
	return r.db.WithContext(ctx).Delete(&model.ImageRecipe{}, id).Error
}

// --- Ingredientes ---

func (r *RecipeRepository) AddIngredient(ctx context.Context, ir *model.IngredientRecipe) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := findOrCreateIngredient(tx, &ir.Ingredient); err != nil {
			return err
		}
		ir.IngredientID = ir.Ingredient.ID
		return tx.Omit(clause.Associations).Create(ir).Error
	})
}

func (r *RecipeRepository) FindIngredientRecipeByID(ctx context.Context, id uint) (*model.IngredientRecipe, error) {
	var ir model.IngredientRecipe
	err := r.db.WithContext(ctx).First(&ir, id).Error
	if err != nil {
		return nil, err
	}
	return &ir, nil
}

func (r *RecipeRepository) DeleteIngredientRecipe(ctx context.Context, id uint) error {
	return r.db.WithContext(ctx).Delete(&model.IngredientRecipe{}, id).Error
}

// --- Preparos ---

func (r *RecipeRepository) AddPreparation(ctx context.Context, p *model.Preparation) error {
	return r.db.WithContext(ctx).Create(p).Error
}

func (r *RecipeRepository) FindPreparationByID(ctx context.Context, id uint) (*model.Preparation, error) {
	var p model.Preparation
	err := r.db.WithContext(ctx).First(&p, id).Error
	if err != nil {
		return nil, err
	}
	return &p, nil
}

func (r *RecipeRepository) UpdatePreparation(ctx context.Context, p *model.Preparation) error {
	return r.db.WithContext(ctx).Model(p).Update("description", p.Description).Error
}

func (r *RecipeRepository) DeletePreparation(ctx context.Context, id uint) error {
	return r.db.WithContext(ctx).Delete(&model.Preparation{}, id).Error
}

// escapeLike faz %, _ e \ do termo buscado valerem como texto no LIKE.
func escapeLike(s string) string {
	return strings.NewReplacer(`\`, `\\`, `%`, `\%`, `_`, `\_`).Replace(s)
}
