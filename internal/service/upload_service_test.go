package service_test

import (
	"bytes"
	"context"
	"testing"

	"panda-cooking-go-api/internal/model"
	"panda-cooking-go-api/internal/service"
	"panda-cooking-go-api/internal/service/mocks"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gorm.io/gorm"
)

// Começo de um PNG: basta para o formato ser reconhecido.
var pngHeader = []byte("\x89PNG\r\n\x1a\n\x00\x00\x00\rIHDR")

func TestUploadService_Create(t *testing.T) {
	var saved *model.Upload
	repo := &mocks.UploadRepoMock{CreateFn: func(_ context.Context, u *model.Upload) error {
		u.ID = "11111111-1111-1111-1111-111111111111"
		saved = u
		return nil
	}}

	resp, err := service.NewUploadService(repo).Create(ctx, userID, pngHeader)
	require.NoError(t, err)
	assert.Equal(t, "/api/uploads/11111111-1111-1111-1111-111111111111", resp.URL)
	assert.Equal(t, userID, saved.UserID)
	assert.Equal(t, "image/png", saved.ContentType)
}

func TestUploadService_CreateRecusa(t *testing.T) {
	tests := []struct {
		name    string
		data    []byte
		wantErr error
	}{
		{"vazio", nil, service.ErrUploadEmpty},
		{"grande demais", append(bytes.Clone(pngHeader), make([]byte, service.MaxUploadBytes)...), service.ErrUploadTooLarge},
		{"não é foto", []byte("<svg xmlns='http://www.w3.org/2000/svg'><script>alert(1)</script></svg>"), service.ErrUploadType},
		{"texto com nome de foto", []byte("só texto"), service.ErrUploadType},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create no repositório não pode ser chamado: o mock entra em pânico.
			_, err := service.NewUploadService(&mocks.UploadRepoMock{}).Create(ctx, userID, tt.data)
			assert.ErrorIs(t, err, tt.wantErr)
		})
	}
}

func TestUploadService_CreateContaApagada(t *testing.T) {
	repo := &mocks.UploadRepoMock{CreateFn: func(context.Context, *model.Upload) error { return gorm.ErrForeignKeyViolated }}
	_, err := service.NewUploadService(repo).Create(ctx, userID, pngHeader)
	assert.ErrorIs(t, err, service.ErrUserNotFound)
}

func TestUploadService_Get(t *testing.T) {
	repo := &mocks.UploadRepoMock{FindByIDFn: func(_ context.Context, id string) (*model.Upload, error) {
		if id == "existe" {
			return &model.Upload{ID: id, ContentType: "image/png", Data: pngHeader}, nil
		}
		return nil, gorm.ErrRecordNotFound
	}}
	svc := service.NewUploadService(repo)

	upload, err := svc.Get(ctx, "existe")
	require.NoError(t, err)
	assert.Equal(t, "image/png", upload.ContentType)

	_, err = svc.Get(ctx, "outra")
	assert.ErrorIs(t, err, service.ErrUploadNotFound)
}
