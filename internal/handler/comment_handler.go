package handler

import (
	"net/http"
	"panda-cooking-go-api/internal/middleware"
	"panda-cooking-go-api/internal/service"
	"strconv"

	"github.com/gin-gonic/gin"
)

type CommentHandler struct {
	service *service.CommentService
}

func NewCommentHandler(service *service.CommentService) *CommentHandler {
	return &CommentHandler{service: service}
}

func (h *CommentHandler) RegisterRoutes(r *gin.RouterGroup, authMiddleware gin.HandlerFunc) {
	r.POST("", authMiddleware, h.create)
	r.GET("", h.getAll)
	r.PATCH("/:id", authMiddleware, h.update)
	r.DELETE("/:id", authMiddleware, h.delete)
}

func (h *CommentHandler) create(c *gin.Context) {
	var input service.CreateCommentInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userID := c.GetString(middleware.UserIDKey)
	comment, err := h.service.Create(userID, input)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, comment)
}

func (h *CommentHandler) getAll(c *gin.Context) {
	comments, err := h.service.GetAll()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, comments)
}

func (h *CommentHandler) update(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "id inválido"})
		return
	}

	var input service.UpdateCommentInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userID := c.GetString(middleware.UserIDKey)
	comment, err := h.service.Update(uint(id), userID, input)
	if err != nil {
		c.JSON(statusFromErr(err), gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, comment)
}

func (h *CommentHandler) delete(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "id inválido"})
		return
	}

	userID := c.GetString(middleware.UserIDKey)
	isAdm := c.GetBool(middleware.IsAdmKey)

	if err := h.service.Delete(uint(id), userID, isAdm); err != nil {
		c.JSON(statusFromErr(err), gin.H{"error": err.Error()})
		return
	}

	c.Status(http.StatusNoContent)
}
