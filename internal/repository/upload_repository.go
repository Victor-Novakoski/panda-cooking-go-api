package repository

import (
	"context"

	"panda-cooking-go-api/internal/model"

	"gorm.io/gorm"
)

type UploadRepository struct {
	db *gorm.DB
}

func NewUploadRepository(db *gorm.DB) *UploadRepository {
	return &UploadRepository{db: db}
}

func (r *UploadRepository) Create(ctx context.Context, upload *model.Upload) error {
	return r.db.WithContext(ctx).Create(upload).Error
}

func (r *UploadRepository) FindByID(ctx context.Context, id string) (*model.Upload, error) {
	if !IsUUID(id) {
		return nil, gorm.ErrRecordNotFound
	}
	var upload model.Upload
	err := r.db.WithContext(ctx).First(&upload, "id = ?", id).Error
	if err != nil {
		return nil, err
	}
	return &upload, nil
}
