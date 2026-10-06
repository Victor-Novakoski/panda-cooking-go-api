package repository

import (
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

func (r *RecipeRepository) Create(recipe *model.Recipe) error {
	return r.db.Create(recipe).Error
}

func (r *RecipeRepository) FindAll() ([]model.Recipe, error) {
	var recipes []model.Recipe
	err := r.db.
		Preload("User").
		Preload("Category").
		Preload("Images", byID).
		Preload("Ingredients", byID).
		Preload("Ingredients.Ingredient").
		Preload("Preparations", byID).
		Find(&recipes).Error
	return recipes, err
}

func (r *RecipeRepository) FindByID(id string) (*model.Recipe, error) {
	if !isUUID(id) {
		return nil, gorm.ErrRecordNotFound
	}
	var recipe model.Recipe
	err := r.db.
		Preload("User").
		Preload("Category").
		Preload("Images", byID).
		Preload("Ingredients", byID).
		Preload("Ingredients.Ingredient").
		Preload("Preparations", byID).
		First(&recipe, "id = ?", id).Error
	if err != nil {
		return nil, err
	}
	return &recipe, nil
}

// Update grava só as colunas da receita. Sem o Omit, o Save também gravava a
// Category carregada pelo Preload e voltava o category_id para o antigo.
func (r *RecipeRepository) Update(recipe *model.Recipe) error {
	return r.db.Omit(clause.Associations).Save(recipe).Error
}

// Replace troca dados, fotos, ingredientes e passos da receita numa transação.
// Os itens antigos são apagados e os novos criados na ordem recebida.
func (r *RecipeRepository) Replace(recipe *model.Recipe) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		err := tx.Model(&model.Recipe{ID: recipe.ID}).
			Select("Name", "Description", "Time", "Portions", "CategoryID").
			Updates(recipe).Error
		if err != nil {
			return err
		}

		for _, item := range []any{&model.ImageRecipe{}, &model.IngredientRecipe{}, &model.Preparation{}} {
			if err := tx.Where("recipe_id = ?", recipe.ID).Delete(item).Error; err != nil {
				return err
			}
		}

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
			name := ir.Ingredient.Name
			if err := tx.Where("name = ?", name).FirstOrCreate(&ir.Ingredient, model.Ingredient{Name: name}).Error; err != nil {
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
	})
}

func (r *RecipeRepository) Delete(id string) error {
	return r.db.Delete(&model.Recipe{}, "id = ?", id).Error
}

// --- Imagens ---

func (r *RecipeRepository) AddImage(image *model.ImageRecipe) error {
	return r.db.Create(image).Error
}

func (r *RecipeRepository) FindImageByID(id uint) (*model.ImageRecipe, error) {
	var image model.ImageRecipe
	err := r.db.First(&image, id).Error
	if err != nil {
		return nil, err
	}
	return &image, nil
}

func (r *RecipeRepository) UpdateImage(image *model.ImageRecipe) error {
	return r.db.Save(image).Error
}

func (r *RecipeRepository) DeleteImage(id uint) error {
	return r.db.Delete(&model.ImageRecipe{}, id).Error
}

// --- Ingredientes ---

func (r *RecipeRepository) FindOrCreateIngredient(name string) (*model.Ingredient, error) {
	var ingredient model.Ingredient
	err := r.db.Where("name = ?", name).FirstOrCreate(&ingredient, model.Ingredient{Name: name}).Error
	if err != nil {
		return nil, err
	}
	return &ingredient, nil
}

func (r *RecipeRepository) AddIngredient(ir *model.IngredientRecipe) error {
	return r.db.Create(ir).Error
}

func (r *RecipeRepository) FindIngredientRecipeByID(id uint) (*model.IngredientRecipe, error) {
	var ir model.IngredientRecipe
	err := r.db.First(&ir, id).Error
	if err != nil {
		return nil, err
	}
	return &ir, nil
}

func (r *RecipeRepository) DeleteIngredientRecipe(id uint) error {
	return r.db.Delete(&model.IngredientRecipe{}, id).Error
}

// --- Preparos ---

func (r *RecipeRepository) AddPreparation(p *model.Preparation) error {
	return r.db.Create(p).Error
}

func (r *RecipeRepository) FindPreparationByID(id uint) (*model.Preparation, error) {
	var p model.Preparation
	err := r.db.First(&p, id).Error
	if err != nil {
		return nil, err
	}
	return &p, nil
}

func (r *RecipeRepository) UpdatePreparation(p *model.Preparation) error {
	return r.db.Save(p).Error
}

func (r *RecipeRepository) DeletePreparation(id uint) error {
	return r.db.Delete(&model.Preparation{}, id).Error
}
