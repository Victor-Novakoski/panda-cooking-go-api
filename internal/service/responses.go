package service

import (
	"time"

	"panda-cooking-go-api/internal/model"
)

// Respostas da API. Os services nunca devolvem o model direto, para não
// vazar campo que não deve sair (hash da senha, e-mail de outra pessoa).

type UserResponse struct {
	ID           string    `json:"id"`
	Name         string    `json:"name"`
	Email        string    `json:"email"`
	ImageProfile string    `json:"image_profile"`
	IsAdm        bool      `json:"is_adm"`
	CreatedAt    time.Time `json:"created_at"`
}

// AuthorResponse é o autor de receita ou comentário: só o que é público.
type AuthorResponse struct {
	ID           string `json:"id"`
	Name         string `json:"name"`
	ImageProfile string `json:"image_profile"`
}

type CategoryResponse struct {
	ID   uint   `json:"id"`
	Name string `json:"name"`
}

// RecipeSummary é a receita na listagem: o que o card mostra.
type RecipeSummary struct {
	ID          string           `json:"id"`
	Name        string           `json:"name"`
	Description string           `json:"description"`
	Time        string           `json:"time"`
	Portions    int              `json:"portions"`
	ImageURL    string           `json:"image_url"`
	Category    CategoryResponse `json:"category"`
	Author      AuthorResponse   `json:"author"`
	CreatedAt   time.Time        `json:"created_at"`
}

type RecipeResponse struct {
	ID           string                     `json:"id"`
	Name         string                     `json:"name"`
	Description  string                     `json:"description"`
	Time         string                     `json:"time"`
	Portions     int                        `json:"portions"`
	Category     CategoryResponse           `json:"category"`
	Author       AuthorResponse             `json:"author"`
	Images       []ImageRecipeResponse      `json:"images"`
	Ingredients  []IngredientRecipeResponse `json:"ingredients"`
	Preparations []PreparationResponse      `json:"preparations"`
	CreatedAt    time.Time                  `json:"created_at"`
	UpdatedAt    time.Time                  `json:"updated_at"`
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

type CommentResponse struct {
	ID          uint           `json:"id"`
	Description string         `json:"description"`
	RecipeID    string         `json:"recipe_id"`
	User        AuthorResponse `json:"user"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
}

func toUserResponse(u *model.User) UserResponse {
	return UserResponse{
		ID:           u.ID,
		Name:         u.Name,
		Email:        u.Email,
		ImageProfile: u.ImageProfile,
		IsAdm:        u.IsAdm,
		CreatedAt:    u.CreatedAt,
	}
}

func toAuthor(u model.User) AuthorResponse {
	return AuthorResponse{ID: u.ID, Name: u.Name, ImageProfile: u.ImageProfile}
}

func toRecipeSummary(r model.Recipe) RecipeSummary {
	s := RecipeSummary{
		ID:          r.ID,
		Name:        r.Name,
		Description: r.Description,
		Time:        r.Time,
		Portions:    r.Portions,
		Category:    CategoryResponse{ID: r.Category.ID, Name: r.Category.Name},
		Author:      toAuthor(r.User),
		CreatedAt:   r.CreatedAt,
	}
	if len(r.Images) > 0 {
		s.ImageURL = r.Images[0].URL
	}
	return s
}

func toRecipeSummaries(recipes []model.Recipe) []RecipeSummary {
	result := make([]RecipeSummary, len(recipes))
	for i, r := range recipes {
		result[i] = toRecipeSummary(r)
	}
	return result
}

// toRecipeResponse monta a receita completa. Listas vazias saem como [] no
// JSON, nunca null, para o front não precisar tratar os dois casos.
func toRecipeResponse(r model.Recipe) RecipeResponse {
	resp := RecipeResponse{
		ID:           r.ID,
		Name:         r.Name,
		Description:  r.Description,
		Time:         r.Time,
		Portions:     r.Portions,
		Category:     CategoryResponse{ID: r.Category.ID, Name: r.Category.Name},
		Author:       toAuthor(r.User),
		Images:       make([]ImageRecipeResponse, 0, len(r.Images)),
		Ingredients:  make([]IngredientRecipeResponse, 0, len(r.Ingredients)),
		Preparations: make([]PreparationResponse, 0, len(r.Preparations)),
		CreatedAt:    r.CreatedAt,
		UpdatedAt:    r.UpdatedAt,
	}
	for _, img := range r.Images {
		resp.Images = append(resp.Images, ImageRecipeResponse{ID: img.ID, URL: img.URL})
	}
	for _, ir := range r.Ingredients {
		resp.Ingredients = append(resp.Ingredients, IngredientRecipeResponse{ID: ir.ID, Amount: ir.Amount, Name: ir.Ingredient.Name})
	}
	for _, p := range r.Preparations {
		resp.Preparations = append(resp.Preparations, PreparationResponse{ID: p.ID, Description: p.Description})
	}
	return resp
}

func toCommentResponse(c model.Comment) CommentResponse {
	return CommentResponse{
		ID:          c.ID,
		Description: c.Description,
		RecipeID:    c.RecipeID,
		User:        toAuthor(c.User),
		CreatedAt:   c.CreatedAt,
		UpdatedAt:   c.UpdatedAt,
	}
}
