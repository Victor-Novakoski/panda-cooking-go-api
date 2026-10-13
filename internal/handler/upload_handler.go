package handler

import (
	"errors"
	"io"
	"net/http"

	"panda-cooking-go-api/internal/middleware"
	"panda-cooking-go-api/internal/service"

	"github.com/gin-gonic/gin"
)

type UploadHandler struct {
	service *service.UploadService
}

func NewUploadHandler(svc *service.UploadService) *UploadHandler {
	return &UploadHandler{service: svc}
}

// Create recebe uma foto (multipart, campo "file") e devolve o endereço dela.
func (h *UploadHandler) Create(c *gin.Context) {
	file, _, err := c.Request.FormFile("file")
	if err != nil {
		var maxErr *http.MaxBytesError
		if errors.As(err, &maxErr) {
			c.AbortWithStatusJSON(http.StatusRequestEntityTooLarge, ErrorBody{Error: "corpo da requisição grande demais"})
			return
		}
		respondError(c, service.ErrUploadEmpty)
		return
	}
	defer file.Close()

	// Lê um byte além do limite só para saber que passou dele.
	data, err := io.ReadAll(io.LimitReader(file, service.MaxUploadBytes+1))
	if err != nil {
		respondError(c, err)
		return
	}

	resp, err := h.service.Create(c.Request.Context(), c.GetString(middleware.UserIDKey), data)
	if err != nil {
		respondError(c, err)
		return
	}
	c.JSON(http.StatusCreated, resp)
}

// Get serve a foto. O endereço nunca muda de conteúdo, então o navegador
// pode guardar a foto por um ano. O CSP sandbox impede que um arquivo
// aberto direto na aba rode qualquer coisa.
func (h *UploadHandler) Get(c *gin.Context) {
	upload, err := h.service.Get(c.Request.Context(), c.Param("id"))
	if err != nil {
		respondError(c, err)
		return
	}
	c.Header("Cache-Control", "public, max-age=31536000, immutable")
	c.Header("Content-Security-Policy", "default-src 'none'; sandbox")
	c.Header("X-Content-Type-Options", "nosniff")
	c.Data(http.StatusOK, upload.ContentType, upload.Data)
}
