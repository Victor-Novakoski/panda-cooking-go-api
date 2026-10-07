package service

import (
	"context"
	"errors"
	"strings"

	"panda-cooking-go-api/internal/model"
	"panda-cooking-go-api/internal/repository"

	"gorm.io/gorm"
)

// Limites de itens por receita, iguais aos da validação do corpo (RecipeInput).
const (
	MaxImages       = 10
	MaxIngredients  = 50
	MaxPreparations = 50
)

type RecipeService struct {
	repo       repository.RecipeRepo
	categories repository.CategoryRepo
}

func NewRecipeService(repo repository.RecipeRepo, categories repository.CategoryRepo) *RecipeService {
	return &RecipeService{repo: repo, categories: categories}
}

// --- Receita ---

func (s *RecipeService) Create(ctx context.Context, userID string, input RecipeInput) (*RecipeResponse, error) {
	if err := s.checkCategory(ctx, input.CategoryID); err != nil {
		return nil, err
	}

	recipe := toRecipeModel(input)
	recipe.UserID = userID
	if err := s.repo.Create(ctx, recipe); err != nil {
		// a categoria já foi conferida: sobra a conta apagada com o access
		// token ainda válido
		if errors.Is(err, gorm.ErrForeignKeyViolated) {
			return nil, ErrUserNotFound
		}
		return nil, err
	}
	return s.GetByID(ctx, recipe.ID)
}

func (s *RecipeService) List(ctx context.Context, q ListRecipesQuery) (Page[RecipeSummary], error) {
	page := PageQuery{Page: q.Page, PerPage: q.PerPage}.repo(DefaultRecipesPerPage)
	filter := repository.RecipeFilter{Search: q.Search, CategoryID: q.CategoryID, UserID: q.UserID}

	recipes, total, err := s.repo.List(ctx, filter, page)
	if err != nil {
		return Page[RecipeSummary]{}, err
	}
	return newPage(toRecipeSummaries(recipes), page, total), nil
}

func (s *RecipeService) GetByID(ctx context.Context, id string) (*RecipeResponse, error) {
	recipe, err := s.find(ctx, id)
	if err != nil {
		return nil, err
	}
	resp := toRecipeResponse(*recipe)
	return &resp, nil
}

// Update (PATCH) muda só os campos enviados.
func (s *RecipeService) Update(ctx context.Context, recipeID, userID string, input UpdateRecipeInput) (*RecipeResponse, error) {
	recipe, err := s.findOwned(ctx, recipeID, userID, ErrForbiddenEditRecipe)
	if err != nil {
		return nil, err
	}

	if input.Name != nil {
		recipe.Name = *input.Name
	}
	if input.Description != nil {
		recipe.Description = *input.Description
	}
	if input.Time != nil {
		recipe.Time = *input.Time
	}
	if input.Portions != nil {
		recipe.Portions = *input.Portions
	}
	if input.CategoryID != nil && *input.CategoryID != recipe.CategoryID {
		if err := s.checkCategory(ctx, *input.CategoryID); err != nil {
			return nil, err
		}
		recipe.CategoryID = *input.CategoryID
	}

	if err := s.repo.Update(ctx, recipe); err != nil {
		return nil, err
	}
	return s.GetByID(ctx, recipe.ID)
}

// Replace (PUT) troca a receita inteira de uma vez (dados, fotos,
// ingredientes e passos). É o que a tela de edição usa: numa transação só,
// ou muda tudo ou nada, sem deixar a receita pela metade.
func (s *RecipeService) Replace(ctx context.Context, recipeID, userID string, input RecipeInput) (*RecipeResponse, error) {
	recipe, err := s.findOwned(ctx, recipeID, userID, ErrForbiddenEditRecipe)
	if err != nil {
		return nil, err
	}
	if err := s.checkCategory(ctx, input.CategoryID); err != nil {
		return nil, err
	}

	updated := toRecipeModel(input)
	updated.ID = recipe.ID
	updated.UserID = recipe.UserID
	if err := s.repo.Replace(ctx, updated); err != nil {
		return nil, err
	}
	return s.GetByID(ctx, recipe.ID)
}

