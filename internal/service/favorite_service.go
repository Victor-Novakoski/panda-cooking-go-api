package service

import (
	"errors"

	"panda-cooking-go-api/internal/model"
	"panda-cooking-go-api/internal/repository"

	"gorm.io/gorm"
)

type FavoriteService struct {
	repo       repository.FavoriteRepo
	recipeRepo repository.RecipeRepo
}

func NewFavoriteService(repo repository.FavoriteRepo, recipeRepo repository.RecipeRepo) *FavoriteService {
	return &FavoriteService{repo: repo, recipeRepo: recipeRepo}
}

type FavoriteResponse struct {
	ID     uint           `json:"id"`
	Recipe RecipeResponse `json:"recipe"`
}

func (s *FavoriteService) Add(userID, recipeID string) (*FavoriteResponse, error) {
	if _, err := s.recipeRepo.FindByID(recipeID); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("receita não encontrada")
		}
		return nil, err
	}

	// impede favorito duplicado
	if _, err := s.repo.Find(userID, recipeID); err == nil {
		return nil, errors.New("receita já está nos favoritos")
	}

	fav := &model.FavoriteRecipe{UserID: userID, RecipeID: recipeID}
	if err := s.repo.Create(fav); err != nil {
		return nil, err
	}

	recipe, err := s.recipeRepo.FindByID(recipeID)
	if err != nil {
		return nil, err
	}

	resp := toRecipeResponse(*recipe)
	return &FavoriteResponse{ID: fav.ID, Recipe: resp}, nil
}

func (s *FavoriteService) Remove(userID, recipeID string) error {
	fav, err := s.repo.Find(userID, recipeID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("receita não está nos favoritos")
		}
		return err
	}

	return s.repo.Delete(fav.ID)
}
