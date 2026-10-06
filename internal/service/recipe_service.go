package service

import (
	"errors"
	"strings"

	"panda-cooking-go-api/internal/model"
	"panda-cooking-go-api/internal/repository"

	"gorm.io/gorm"
)

type RecipeService struct {
	repo repository.RecipeRepo
}

func NewRecipeService(repo repository.RecipeRepo) *RecipeService {
	return &RecipeService{repo: repo}
}

// --- DTOs de entrada ---

type CreateRecipeInput struct {
	Name         string                  `json:"name" binding:"required"`
	Description  string                  `json:"description" binding:"required"`
	Time         string                  `json:"time" binding:"required"`
	Portions     int                     `json:"portions" binding:"required,min=1"`
	CategoryID   uint                    `json:"category_id" binding:"required"`
	Images       []ImageRecipeInput      `json:"images"`
	Ingredients  []IngredientRecipeInput `json:"ingredients"`
	Preparations []PreparationInput      `json:"preparations"`
}

type UpdateRecipeInput struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	Time        string `json:"time"`
	Portions    int    `json:"portions"`
	CategoryID  uint   `json:"category_id"`
}

type ImageRecipeInput struct {
	URL string `json:"url" binding:"required,url"`
}

type IngredientRecipeInput struct {
	Name   string `json:"name" binding:"required"`
	Amount string `json:"amount" binding:"required"`
}

type PreparationInput struct {
	Description string `json:"description" binding:"required"`
}

// --- DTOs de saída ---

type RecipeResponse struct {
	ID           string                     `json:"id"`
	Name         string                     `json:"name"`
	Description  string                     `json:"description"`
	Time         string                     `json:"time"`
	Portions     int                        `json:"portions"`
	UserID       string                     `json:"user_id"`
	Category     CategoryResponse           `json:"category"`
	Images       []ImageRecipeResponse      `json:"images"`
	Ingredients  []IngredientRecipeResponse `json:"ingredients"`
	Preparations []PreparationResponse      `json:"preparations"`
}

type CategoryResponse struct {
	ID   uint   `json:"id"`
	Name string `json:"name"`
}

type ImageRecipeResponse struct {
	ID  uint   `json:"id"`
	URL string `json:"url"`
}

type IngredientRecipeResponse struct {
	ID     uint   `json:"id"`
	Amount string `json:"amount"`
	Name   string `json:"name"`
}

type PreparationResponse struct {
	ID          uint   `json:"id"`
	Description string `json:"description"`
}

// --- Métodos de receita ---

func (s *RecipeService) Create(userID string, input CreateRecipeInput) (*RecipeResponse, error) {
	recipe := &model.Recipe{
		Name:        input.Name,
		Description: input.Description,
		Time:        input.Time,
		Portions:    input.Portions,
		UserID:      userID,
		CategoryID:  input.CategoryID,
	}

	for _, img := range input.Images {
		recipe.Images = append(recipe.Images, model.ImageRecipe{URL: img.URL})
	}

	for _, prep := range input.Preparations {
		recipe.Preparations = append(recipe.Preparations, model.Preparation{Description: prep.Description})
	}

	if err := s.repo.Create(recipe); err != nil {
		return nil, err
	}

	// ingredientes precisam de tratamento especial (FirstOrCreate)
	for _, ing := range input.Ingredients {
		ingredient, err := s.repo.FindOrCreateIngredient(normalizeIngredient(ing.Name))
		if err != nil {
			return nil, err
		}
		ir := &model.IngredientRecipe{
			RecipeID:     recipe.ID,
			IngredientID: ingredient.ID,
			Amount:       ing.Amount,
		}
		if err := s.repo.AddIngredient(ir); err != nil {
			return nil, err
		}
	}

	return s.GetByID(recipe.ID)
}

func (s *RecipeService) GetAll() ([]RecipeResponse, error) {
	recipes, err := s.repo.FindAll()
	if err != nil {
		return nil, err
	}

	result := make([]RecipeResponse, len(recipes))
	for i, r := range recipes {
		result[i] = toRecipeResponse(r)
	}
	return result, nil
}

func (s *RecipeService) GetByID(id string) (*RecipeResponse, error) {
	recipe, err := s.repo.FindByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("receita não encontrada")
		}
		return nil, err
	}
	resp := toRecipeResponse(*recipe)
	return &resp, nil
}

func (s *RecipeService) Update(recipeID, userID string, input UpdateRecipeInput) (*RecipeResponse, error) {
	recipe, err := s.repo.FindByID(recipeID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("receita não encontrada")
		}
		return nil, err
	}

	if recipe.UserID != userID {
		return nil, errors.New("sem permissão para editar esta receita")
	}

	if input.Name != "" {
		recipe.Name = input.Name
	}
	if input.Description != "" {
		recipe.Description = input.Description
	}
	if input.Time != "" {
		recipe.Time = input.Time
	}
	if input.Portions > 0 {
		recipe.Portions = input.Portions
	}
	if input.CategoryID > 0 {
		recipe.CategoryID = input.CategoryID
	}

	if err := s.repo.Update(recipe); err != nil {
		return nil, err
	}

	return s.GetByID(recipe.ID)
}

func (s *RecipeService) Delete(recipeID, userID string) error {
	recipe, err := s.repo.FindByID(recipeID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("receita não encontrada")
		}
		return err
	}

	if recipe.UserID != userID {
		return errors.New("sem permissão para deletar esta receita")
	}

	return s.repo.Delete(recipeID)
}

// --- Imagens ---

