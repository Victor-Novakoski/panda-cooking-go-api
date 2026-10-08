package handler

import (
	"net/http"

	"panda-cooking-go-api/internal/middleware"
	"panda-cooking-go-api/internal/service"

	"github.com/gin-gonic/gin"
)

type RecipeHandler struct {
	service *service.RecipeService
}

func NewRecipeHandler(svc *service.RecipeService) *RecipeHandler {
	return &RecipeHandler{service: svc}
}

// List é a listagem pública: busca por nome e descrição, filtro por
// categoria e por autor, paginada.
func (h *RecipeHandler) List(c *gin.Context) {
	var q service.ListRecipesQuery
	if !bindQuery(c, &q) {
		return
	}

	page, err := h.service.List(c.Request.Context(), q)
	if err != nil {
		respondError(c, err)
		return
	}
	c.JSON(http.StatusOK, page)
}

func (h *RecipeHandler) Get(c *gin.Context) {
	recipe, err := h.service.GetByID(c.Request.Context(), c.Param("id"))
	if err != nil {
		respondError(c, err)
		return
	}
	c.JSON(http.StatusOK, recipe)
}

func (h *RecipeHandler) Create(c *gin.Context) {
	var input service.RecipeInput
	if !bindJSON(c, &input) {
		return
	}

	recipe, err := h.service.Create(c.Request.Context(), c.GetString(middleware.UserIDKey), input)
	if err != nil {
		respondError(c, err)
		return
	}
	c.JSON(http.StatusCreated, recipe)
}

// Update (PATCH) muda só os campos enviados.
func (h *RecipeHandler) Update(c *gin.Context) {
	var input service.UpdateRecipeInput
	if !bindJSON(c, &input) {
		return
	}

	recipe, err := h.service.Update(c.Request.Context(), c.Param("id"), c.GetString(middleware.UserIDKey), input)
	if err != nil {
		respondError(c, err)
		return
	}
	c.JSON(http.StatusOK, recipe)
}

// Replace (PUT) troca a receita inteira numa transação.
func (h *RecipeHandler) Replace(c *gin.Context) {
	var input service.RecipeInput
	if !bindJSON(c, &input) {
		return
	}

	recipe, err := h.service.Replace(c.Request.Context(), c.Param("id"), c.GetString(middleware.UserIDKey), input)
	if err != nil {
		respondError(c, err)
		return
	}
	c.JSON(http.StatusOK, recipe)
}

func (h *RecipeHandler) Delete(c *gin.Context) {
	if err := h.service.Delete(c.Request.Context(), c.Param("id"), c.GetString(middleware.UserIDKey)); err != nil {
		respondError(c, err)
		return
	}
	c.Status(http.StatusNoContent)
}

// --- Imagens ---

func (h *RecipeHandler) AddImage(c *gin.Context) {
	var input service.ImageRecipeInput
	if !bindJSON(c, &input) {
		return
	}

	img, err := h.service.AddImage(c.Request.Context(), c.Param("id"), c.GetString(middleware.UserIDKey), input)
	if err != nil {
		respondError(c, err)
		return
	}
	c.JSON(http.StatusCreated, img)
}

func (h *RecipeHandler) UpdateImage(c *gin.Context) {
	imageID, ok := idParam(c, "imageID", service.ErrImageNotFound)
	if !ok {
		return
	}
	var input service.ImageRecipeInput
	if !bindJSON(c, &input) {
		return
	}

	img, err := h.service.UpdateImage(c.Request.Context(), c.Param("id"), c.GetString(middleware.UserIDKey), imageID, input)
	if err != nil {
		respondError(c, err)
		return
	}
	c.JSON(http.StatusOK, img)
}

func (h *RecipeHandler) DeleteImage(c *gin.Context) {
	imageID, ok := idParam(c, "imageID", service.ErrImageNotFound)
	if !ok {
		return
	}

	if err := h.service.DeleteImage(c.Request.Context(), c.Param("id"), c.GetString(middleware.UserIDKey), imageID); err != nil {
		respondError(c, err)
		return
	}
	c.Status(http.StatusNoContent)
}

// --- Ingredientes ---

func (h *RecipeHandler) AddIngredient(c *gin.Context) {
	var input service.IngredientRecipeInput
	if !bindJSON(c, &input) {
		return
	}

	ir, err := h.service.AddIngredient(c.Request.Context(), c.Param("id"), c.GetString(middleware.UserIDKey), input)
	if err != nil {
		respondError(c, err)
		return
	}
	c.JSON(http.StatusCreated, ir)
}

func (h *RecipeHandler) DeleteIngredient(c *gin.Context) {
	ingredientID, ok := idParam(c, "ingredientID", service.ErrIngredientNotFound)
	if !ok {
		return
	}

	if err := h.service.DeleteIngredient(c.Request.Context(), c.Param("id"), c.GetString(middleware.UserIDKey), ingredientID); err != nil {
		respondError(c, err)
		return
	}
	c.Status(http.StatusNoContent)
}

// --- Preparos ---

func (h *RecipeHandler) AddPreparation(c *gin.Context) {
	var input service.PreparationInput
	if !bindJSON(c, &input) {
		return
	}

	p, err := h.service.AddPreparation(c.Request.Context(), c.Param("id"), c.GetString(middleware.UserIDKey), input)
	if err != nil {
		respondError(c, err)
		return
	}
	c.JSON(http.StatusCreated, p)
}

func (h *RecipeHandler) UpdatePreparation(c *gin.Context) {
	prepID, ok := idParam(c, "prepID", service.ErrPreparationNotFound)
	if !ok {
		return
	}
	var input service.PreparationInput
	if !bindJSON(c, &input) {
		return
	}

	p, err := h.service.UpdatePreparation(c.Request.Context(), c.Param("id"), c.GetString(middleware.UserIDKey), prepID, input)
	if err != nil {
		respondError(c, err)
		return
	}
	c.JSON(http.StatusOK, p)
}

func (h *RecipeHandler) DeletePreparation(c *gin.Context) {
	prepID, ok := idParam(c, "prepID", service.ErrPreparationNotFound)
	if !ok {
		return
	}

	if err := h.service.DeletePreparation(c.Request.Context(), c.Param("id"), c.GetString(middleware.UserIDKey), prepID); err != nil {
		respondError(c, err)
		return
	}
	c.Status(http.StatusNoContent)
}
