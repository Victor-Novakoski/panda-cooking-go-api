package service

import (
	"errors"

	"panda-cooking-go-api/internal/model"
	"panda-cooking-go-api/internal/repository"

	"gorm.io/gorm"
)

type CommentService struct {
	repo       repository.CommentRepo
	recipeRepo repository.RecipeRepo
}

func NewCommentService(repo repository.CommentRepo, recipeRepo repository.RecipeRepo) *CommentService {
	return &CommentService{repo: repo, recipeRepo: recipeRepo}
}

// --- DTOs ---

type CreateCommentInput struct {
	Description string `json:"description" binding:"required"`
	RecipeID    string `json:"recipe_id" binding:"required"`
}

type UpdateCommentInput struct {
	Description string `json:"description" binding:"required"`
}

type CommentResponse struct {
	ID          uint         `json:"id"`
	Description string       `json:"description"`
	RecipeID    string       `json:"recipe_id"`
	User        UserResponse `json:"user"`
}

// --- Métodos ---

func (s *CommentService) Create(userID string, input CreateCommentInput) (*CommentResponse, error) {
	// garante que a receita existe antes de comentar
	if _, err := s.recipeRepo.FindByID(input.RecipeID); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrRecipeNotFound
		}
		return nil, err
	}

	comment := &model.Comment{
		Description: input.Description,
		RecipeID:    input.RecipeID,
		UserID:      userID,
	}

	if err := s.repo.Create(comment); err != nil {
		return nil, err
	}

	// recarrega com o User para montar a resposta
	created, err := s.repo.FindByID(comment.ID)
	if err != nil {
		return nil, err
	}

	return toCommentResponse(*created), nil
}

func (s *CommentService) GetAll() ([]CommentResponse, error) {
	comments, err := s.repo.FindAll()
	if err != nil {
		return nil, err
	}

	result := make([]CommentResponse, len(comments))
	for i, c := range comments {
		result[i] = *toCommentResponse(c)
	}
	return result, nil
}

func (s *CommentService) Update(commentID uint, userID string, input UpdateCommentInput) (*CommentResponse, error) {
	comment, err := s.repo.FindByID(commentID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrCommentNotFound
		}
		return nil, err
	}

	if comment.UserID != userID {
		return nil, ErrForbiddenEditComment
	}

	comment.Description = input.Description
	if err := s.repo.Update(comment); err != nil {
		return nil, err
	}

	updated, err := s.repo.FindByID(comment.ID)
	if err != nil {
		return nil, err
	}

	return toCommentResponse(*updated), nil
}

func (s *CommentService) Delete(commentID uint, userID string, isAdm bool) error {
	comment, err := s.repo.FindByID(commentID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return ErrCommentNotFound
		}
		return err
	}

	// admin pode deletar qualquer comentário, usuário só o próprio
	if !isAdm && comment.UserID != userID {
		return ErrForbiddenDelComment
	}

	return s.repo.Delete(commentID)
}

func toCommentResponse(c model.Comment) *CommentResponse {
	return &CommentResponse{
		ID:          c.ID,
		Description: c.Description,
		RecipeID:    c.RecipeID,
		User:        *toUserResponse(&c.User),
	}
}