func (s *RecipeService) AddImage(recipeID, userID string, input ImageRecipeInput) (*ImageRecipeResponse, error) {
	recipe, err := s.repo.FindByID(recipeID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("receita não encontrada")
		}
		return nil, err
	}

	if recipe.UserID != userID {
		return nil, errors.New("sem permissão")
	}

	image := &model.ImageRecipe{URL: input.URL, RecipeID: recipeID}
	if err := s.repo.AddImage(image); err != nil {
		return nil, err
	}

	return &ImageRecipeResponse{ID: image.ID, URL: image.URL}, nil
}

func (s *RecipeService) UpdateImage(recipeID, userID string, imageID uint, input ImageRecipeInput) (*ImageRecipeResponse, error) {
	recipe, err := s.repo.FindByID(recipeID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("receita não encontrada")
		}
		return nil, err
	}

	if recipe.UserID != userID {
		return nil, errors.New("sem permissão")
	}

	image, err := s.repo.FindImageByID(imageID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("imagem não encontrada")
		}
		return nil, err
	}

	image.URL = input.URL
	if err := s.repo.UpdateImage(image); err != nil {
		return nil, err
	}

	return &ImageRecipeResponse{ID: image.ID, URL: image.URL}, nil
}

func (s *RecipeService) DeleteImage(recipeID, userID string, imageID uint) error {
	recipe, err := s.repo.FindByID(recipeID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("receita não encontrada")
		}
		return err
	}

	if recipe.UserID != userID {
		return errors.New("sem permissão")
	}

	if _, err := s.repo.FindImageByID(imageID); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("imagem não encontrada")
		}
		return err
	}

	return s.repo.DeleteImage(imageID)
}

// --- Ingredientes ---

func (s *RecipeService) AddIngredient(recipeID, userID string, input IngredientRecipeInput) (*IngredientRecipeResponse, error) {
	recipe, err := s.repo.FindByID(recipeID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("receita não encontrada")
		}
		return nil, err
	}

	if recipe.UserID != userID {
		return nil, errors.New("sem permissão")
	}

	ingredient, err := s.repo.FindOrCreateIngredient(normalizeIngredient(input.Name))
	if err != nil {
		return nil, err
	}

	ir := &model.IngredientRecipe{
		RecipeID:     recipeID,
		IngredientID: ingredient.ID,
		Amount:       input.Amount,
	}
	if err := s.repo.AddIngredient(ir); err != nil {
		return nil, err
	}

	return &IngredientRecipeResponse{ID: ir.ID, Amount: ir.Amount, Name: ingredient.Name}, nil
}

func (s *RecipeService) DeleteIngredient(recipeID, userID string, ingredientRecipeID uint) error {
	recipe, err := s.repo.FindByID(recipeID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("receita não encontrada")
		}
		return err
	}

	if recipe.UserID != userID {
		return errors.New("sem permissão")
	}

	if _, err := s.repo.FindIngredientRecipeByID(ingredientRecipeID); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("ingrediente não encontrado nesta receita")
		}
		return err
	}

	return s.repo.DeleteIngredientRecipe(ingredientRecipeID)
}

// --- Preparos ---

func (s *RecipeService) AddPreparation(recipeID, userID string, input PreparationInput) (*PreparationResponse, error) {
	recipe, err := s.repo.FindByID(recipeID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("receita não encontrada")
		}
		return nil, err
	}

	if recipe.UserID != userID {
		return nil, errors.New("sem permissão")
	}

	p := &model.Preparation{Description: input.Description, RecipeID: recipeID}
	if err := s.repo.AddPreparation(p); err != nil {
		return nil, err
	}

	return &PreparationResponse{ID: p.ID, Description: p.Description}, nil
}

func (s *RecipeService) UpdatePreparation(recipeID, userID string, prepID uint, input PreparationInput) (*PreparationResponse, error) {
	recipe, err := s.repo.FindByID(recipeID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("receita não encontrada")
		}
		return nil, err
	}

	if recipe.UserID != userID {
		return nil, errors.New("sem permissão")
	}

	p, err := s.repo.FindPreparationByID(prepID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("passo de preparo não encontrado")
		}
		return nil, err
	}

	p.Description = input.Description
	if err := s.repo.UpdatePreparation(p); err != nil {
		return nil, err
	}

	return &PreparationResponse{ID: p.ID, Description: p.Description}, nil
}

func (s *RecipeService) DeletePreparation(recipeID, userID string, prepID uint) error {
	recipe, err := s.repo.FindByID(recipeID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("receita não encontrada")
		}
		return err
	}

	if recipe.UserID != userID {
		return errors.New("sem permissão")
	}

	if _, err := s.repo.FindPreparationByID(prepID); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("passo de preparo não encontrado")
		}
		return err
	}

	return s.repo.DeletePreparation(prepID)
}

func normalizeIngredient(name string) string {
	return strings.ToLower(strings.TrimSpace(name))
}

// --- Conversor model → DTO ---

func toRecipeResponse(r model.Recipe) RecipeResponse {
	resp := RecipeResponse{
		ID:          r.ID,
		Name:        r.Name,
		Description: r.Description,
		Time:        r.Time,
		Portions:    r.Portions,
		UserID:      r.UserID,
		Category:    CategoryResponse{ID: r.Category.ID, Name: r.Category.Name},
	}

	for _, img := range r.Images {
		resp.Images = append(resp.Images, ImageRecipeResponse{ID: img.ID, URL: img.URL})
	}

	for _, ir := range r.Ingredients {
		resp.Ingredients = append(resp.Ingredients, IngredientRecipeResponse{
			ID:     ir.ID,
			Amount: ir.Amount,
			Name:   ir.Ingredient.Name,
		})
	}

	for _, p := range r.Preparations {
		resp.Preparations = append(resp.Preparations, PreparationResponse{ID: p.ID, Description: p.Description})
	}

	return resp
}
