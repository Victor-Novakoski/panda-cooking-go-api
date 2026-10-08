package handler

import (
	"net/http"

	"panda-cooking-go-api/internal/middleware"
	"panda-cooking-go-api/internal/service"

	"github.com/gin-gonic/gin"
)

type FavoriteHandler struct {
	service *service.FavoriteService
}

func NewFavoriteHandler(svc *service.FavoriteService) *FavoriteHandler {
	return &FavoriteHandler{service: svc}
}

// Status diz se a receita está nos favoritos de quem está logado.
func (h *FavoriteHandler) Status(c *gin.Context) {
	status, err := h.service.Status(c.Request.Context(), c.GetString(middleware.UserIDKey), c.Param("recipeID"))
	if err != nil {
		respondError(c, err)
		return
	}
	c.JSON(http.StatusOK, status)
}

func (h *FavoriteHandler) Add(c *gin.Context) {
	status, err := h.service.Add(c.Request.Context(), c.GetString(middleware.UserIDKey), c.Param("recipeID"))
	if err != nil {
		respondError(c, err)
		return
	}
	c.JSON(http.StatusCreated, status)
}

func (h *FavoriteHandler) Remove(c *gin.Context) {
	if err := h.service.Remove(c.Request.Context(), c.GetString(middleware.UserIDKey), c.Param("recipeID")); err != nil {
		respondError(c, err)
		return
	}
	c.Status(http.StatusNoContent)
}
