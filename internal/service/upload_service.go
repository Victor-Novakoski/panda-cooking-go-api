package service

import (
	"context"
	"errors"
	"net/http"

	"panda-cooking-go-api/internal/model"
	"panda-cooking-go-api/internal/repository"

	"gorm.io/gorm"
)

// UploadPathPrefix é o começo do endereço de uma foto enviada. O endereço é
// relativo: o navegador pede a foto ao próprio site, que repassa para a API.
const UploadPathPrefix = "/api/uploads/"

// MaxUploadBytes é o tamanho máximo de uma foto. O front reduz a foto antes
// de enviar; o limite do corpo da requisição (1 MB) cobre isso com folga.
const MaxUploadBytes = 900 << 10

// Formatos aceitos, conferidos pelos primeiros bytes do arquivo e não pelo
// nome ou pelo Content-Type que o navegador mandou.
var uploadTypes = map[string]bool{
	"image/jpeg": true,
	"image/png":  true,
	"image/webp": true,
	"image/gif":  true,
}

var (
	ErrUploadNotFound = newError(KindNotFound, "foto não encontrada")
	ErrUploadEmpty    = fieldError(KindInvalid, "file", "envie uma foto")
	ErrUploadTooLarge = fieldError(KindInvalid, "file", "a foto pode ter no máximo 900 KB")
	ErrUploadType     = fieldError(KindInvalid, "file", "use uma foto JPG, PNG, WebP ou GIF")
)

type UploadService struct {
	repo repository.UploadRepo
}

func NewUploadService(repo repository.UploadRepo) *UploadService {
	return &UploadService{repo: repo}
}

// UploadResponse é o endereço da foto, para usar como url da foto da
// receita ou como image_profile.
type UploadResponse struct {
	URL string `json:"url"`
}

// Create guarda a foto de quem está logado.
func (s *UploadService) Create(ctx context.Context, userID string, data []byte) (*UploadResponse, error) {
	if len(data) == 0 {
		return nil, ErrUploadEmpty
	}
	if len(data) > MaxUploadBytes {
		return nil, ErrUploadTooLarge
	}
	contentType := http.DetectContentType(data)
	if !uploadTypes[contentType] {
		return nil, ErrUploadType
	}

	upload := &model.Upload{UserID: userID, ContentType: contentType, Data: data}
	if err := s.repo.Create(ctx, upload); err != nil {
		if errors.Is(err, gorm.ErrForeignKeyViolated) {
			return nil, ErrUserNotFound
		}
		return nil, err
	}
	return &UploadResponse{URL: UploadPathPrefix + upload.ID}, nil
}

func (s *UploadService) Get(ctx context.Context, id string) (*model.Upload, error) {
	upload, err := s.repo.FindByID(ctx, id)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, ErrUploadNotFound
	}
	return upload, err
}
