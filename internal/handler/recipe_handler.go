package handler

import (
	"net/http"
	"strconv"

	"panda-cooking-go-api/internal/middleware"
	"panda-cooking-go-api/internal/service"

	"github.com/gin-gonic/gin"
)

type RecipeHandler struct {
	service *service.RecipeService
}

func NewRecipeHandler(service *service.RecipeService) *RecipeHandler {
	return &RecipeHandler{service: service}
}

func (h *RecipeHandler) RegisterRoutes(r *gin.RouterGroup, authMiddleware gin.HandlerFunc) {
	r.GET("", h.getAll)
	r.GET("/:id", h.getByID)
	r.POST("", authMiddleware, h.create)
	r.PATCH("/:id", authMiddleware, h.update)
	r.DELETE("/:id", authMiddleware, h.delete)

	r.POST("/:id/images", authMiddleware, h.addImage)
	r.PATCH("/:id/images/:imageID", authMiddleware, h.updateImage)
	r.DELETE("/:id/images/:imageID", authMiddleware, h.deleteImage)

	r.POST("/:id/ingredients", authMiddleware, h.addIngredient)
	r.DELETE("/:id/ingredients/:ingredientID", authMiddleware, h.deleteIngredient)

	r.POST("/:id/preparations", authMiddleware, h.addPreparation)
	r.PATCH("/:id/preparations/:prepID", authMiddleware, h.updatePreparation)
	r.DELETE("/:id/preparations/:prepID", authMiddleware, h.deletePreparation)
}

func (h *RecipeHandler) getAll(c *gin.Context) {
	recipes, err := h.service.GetAll()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, recipes)
}

func (h *RecipeHandler) getByID(c *gin.Context) {
	recipe, err := h.service.GetByID(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, recipe)
}

func (h *RecipeHandler) create(c *gin.Context) {
	var input service.CreateRecipeInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userID := c.GetString(middleware.UserIDKey)
	recipe, err := h.service.Create(userID, input)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, recipe)
}

func (h *RecipeHandler) update(c *gin.Context) {
	var input service.UpdateRecipeInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userID := c.GetString(middleware.UserIDKey)
	recipe, err := h.service.Update(c.Param("id"), userID, input)
	if err != nil {
		status := http.StatusInternalServerError
		if err.Error() == "sem permissão para editar esta receita" {
			status = http.StatusForbidden
		} else if err.Error() == "receita não encontrada" {
			status = http.StatusNotFound
		}
		c.JSON(status, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, recipe)
}

func (h *RecipeHandler) delete(c *gin.Context) {
	userID := c.GetString(middleware.UserIDKey)
	if err := h.service.Delete(c.Param("id"), userID); err != nil {
		status := http.StatusInternalServerError
		if err.Error() == "sem permissão para deletar esta receita" {
			status = http.StatusForbidden
		} else if err.Error() == "receita não encontrada" {
			status = http.StatusNotFound
		}
		c.JSON(status, gin.H{"error": err.Error()})
		return
	}
	c.Status(http.StatusNoContent)
}

// --- Imagens ---

func (h *RecipeHandler) addImage(c *gin.Context) {
	var input service.ImageRecipeInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userID := c.GetString(middleware.UserIDKey)
	img, err := h.service.AddImage(c.Param("id"), userID, input)
	if err != nil {
		c.JSON(statusFromErr(err), gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, img)
}

func (h *RecipeHandler) updateImage(c *gin.Context) {
	var input service.ImageRecipeInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	imageID, err := parseUintParam(c, "imageID")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "imageID inválido"})
		return
	}

	userID := c.GetString(middleware.UserIDKey)
	img, err := h.service.UpdateImage(c.Param("id"), userID, imageID, input)
	if err != nil {
		c.JSON(statusFromErr(err), gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, img)
}

func (h *RecipeHandler) deleteImage(c *gin.Context) {
	imageID, err := parseUintParam(c, "imageID")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "imageID inválido"})
		return
	}

	userID := c.GetString(middleware.UserIDKey)
	if err := h.service.DeleteImage(c.Param("id"), userID, imageID); err != nil {
		c.JSON(statusFromErr(err), gin.H{"error": err.Error()})
		return
	}

	c.Status(http.StatusNoContent)
}

// --- Ingredientes ---

func (h *RecipeHandler) addIngredient(c *gin.Context) {
	var input service.IngredientRecipeInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userID := c.GetString(middleware.UserIDKey)
	ir, err := h.service.AddIngredient(c.Param("id"), userID, input)
	if err != nil {
		c.JSON(statusFromErr(err), gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, ir)
}

func (h *RecipeHandler) deleteIngredient(c *gin.Context) {
	ingredientID, err := parseUintParam(c, "ingredientID")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ingredientID inválido"})
		return
	}

	userID := c.GetString(middleware.UserIDKey)
	if err := h.service.DeleteIngredient(c.Param("id"), userID, ingredientID); err != nil {
		c.JSON(statusFromErr(err), gin.H{"error": err.Error()})
		return
	}

	c.Status(http.StatusNoContent)
}

// --- Preparos ---

func (h *RecipeHandler) addPreparation(c *gin.Context) {
	var input service.PreparationInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userID := c.GetString(middleware.UserIDKey)
	p, err := h.service.AddPreparation(c.Param("id"), userID, input)
	if err != nil {
		c.JSON(statusFromErr(err), gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, p)
}

func (h *RecipeHandler) updatePreparation(c *gin.Context) {
	var input service.PreparationInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	prepID, err := parseUintParam(c, "prepID")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "prepID inválido"})
		return
	}

	userID := c.GetString(middleware.UserIDKey)
	p, err := h.service.UpdatePreparation(c.Param("id"), userID, prepID, input)
	if err != nil {
		c.JSON(statusFromErr(err), gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, p)
}

func (h *RecipeHandler) deletePreparation(c *gin.Context) {
	prepID, err := parseUintParam(c, "prepID")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "prepID inválido"})
		return
	}

	userID := c.GetString(middleware.UserIDKey)
	if err := h.service.DeletePreparation(c.Param("id"), userID, prepID); err != nil {
		c.JSON(statusFromErr(err), gin.H{"error": err.Error()})
		return
	}

	c.Status(http.StatusNoContent)
}

func parseUintParam(c *gin.Context, name string) (uint, error) {
	v, err := strconv.ParseUint(c.Param(name), 10, 64)
	return uint(v), err
}
