// Package server monta a API: os middlewares que valem para toda requisição,
// a tabela de rotas e a ligação de cada handler com o banco.
package server

import (
	"context"
	"log/slog"
	"net/http"
	"time"

	"panda-cooking-go-api/internal/config"
	"panda-cooking-go-api/internal/database"
	"panda-cooking-go-api/internal/handler"
	"panda-cooking-go-api/internal/middleware"
	"panda-cooking-go-api/internal/ratelimit"
	"panda-cooking-go-api/internal/repository"
	"panda-cooking-go-api/internal/service"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// Limites de cada requisição e de requisições por IP.
const (
	MaxBodyBytes   = 1 << 20 // 1 MB: uma receita com 50 ingredientes e 50 passos cabe com folga
	RequestTimeout = 10 * time.Second

	GlobalLimit  = 300 // por minuto em /api, para quem varre a API
	LoginLimit   = 10  // por minuto; o bloqueio por e-mail fica no AuthService
	RefreshLimit = 30  // por minuto
	SignupLimit  = 10  // por hora
)

// Handlers são os handlers que as rotas usam.
type Handlers struct {
	Auth     *handler.AuthHandler
	User     *handler.UserHandler
	Recipe   *handler.RecipeHandler
	Category *handler.CategoryHandler
	Comment  *handler.CommentHandler
	Favorite *handler.FavoriteHandler
	Health   gin.HandlerFunc
	OpenAPI  gin.HandlerFunc
}

// Guards são os middlewares de rota: login obrigatório e limites por IP.
type Guards struct {
	Auth    gin.HandlerFunc
	Global  gin.HandlerFunc
	Login   gin.HandlerFunc
	Refresh gin.HandlerFunc
	Signup  gin.HandlerFunc
}

// Deps é o que a API precisa para rodar: os repositórios, o teste de vida
// do banco (para o /health), o logger e a especificação OpenAPI servida em
// /api/openapi.yaml. Em produção vem de FromDB; nos testes, de mocks.
type Deps struct {
	Users      repository.UserRepo
	Recipes    repository.RecipeRepo
	Categories repository.CategoryRepo
	Comments   repository.CommentRepo
	Favorites  repository.FavoriteRepo
	Sessions   repository.SessionRepo
	Ping       func(context.Context) error
	Log        *slog.Logger
	Spec       []byte
}

// FromDB liga os repositórios ao banco.
func FromDB(db *gorm.DB, log *slog.Logger, spec []byte) Deps {
	return Deps{
		Users:      repository.NewUserRepository(db),
		Recipes:    repository.NewRecipeRepository(db),
		Categories: repository.NewCategoryRepository(db),
		Comments:   repository.NewCommentRepository(db),
		Favorites:  repository.NewFavoriteRepository(db),
		Sessions:   repository.NewSessionRepository(db),
		Ping:       func(ctx context.Context) error { return database.Ping(ctx, db) },
		Log:        log,
		Spec:       spec,
	}
}

// New monta a API: services, handlers, middlewares e rotas.
func New(cfg config.Config, d Deps) (*gin.Engine, error) {
	favorites := service.NewFavoriteService(d.Favorites, d.Recipes)
	secure := cfg.IsProduction()
	h := Handlers{
		Auth:     handler.NewAuthHandler(service.NewAuthService(d.Users, d.Sessions, cfg.SecretKey, d.Log), secure),
		User:     handler.NewUserHandler(service.NewUserService(d.Users), favorites, secure),
		Recipe:   handler.NewRecipeHandler(service.NewRecipeService(d.Recipes, d.Categories)),
		Category: handler.NewCategoryHandler(service.NewCategoryService(d.Categories)),
		Comment:  handler.NewCommentHandler(service.NewCommentService(d.Comments, d.Recipes)),
		Favorite: handler.NewFavoriteHandler(favorites),
		Health:   handler.Health(d.Ping),
		OpenAPI:  handler.OpenAPI(d.Spec),
	}

	r, err := Engine(cfg, d.Log)
	if err != nil {
		return nil, err
	}
	Routes(r, h, NewGuards(cfg.SecretKey))
	return r, nil
}

// NewGuards cria o login obrigatório e os limitadores por IP, cada um com o
// seu contador.
func NewGuards(secretKey string) Guards {
	return Guards{
		Auth:    middleware.Auth(secretKey),
		Global:  middleware.RateLimitByIP(ratelimit.New(GlobalLimit, time.Minute)),
		Login:   middleware.RateLimitByIP(ratelimit.New(LoginLimit, time.Minute)),
		Refresh: middleware.RateLimitByIP(ratelimit.New(RefreshLimit, time.Minute)),
		Signup:  middleware.RateLimitByIP(ratelimit.New(SignupLimit, time.Hour)),
	}
}

// Engine cria o roteador com os middlewares que valem para toda requisição,
// na ordem em que rodam: recuperação de panic, id da requisição, log de
// acesso, cabeçalhos de segurança, proteção contra CSRF, CORS, limite de
// corpo e prazo.
func Engine(cfg config.Config, log *slog.Logger) (*gin.Engine, error) {
	if cfg.IsProduction() {
		gin.SetMode(gin.ReleaseMode)
	}

	handler.SetupValidator()

	r := gin.New()
	r.HandleMethodNotAllowed = true
	// Só os proxies listados podem dizer o IP do cliente pelo
	// X-Forwarded-For. Sem lista, vale o IP da conexão: ninguém burla o
	// limite por IP mandando o cabeçalho.
	if err := r.SetTrustedProxies(cfg.TrustedProxies); err != nil {
		return nil, err
	}

	crossOrigin, err := middleware.CrossOrigin(cfg.CORSOrigins)
	if err != nil {
		return nil, err
	}

	r.Use(
		middleware.Recovery(log),
		middleware.RequestID(),
		middleware.AccessLog(log),
		middleware.SecurityHeaders(cfg.IsProduction()),
		// antes do CORS, para a recusa sair no formato de erro da API
		crossOrigin,
		cors.New(cors.Config{
			AllowOrigins:     cfg.CORSOrigins,
			AllowMethods:     []string{http.MethodGet, http.MethodPost, http.MethodPut, http.MethodPatch, http.MethodDelete},
			AllowHeaders:     []string{"Content-Type", "Authorization", middleware.RequestIDHeader},
			ExposeHeaders:    []string{middleware.RequestIDHeader, "Retry-After"},
			AllowCredentials: true,
			MaxAge:           12 * time.Hour,
		}),
		middleware.BodyLimit(MaxBodyBytes),
		middleware.Timeout(RequestTimeout),
	)
	return r, nil
}

// Routes registra todas as rotas. Os testes usam a mesma função, então a
// tabela testada é a que roda.
func Routes(r *gin.Engine, h Handlers, g Guards) {
	r.GET("/health", h.Health)
	r.NoRoute(handler.NotFound)
	r.NoMethod(handler.MethodNotAllowed)

	api := r.Group("/api", g.Global)
	api.GET("/openapi.yaml", h.OpenAPI)

	auth := api.Group("/auth")
	auth.POST("/login", g.Login, h.Auth.Login)
	auth.POST("/refresh", g.Refresh, h.Auth.Refresh)
	auth.POST("/logout", h.Auth.Logout)

	api.POST("/users", g.Signup, h.User.Create)
	profile := api.Group("/users/profile", g.Auth)
	profile.GET("", h.User.GetProfile)
	profile.PATCH("", h.User.Update)
	profile.DELETE("", h.User.Delete)
	profile.GET("/favorites", h.User.Favorites)

	api.GET("/categories", h.Category.List)

	api.GET("/recipes", h.Recipe.List)
	api.GET("/recipes/:id", h.Recipe.Get)
	api.GET("/recipes/:id/comments", h.Comment.ListByRecipe)

	recipes := api.Group("/recipes", g.Auth)
	recipes.POST("", h.Recipe.Create)
	recipes.PUT("/:id", h.Recipe.Replace)
	recipes.PATCH("/:id", h.Recipe.Update)
	recipes.DELETE("/:id", h.Recipe.Delete)
	recipes.POST("/:id/images", h.Recipe.AddImage)
	recipes.PATCH("/:id/images/:imageID", h.Recipe.UpdateImage)
	recipes.DELETE("/:id/images/:imageID", h.Recipe.DeleteImage)
	recipes.POST("/:id/ingredients", h.Recipe.AddIngredient)
	recipes.DELETE("/:id/ingredients/:ingredientID", h.Recipe.DeleteIngredient)
	recipes.POST("/:id/preparations", h.Recipe.AddPreparation)
	recipes.PATCH("/:id/preparations/:prepID", h.Recipe.UpdatePreparation)
	recipes.DELETE("/:id/preparations/:prepID", h.Recipe.DeletePreparation)
	recipes.POST("/:id/comments", h.Comment.Create)

	comments := api.Group("/comments", g.Auth)
	comments.PATCH("/:id", h.Comment.Update)
	comments.DELETE("/:id", h.Comment.Delete)

	favorites := api.Group("/favorites", g.Auth)
	favorites.GET("/:recipeID", h.Favorite.Status)
	favorites.POST("/:recipeID", h.Favorite.Add)
	favorites.DELETE("/:recipeID", h.Favorite.Remove)
}
