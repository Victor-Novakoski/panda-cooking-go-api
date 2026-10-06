package service

import (
	"panda-cooking-go-api/internal/repository"
)

type CategoryService struct {
	repo repository.CategoryRepo
}

func NewCategoryService(repo repository.CategoryRepo) *CategoryService {
	return &CategoryService{repo: repo}
}

func (s *CategoryService) GetAll() ([]CategoryResponse, error) {
	categories, err := s.repo.FindAll()
	if err != nil {
		return nil, err
	}

	result := make([]CategoryResponse, len(categories))
	for i, c := range categories {
		result[i] = CategoryResponse{ID: c.ID, Name: c.Name}
	}
	return result, nil
}
