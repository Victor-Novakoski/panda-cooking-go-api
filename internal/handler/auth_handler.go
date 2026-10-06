package handler

import (
	"net/http"

	"panda-cooking-go-api/internal/service"

	"github.com/gin-gonic/gin"
)

type AuthHandler struct {
	service *service.UserService
}

func NewAuthHandler(service *service.UserService) *AuthHandler {
	return &AuthHandler{service: service}
}

// RegisterRoutes registra o login. Os middlewares recebidos (ex.: limite por IP)
// rodam antes dele.
func (h *AuthHandler) RegisterRoutes(r *gin.RouterGroup, middlewares ...gin.HandlerFunc) {
	r.POST("", append(middlewares, h.login)...)
}

func (h *AuthHandler) login(c *gin.Context) {
	var input service.LoginInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	resp, err := h.service.Login(input)
	if err != nil {
		respondError(c, err)
		return
	}

	c.JSON(http.StatusOK, resp)
}
