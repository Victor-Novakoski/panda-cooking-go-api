package handler

import (
	"net/http"

	"panda-cooking-go-api/internal/service"

	"github.com/gin-gonic/gin"
)

type CategoryHandler struct {
	service *service.CategoryService
}

func NewCategoryHandler(service *service.CategoryService) *CategoryHandler {
	return &CategoryHandler{service: service}
}

func (h *CategoryHandler) RegisterRoutes(r *gin.RouterGroup) {
	r.GET("", h.getAll)
}

func (h *CategoryHandler) getAll(c *gin.Context) {
	categories, err := h.service.GetAll()
	if err != nil {
		respondError(c, err)
		return
	}
	c.JSON(http.StatusOK, categories)
}
