package main

import (
	"log"
	"panda-cooking-go-api/internal/config"
	"panda-cooking-go-api/internal/database"
	"panda-cooking-go-api/internal/handler"
	"panda-cooking-go-api/internal/middleware"
	"panda-cooking-go-api/internal/repository"
	"panda-cooking-go-api/internal/service"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Println("arquivo .env não encontrado, usando variáveis de ambiente do sistema")
	}

	cfg := config.Load()

	db, err := database.Connect(cfg.DB)
	if err != nil {
		log.Fatalf("erro ao conectar ao banco: %v", err)
	}
	log.Println("banco de dados conectado")
	database.Seed(db)

	// Repositories
	userRepo     := repository.NewUserRepository(db)
	recipeRepo   := repository.NewRecipeRepository(db)
	categoryRepo := repository.NewCategoryRepository(db)
	commentRepo  := repository.NewCommentRepository(db)
	favoriteRepo := repository.NewFavoriteRepository(db)

	// Services
	userService     := service.NewUserService(userRepo, cfg.SecretKey)
	recipeService   := service.NewRecipeService(recipeRepo)
	categoryService := service.NewCategoryService(categoryRepo)
	commentService  := service.NewCommentService(commentRepo, recipeRepo)
	favoriteService := service.NewFavoriteService(favoriteRepo, recipeRepo)

	// Handlers
	authHandler     := handler.NewAuthHandler(userService)
	userHandler     := handler.NewUserHandler(userService)
	recipeHandler   := handler.NewRecipeHandler(recipeService)
	categoryHandler := handler.NewCategoryHandler(categoryService)
	commentHandler  := handler.NewCommentHandler(commentService)
	favoriteHandler := handler.NewFavoriteHandler(favoriteService)

	authMiddleware := middleware.Auth(cfg.SecretKey)

	r := gin.Default()

	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:3000", "http://localhost:3001"},
		AllowMethods:     []string{"GET", "POST", "PATCH", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		AllowCredentials: true,
	}))

	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok"})
	})

	authHandler.RegisterRoutes(r.Group("/auth"))
	userHandler.RegisterRoutes(r.Group("/users"), authMiddleware)
	recipeHandler.RegisterRoutes(r.Group("/recipes"), authMiddleware)
	categoryHandler.RegisterRoutes(r.Group("/categories"))
	commentHandler.RegisterRoutes(r.Group("/comments"), authMiddleware)
	favoriteHandler.RegisterRoutes(r.Group("/favorites"), authMiddleware)

	log.Printf("servidor rodando na porta %s", cfg.Port)
	if err := r.Run(":" + cfg.Port); err != nil {
		log.Fatalf("erro ao iniciar servidor: %v", err)
	}
}
