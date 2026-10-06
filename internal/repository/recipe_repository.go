package repository

import (
	"panda-cooking-go-api/internal/model"

	"gorm.io/gorm"
)

type RecipeRepository struct {
	db *gorm.DB
}

func NewRecipeRepository(db *gorm.DB) *RecipeRepository {
	return &RecipeRepository{db: db}
}

func (r *RecipeRepository) Create(recipe *model.Recipe) error {
	return r.db.Create(recipe).Error
}

func (r *RecipeRepository) FindAll() ([]model.Recipe, error) {
	var recipes []model.Recipe
	err := r.db.
		Preload("User").
		Preload("Category").
		Preload("Images").
		Preload("Ingredients").
		Preload("Ingredients.Ingredient").
		Preload("Preparations").
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
		Preload("Images").
		Preload("Ingredients").
		Preload("Ingredients.Ingredient").
		Preload("Preparations").
		First(&recipe, "id = ?", id).Error
	if err != nil {
		return nil, err
	}
	return &recipe, nil
}

func (r *RecipeRepository) Update(recipe *model.Recipe) error {
	return r.db.Save(recipe).Error
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