func (s *RecipeService) Delete(ctx context.Context, recipeID, userID string) error {
	if _, err := s.findOwned(ctx, recipeID, userID, ErrForbiddenDelRecipe); err != nil {
		return err
	}
	return s.repo.Delete(ctx, recipeID)
}

// --- Imagens ---

func (s *RecipeService) AddImage(ctx context.Context, recipeID, userID string, input ImageRecipeInput) (*ImageRecipeResponse, error) {
	recipe, err := s.findOwned(ctx, recipeID, userID, ErrForbidden)
	if err != nil {
		return nil, err
	}
	if len(recipe.Images) >= MaxImages {
		return nil, ErrTooManyImages
	}

	image := &model.ImageRecipe{URL: input.URL, RecipeID: recipe.ID}
	if err := s.repo.AddImage(ctx, image); err != nil {
		return nil, err
	}
	return &ImageRecipeResponse{ID: image.ID, URL: image.URL}, nil
}

func (s *RecipeService) UpdateImage(ctx context.Context, recipeID, userID string, imageID uint, input ImageRecipeInput) (*ImageRecipeResponse, error) {
	recipe, err := s.findOwned(ctx, recipeID, userID, ErrForbidden)
	if err != nil {
		return nil, err
	}

	image, err := s.repo.FindImageByID(ctx, imageID)
	// a imagem precisa ser desta receita, senão o dono de uma receita mexeria em outra
	if notFound(err) || (err == nil && image.RecipeID != recipe.ID) {
		return nil, ErrImageNotFound
	}
	if err != nil {
		return nil, err
	}

	image.URL = input.URL
	if err := s.repo.UpdateImage(ctx, image); err != nil {
		return nil, err
	}
	return &ImageRecipeResponse{ID: image.ID, URL: image.URL}, nil
}

func (s *RecipeService) DeleteImage(ctx context.Context, recipeID, userID string, imageID uint) error {
	recipe, err := s.findOwned(ctx, recipeID, userID, ErrForbidden)
	if err != nil {
		return err
	}

	image, err := s.repo.FindImageByID(ctx, imageID)
	if notFound(err) || (err == nil && image.RecipeID != recipe.ID) {
		return ErrImageNotFound
	}
	if err != nil {
		return err
	}
	return s.repo.DeleteImage(ctx, imageID)
}

// --- Ingredientes ---

func (s *RecipeService) AddIngredient(ctx context.Context, recipeID, userID string, input IngredientRecipeInput) (*IngredientRecipeResponse, error) {
	recipe, err := s.findOwned(ctx, recipeID, userID, ErrForbidden)
	if err != nil {
		return nil, err
	}
	if len(recipe.Ingredients) >= MaxIngredients {
		return nil, ErrTooManyItems
	}

	ir := &model.IngredientRecipe{
		RecipeID:   recipe.ID,
		Amount:     input.Amount,
		Ingredient: model.Ingredient{Name: NormalizeIngredient(input.Name)},
	}
	if err := s.repo.AddIngredient(ctx, ir); err != nil {
		return nil, err
	}
	return &IngredientRecipeResponse{ID: ir.ID, Amount: ir.Amount, Name: ir.Ingredient.Name}, nil
}

func (s *RecipeService) DeleteIngredient(ctx context.Context, recipeID, userID string, ingredientRecipeID uint) error {
	recipe, err := s.findOwned(ctx, recipeID, userID, ErrForbidden)
	if err != nil {
		return err
	}

	ir, err := s.repo.FindIngredientRecipeByID(ctx, ingredientRecipeID)
	if notFound(err) || (err == nil && ir.RecipeID != recipe.ID) {
		return ErrIngredientNotFound
	}
	if err != nil {
		return err
	}
	return s.repo.DeleteIngredientRecipe(ctx, ingredientRecipeID)
}

// --- Preparos ---

