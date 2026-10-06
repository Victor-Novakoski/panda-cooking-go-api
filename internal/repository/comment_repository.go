package repository

import (
	"panda-cooking-go-api/internal/model"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type CommentRepository struct {
	db *gorm.DB
}

func NewCommentRepository(db *gorm.DB) *CommentRepository {
	return &CommentRepository{db: db}
}

func (r *CommentRepository) Create(comment *model.Comment) error {
	return r.db.Create(comment).Error
}

func (r *CommentRepository) FindAll() ([]model.Comment, error) {
	var comments []model.Comment
	err := r.db.Preload("User").Find(&comments).Error
	return comments, err
}

// FindByRecipe lista os comentários da receita, do mais antigo ao mais novo.
func (r *CommentRepository) FindByRecipe(recipeID string) ([]model.Comment, error) {
	var comments []model.Comment
	err := r.db.Preload("User").
		Where("recipe_id = ?", recipeID).
		Order("created_at, id").
		Find(&comments).Error
	return comments, err
}

func (r *CommentRepository) FindByID(id uint) (*model.Comment, error) {
	var comment model.Comment
	err := r.db.Preload("User").First(&comment, id).Error
	if err != nil {
		return nil, err
	}
	return &comment, nil
}

func (r *CommentRepository) Update(comment *model.Comment) error {
	return r.db.Omit(clause.Associations).Save(comment).Error
}

func (r *CommentRepository) Delete(id uint) error {
	return r.db.Delete(&model.Comment{}, id).Error
}
