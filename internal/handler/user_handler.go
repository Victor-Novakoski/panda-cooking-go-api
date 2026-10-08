package handler

import (
	"net/http"

	"panda-cooking-go-api/internal/middleware"
	"panda-cooking-go-api/internal/service"

	"github.com/gin-gonic/gin"
)

type UserHandler struct {
	users         *service.UserService
	favorites     *service.FavoriteService
	secureCookies bool
}

func NewUserHandler(users *service.UserService, favorites *service.FavoriteService, secureCookies bool) *UserHandler {
	return &UserHandler{users: users, favorites: favorites, secureCookies: secureCookies}
}

// Create é o cadastro. Não abre sessão: o front manda para o login.
func (h *UserHandler) Create(c *gin.Context) {
	var input service.CreateUserInput
	if !bindJSON(c, &input) {
		return
	}

	user, err := h.users.Create(c.Request.Context(), input)
	if err != nil {
		respondError(c, err)
		return
	}
	c.JSON(http.StatusCreated, user)
}

func (h *UserHandler) GetProfile(c *gin.Context) {
	user, err := h.users.GetProfile(c.Request.Context(), c.GetString(middleware.UserIDKey))
	if err != nil {
		respondError(c, err)
		return
	}
	c.JSON(http.StatusOK, user)
}

func (h *UserHandler) Update(c *gin.Context) {
	var input service.UpdateUserInput
	if !bindJSON(c, &input) {
		return
	}

	user, err := h.users.Update(c.Request.Context(), c.GetString(middleware.UserIDKey), input)
	if err != nil {
		respondError(c, err)
		return
	}
	c.JSON(http.StatusOK, user)
}

// Delete apaga a conta e os cookies de sessão (as sessões somem junto com
// o usuário no banco).
func (h *UserHandler) Delete(c *gin.Context) {
	if err := h.users.Delete(c.Request.Context(), c.GetString(middleware.UserIDKey)); err != nil {
		respondError(c, err)
		return
	}
	clearCookies(c, h.secureCookies)
	c.Status(http.StatusNoContent)
}

// Favorites lista as receitas favoritas de quem está logado.
func (h *UserHandler) Favorites(c *gin.Context) {
	var q service.PageQuery
	if !bindQuery(c, &q) {
		return
	}

	page, err := h.favorites.List(c.Request.Context(), c.GetString(middleware.UserIDKey), q)
	if err != nil {
		respondError(c, err)
		return
	}
	c.JSON(http.StatusOK, page)
}
