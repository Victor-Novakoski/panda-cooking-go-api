package service

import "strings"

// Entradas da API. As tags binding são conferidas no handler (erro 422 com
// a mensagem de cada campo); Normalize roda antes, para a validação já ver
// o texto sem espaço nas pontas e o e-mail em minúsculas.
//
// Limites de tamanho: os mesmos do front (src/lib/limits.ts).

type CreateUserInput struct {
	Name         string `json:"name" binding:"required,min=2,max=80"`
	Email        string `json:"email" binding:"required,max=254,email"`
	Password     string `json:"password" binding:"required,min=10,max=128"`
	ImageProfile string `json:"image_profile" binding:"omitempty,max=2048,httpsurl"`
}

func (in *CreateUserInput) Normalize() {
	in.Name = strings.TrimSpace(in.Name)
	in.Email = normalizeEmail(in.Email)
	in.ImageProfile = strings.TrimSpace(in.ImageProfile)
}

// UpdateUserInput usa ponteiro para separar "não mandou" (nil, fica como
// está) de "mandou vazio": foto vazia remove a foto; nome vazio é recusado.
type UpdateUserInput struct {
	Name         *string `json:"name" binding:"omitempty,min=2,max=80"`
	ImageProfile *string `json:"image_profile" binding:"omitempty,max=2048,httpsurl"`
}

func (in *UpdateUserInput) Normalize() {
	trimPtr(in.Name)
	trimPtr(in.ImageProfile)
}

type LoginInput struct {
	Email    string `json:"email" binding:"required,max=254,email"`
	Password string `json:"password" binding:"required,max=128"`
}

func (in *LoginInput) Normalize() { in.Email = normalizeEmail(in.Email) }

// RecipeInput é a receita inteira: criação (POST) e troca completa (PUT).
// No PUT, foto, ingrediente ou passo que não vier deixa de existir.
type RecipeInput struct {
	Name         string                  `json:"name" binding:"required,min=3,max=120"`
	Description  string                  `json:"description" binding:"required,min=10,max=2000"`
	Time         string                  `json:"time" binding:"required,max=50"`
	Portions     int                     `json:"portions" binding:"required,min=1,max=100"`
	CategoryID   uint                    `json:"category_id" binding:"required,max=9223372036854775807"`
	Images       []ImageRecipeInput      `json:"images" binding:"max=10,dive"`
	Ingredients  []IngredientRecipeInput `json:"ingredients" binding:"required,min=1,max=50,dive"`
	Preparations []PreparationInput      `json:"preparations" binding:"required,min=1,max=50,dive"`
}

func (in *RecipeInput) Normalize() {
	in.Name = strings.TrimSpace(in.Name)
	in.Description = strings.TrimSpace(in.Description)
	in.Time = strings.TrimSpace(in.Time)
	for i := range in.Images {
		in.Images[i].Normalize()
	}
	for i := range in.Ingredients {
		in.Ingredients[i].Normalize()
	}
	for i := range in.Preparations {
		in.Preparations[i].Normalize()
	}
}

// UpdateRecipeInput (PATCH) muda só os campos enviados.
type UpdateRecipeInput struct {
	Name        *string `json:"name" binding:"omitempty,min=3,max=120"`
	Description *string `json:"description" binding:"omitempty,min=10,max=2000"`
	Time        *string `json:"time" binding:"omitempty,min=1,max=50"`
	Portions    *int    `json:"portions" binding:"omitempty,min=1,max=100"`
	CategoryID  *uint   `json:"category_id" binding:"omitempty,min=1,max=9223372036854775807"`
}

func (in *UpdateRecipeInput) Normalize() {
	trimPtr(in.Name)
	trimPtr(in.Description)
	trimPtr(in.Time)
}

type ImageRecipeInput struct {
	URL string `json:"url" binding:"required,max=2048,httpsurl"`
}

func (in *ImageRecipeInput) Normalize() { in.URL = strings.TrimSpace(in.URL) }

type IngredientRecipeInput struct {
	Name   string `json:"name" binding:"required,max=80"`
	Amount string `json:"amount" binding:"required,max=60"`
}

func (in *IngredientRecipeInput) Normalize() {
	in.Name = strings.TrimSpace(in.Name)
	in.Amount = strings.TrimSpace(in.Amount)
}

type PreparationInput struct {
	Description string `json:"description" binding:"required,max=1000"`
}

func (in *PreparationInput) Normalize() { in.Description = strings.TrimSpace(in.Description) }

type CommentInput struct {
	Description string `json:"description" binding:"required,max=1000"`
}

func (in *CommentInput) Normalize() { in.Description = strings.TrimSpace(in.Description) }

// ListRecipesQuery são os filtros da listagem de receitas (query string).
type ListRecipesQuery struct {
	Search     string `form:"search" binding:"max=100"`
	CategoryID uint   `form:"category_id" binding:"max=9223372036854775807"`
	UserID     string `form:"user_id" binding:"omitempty,uuid"`
	Page       int    `form:"page" binding:"omitempty,min=1,max=10000"`
	PerPage    int    `form:"per_page" binding:"omitempty,min=1,max=50"`
}

func (in *ListRecipesQuery) Normalize() {
	in.Search = strings.Join(strings.Fields(in.Search), " ")
	in.UserID = strings.TrimSpace(in.UserID)
}

func normalizeEmail(email string) string { return strings.ToLower(strings.TrimSpace(email)) }

func trimPtr(s *string) {
	if s != nil {
		*s = strings.TrimSpace(*s)
	}
}
