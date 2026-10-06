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

func NewFavoriteHandler(service *service.FavoriteService) *FavoriteHandler {
	return &FavoriteHandler{service: service}
}

func (h *FavoriteHandler) RegisterRoutes(r *gin.RouterGroup, authMiddleware gin.HandlerFunc) {
	r.POST("/:recipeID", authMiddleware, h.add)
	r.DELETE("/:recipeID", authMiddleware, h.remove)
}

func (h *FavoriteHandler) add(c *gin.Context) {
	userID := c.GetString(middleware.UserIDKey)
	recipeID := c.Param("recipeID")

	fav, err := h.service.Add(userID, recipeID)
	if err != nil {
		c.JSON(statusFromErr(err), gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, fav)
}

func (h *FavoriteHandler) remove(c *gin.Context) {
	userID := c.GetString(middleware.UserIDKey)
	recipeID := c.Param("recipeID")

	if err := h.service.Remove(userID, recipeID); err != nil {
		c.JSON(statusFromErr(err), gin.H{"error": err.Error()})
		return
	}

	c.Status(http.StatusNoContent)
}