func (s *RecipeService) AddPreparation(ctx context.Context, recipeID, userID string, input PreparationInput) (*PreparationResponse, error) {
	recipe, err := s.findOwned(ctx, recipeID, userID, ErrForbidden)
	if err != nil {
		return nil, err
	}
	if len(recipe.Preparations) >= MaxPreparations {
		return nil, ErrTooManyItems
	}

	p := &model.Preparation{Description: input.Description, RecipeID: recipe.ID}
	if err := s.repo.AddPreparation(ctx, p); err != nil {
		return nil, err
	}
	return &PreparationResponse{ID: p.ID, Description: p.Description}, nil
}

func (s *RecipeService) UpdatePreparation(ctx context.Context, recipeID, userID string, prepID uint, input PreparationInput) (*PreparationResponse, error) {
	recipe, err := s.findOwned(ctx, recipeID, userID, ErrForbidden)
	if err != nil {
		return nil, err
	}

	p, err := s.repo.FindPreparationByID(ctx, prepID)
	if notFound(err) || (err == nil && p.RecipeID != recipe.ID) {
		return nil, ErrPreparationNotFound
	}
	if err != nil {
		return nil, err
	}

	p.Description = input.Description
	if err := s.repo.UpdatePreparation(ctx, p); err != nil {
		return nil, err
	}
	return &PreparationResponse{ID: p.ID, Description: p.Description}, nil
}

func (s *RecipeService) DeletePreparation(ctx context.Context, recipeID, userID string, prepID uint) error {
	recipe, err := s.findOwned(ctx, recipeID, userID, ErrForbidden)
	if err != nil {
		return err
	}

	p, err := s.repo.FindPreparationByID(ctx, prepID)
	if notFound(err) || (err == nil && p.RecipeID != recipe.ID) {
		return ErrPreparationNotFound
	}
	if err != nil {
		return err
	}
	return s.repo.DeletePreparation(ctx, prepID)
}

// --- Auxiliares ---

func (s *RecipeService) find(ctx context.Context, id string) (*model.Recipe, error) {
	recipe, err := s.repo.FindByID(ctx, id)
	if err != nil {
		if notFound(err) {
			return nil, ErrRecipeNotFound
		}
		return nil, err
	}
	return recipe, nil
}

// findOwned busca a receita e confere se é de quem pede; forbidden é o erro
// devolvido quando não é.
func (s *RecipeService) findOwned(ctx context.Context, recipeID, userID string, forbidden error) (*model.Recipe, error) {
	recipe, err := s.find(ctx, recipeID)
	if err != nil {
		return nil, err
	}
	if recipe.UserID != userID {
		return nil, forbidden
	}
	return recipe, nil
}

// checkCategory confere a categoria antes de gravar: a chave estrangeira
// violada viraria 500, e a mensagem precisa apontar o campo.
func (s *RecipeService) checkCategory(ctx context.Context, id uint) error {
	ok, err := s.categories.Exists(ctx, id)
	if err != nil {
		return err
	}
	if !ok {
		return ErrCategoryNotFound
	}
	return nil
}

func toRecipeModel(input RecipeInput) *model.Recipe {
	recipe := &model.Recipe{
		Name:        input.Name,
		Description: input.Description,
		Time:        input.Time,
		Portions:    input.Portions,
		CategoryID:  input.CategoryID,
	}
	for _, img := range input.Images {
		recipe.Images = append(recipe.Images, model.ImageRecipe{URL: img.URL})
	}
	for _, ing := range input.Ingredients {
		recipe.Ingredients = append(recipe.Ingredients, model.IngredientRecipe{
			Amount:     ing.Amount,
			Ingredient: model.Ingredient{Name: NormalizeIngredient(ing.Name)},
		})
	}
	for _, prep := range input.Preparations {
		recipe.Preparations = append(recipe.Preparations, model.Preparation{Description: prep.Description})
	}
	return recipe
}

// NormalizeIngredient guarda o ingrediente em minúsculas e sem espaço
// sobrando, para "Farinha" e "farinha " serem o mesmo ingrediente. O seed
// usa a mesma regra.
func NormalizeIngredient(name string) string {
	return strings.ToLower(strings.Join(strings.Fields(name), " "))
}

func notFound(err error) bool { return errors.Is(err, gorm.ErrRecordNotFound) }
