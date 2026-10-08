package repository

import (
	"context"

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

func (r *CommentRepository) Create(ctx context.Context, comment *model.Comment) error {
	return r.db.WithContext(ctx).Omit(clause.Associations).Create(comment).Error
}

// ListByRecipe lista os comentários da receita, dos mais novos aos mais antigos.
func (r *CommentRepository) ListByRecipe(ctx context.Context, recipeID string, page Page) ([]model.Comment, int64, error) {
	if !IsUUID(recipeID) {
		return nil, 0, nil
	}
	// Session deixa reaproveitar a consulta no Count e no Find
	db := r.db.WithContext(ctx).Model(&model.Comment{}).Where("recipe_id = ?", recipeID).Session(&gorm.Session{})

	var total int64
	if err := db.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	var comments []model.Comment
	err := db.Preload("User", selectAuthor).
		Order("created_at DESC, id DESC").
		Limit(page.Size).Offset(page.offset()).
		Find(&comments).Error
	return comments, total, err
}

func (r *CommentRepository) FindByID(ctx context.Context, id uint) (*model.Comment, error) {
	var comment model.Comment
	err := r.db.WithContext(ctx).Preload("User", selectAuthor).First(&comment, id).Error
	if err != nil {
		return nil, err
	}
	return &comment, nil
}

func (r *CommentRepository) Update(ctx context.Context, comment *model.Comment) error {
	return r.db.WithContext(ctx).Model(comment).Select("Description", "UpdatedAt").Updates(comment).Error
}

func (r *CommentRepository) Delete(ctx context.Context, id uint) error {
	return r.db.WithContext(ctx).Delete(&model.Comment{}, id).Error
}

// selectAuthor carrega do autor só o que é público: nada de e-mail ou hash.
func selectAuthor(db *gorm.DB) *gorm.DB { return db.Select("id", "name", "image_profile") }
