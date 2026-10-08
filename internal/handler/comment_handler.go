package handler

import (
	"net/http"

	"panda-cooking-go-api/internal/middleware"
	"panda-cooking-go-api/internal/service"

	"github.com/gin-gonic/gin"
)

type CommentHandler struct {
	service *service.CommentService
}

func NewCommentHandler(svc *service.CommentService) *CommentHandler {
	return &CommentHandler{service: svc}
}

// ListByRecipe lista os comentários de uma receita (público, paginado).
func (h *CommentHandler) ListByRecipe(c *gin.Context) {
	var q service.PageQuery
	if !bindQuery(c, &q) {
		return
	}

	page, err := h.service.ListByRecipe(c.Request.Context(), c.Param("id"), q)
	if err != nil {
		respondError(c, err)
		return
	}
	c.JSON(http.StatusOK, page)
}

func (h *CommentHandler) Create(c *gin.Context) {
	var input service.CommentInput
	if !bindJSON(c, &input) {
		return
	}

	comment, err := h.service.Create(c.Request.Context(), c.GetString(middleware.UserIDKey), c.Param("id"), input)
	if err != nil {
		respondError(c, err)
		return
	}
	c.JSON(http.StatusCreated, comment)
}

func (h *CommentHandler) Update(c *gin.Context) {
	id, ok := idParam(c, "id", service.ErrCommentNotFound)
	if !ok {
		return
	}
	var input service.CommentInput
	if !bindJSON(c, &input) {
		return
	}

	comment, err := h.service.Update(c.Request.Context(), id, c.GetString(middleware.UserIDKey), input)
	if err != nil {
		respondError(c, err)
		return
	}
	c.JSON(http.StatusOK, comment)
}

func (h *CommentHandler) Delete(c *gin.Context) {
	id, ok := idParam(c, "id", service.ErrCommentNotFound)
	if !ok {
		return
	}

	err := h.service.Delete(c.Request.Context(), id, c.GetString(middleware.UserIDKey), c.GetBool(middleware.IsAdmKey))
	if err != nil {
		respondError(c, err)
		return
	}
	c.Status(http.StatusNoContent)
}
