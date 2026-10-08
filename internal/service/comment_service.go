package service

import (
	"context"
	"errors"

	"panda-cooking-go-api/internal/model"
	"panda-cooking-go-api/internal/repository"

	"gorm.io/gorm"
)

type CommentService struct {
	repo    repository.CommentRepo
	recipes repository.RecipeRepo
}

func NewCommentService(repo repository.CommentRepo, recipes repository.RecipeRepo) *CommentService {
	return &CommentService{repo: repo, recipes: recipes}
}

func (s *CommentService) Create(ctx context.Context, userID, recipeID string, input CommentInput) (*CommentResponse, error) {
	if err := s.checkRecipe(ctx, recipeID); err != nil {
		return nil, err
	}

	comment := &model.Comment{Description: input.Description, RecipeID: recipeID, UserID: userID}
	if err := s.repo.Create(ctx, comment); err != nil {
		// receita apagada entre a conferência e a gravação, ou conta apagada
		// com o access token ainda válido
		if errors.Is(err, gorm.ErrForeignKeyViolated) {
			return nil, missingRecipeOrUser(s.checkRecipe(ctx, recipeID))
		}
		return nil, err
	}

	// recarrega com o autor para montar a resposta
	created, err := s.repo.FindByID(ctx, comment.ID)
	if err != nil {
		return nil, err
	}
	resp := toCommentResponse(*created)
	return &resp, nil
}

// ListByRecipe lista os comentários da receita, dos mais novos aos mais antigos.
func (s *CommentService) ListByRecipe(ctx context.Context, recipeID string, q PageQuery) (Page[CommentResponse], error) {
	if err := s.checkRecipe(ctx, recipeID); err != nil {
		return Page[CommentResponse]{}, err
	}

	page := q.repo(DefaultCommentsPerPage)
	comments, total, err := s.repo.ListByRecipe(ctx, recipeID, page)
	if err != nil {
		return Page[CommentResponse]{}, err
	}

	items := make([]CommentResponse, len(comments))
	for i, c := range comments {
		items[i] = toCommentResponse(c)
	}
	return newPage(items, page, total), nil
}

func (s *CommentService) Update(ctx context.Context, commentID uint, userID string, input CommentInput) (*CommentResponse, error) {
	comment, err := s.find(ctx, commentID)
	if err != nil {
		return nil, err
	}
	if comment.UserID != userID {
		return nil, ErrForbiddenEditComment
	}

	comment.Description = input.Description
	if err := s.repo.Update(ctx, comment); err != nil {
		return nil, err
	}

	updated, err := s.repo.FindByID(ctx, comment.ID)
	if err != nil {
		return nil, err
	}
	resp := toCommentResponse(*updated)
	return &resp, nil
}

// Delete: o autor apaga o próprio comentário; o admin apaga qualquer um.
func (s *CommentService) Delete(ctx context.Context, commentID uint, userID string, isAdm bool) error {
	comment, err := s.find(ctx, commentID)
	if err != nil {
		return err
	}
	if !isAdm && comment.UserID != userID {
		return ErrForbiddenDelComment
	}
	return s.repo.Delete(ctx, commentID)
}

func (s *CommentService) find(ctx context.Context, id uint) (*model.Comment, error) {
	comment, err := s.repo.FindByID(ctx, id)
	if err != nil {
		if notFound(err) {
			return nil, ErrCommentNotFound
		}
		return nil, err
	}
	return comment, nil
}

func (s *CommentService) checkRecipe(ctx context.Context, recipeID string) error {
	ok, err := s.recipes.Exists(ctx, recipeID)
	if err != nil {
		return err
	}
	if !ok {
		return ErrRecipeNotFound
	}
	return nil
}
