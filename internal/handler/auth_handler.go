package handler

import (
	"net/http"
	"time"

	"panda-cooking-go-api/internal/service"

	"github.com/gin-gonic/gin"
)

// Cookies da sessão.
const (
	// RefreshCookie guarda o refresh token: HttpOnly (o JavaScript não lê,
	// então um XSS não rouba), SameSite=Strict (não vai em requisição vinda
	// de outro site) e Path=/api/auth (só vai para as rotas de sessão).
	RefreshCookie = "panda_refresh"
	// SessionCookie não é credencial (vale "1"): só avisa o front, no
	// servidor, que existe sessão, para ele redirecionar sem piscar a tela.
	SessionCookie = "panda_session"

	refreshCookiePath = "/api/auth"
)

type AuthHandler struct {
	service *service.AuthService
	// secureCookies liga o atributo Secure (cookie só em https); desligado
	// só em desenvolvimento, que roda em http://localhost.
	secureCookies bool
}

func NewAuthHandler(svc *service.AuthService, secureCookies bool) *AuthHandler {
	return &AuthHandler{service: svc, secureCookies: secureCookies}
}

// Login confere e-mail e senha, abre a sessão e devolve o access token.
func (h *AuthHandler) Login(c *gin.Context) {
	var input service.LoginInput
	if !bindJSON(c, &input) {
		return
	}

	res, err := h.service.Login(c.Request.Context(), input)
	if err != nil {
		respondError(c, err)
		return
	}
	h.setSessionCookies(c, res.RefreshToken, res.RefreshExpiresAt)
	c.JSON(http.StatusOK, res)
}

// Refresh troca o refresh token do cookie por um novo e devolve um access
// token novo. É o que o front chama ao abrir a página para recuperar a sessão.
func (h *AuthHandler) Refresh(c *gin.Context) {
	refresh, _ := c.Cookie(RefreshCookie)

	res, err := h.service.Refresh(c.Request.Context(), refresh)
	if err != nil {
		h.clearSessionCookies(c)
		respondError(c, err)
		return
	}
	h.setSessionCookies(c, res.RefreshToken, res.RefreshExpiresAt)
	c.JSON(http.StatusOK, res)
}

// Logout encerra a sessão no banco e apaga os cookies.
func (h *AuthHandler) Logout(c *gin.Context) {
	refresh, _ := c.Cookie(RefreshCookie)
	if err := h.service.Logout(c.Request.Context(), refresh); err != nil {
		respondError(c, err)
		return
	}
	h.clearSessionCookies(c)
	c.Status(http.StatusNoContent)
}

func (h *AuthHandler) setSessionCookies(c *gin.Context, refresh string, expires time.Time) {
	maxAge := int(time.Until(expires).Seconds())
	setCookie(c, RefreshCookie, refresh, refreshCookiePath, maxAge, h.secureCookies, http.SameSiteStrictMode)
	// Lax: o aviso precisa chegar também quando a pessoa abre o site por um
	// link de fora (com Strict, o navegador não mandaria nessa navegação)
	setCookie(c, SessionCookie, "1", "/", maxAge, h.secureCookies, http.SameSiteLaxMode)
}

func (h *AuthHandler) clearSessionCookies(c *gin.Context) {
	clearCookies(c, h.secureCookies)
}

// clearCookies apaga os dois cookies (MaxAge negativo = Max-Age=0).
func clearCookies(c *gin.Context, secure bool) {
	setCookie(c, RefreshCookie, "", refreshCookiePath, -1, secure, http.SameSiteStrictMode)
	setCookie(c, SessionCookie, "", "/", -1, secure, http.SameSiteLaxMode)
}

// setCookie grava um cookie da sessão, sempre HttpOnly. Secure só fica
// desligado em desenvolvimento, que roda em http://localhost.
func setCookie(c *gin.Context, name, value, path string, maxAge int, secure bool, sameSite http.SameSite) {
	http.SetCookie(c.Writer, &http.Cookie{ //nolint:gosec // G124: Secure desligado só em desenvolvimento; o único cookie Lax (panda_session) não é credencial
		Name: name, Value: value, Path: path, MaxAge: maxAge,
		HttpOnly: true, Secure: secure, SameSite: sameSite,
	})
}
