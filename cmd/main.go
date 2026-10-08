// Command api sobe a API do Panda Cooking.
//
// Com o argumento healthcheck, só confere se a API que já está rodando
// responde no /health: é o HEALTHCHECK do Docker, já que a imagem
// distroless não tem curl nem wget.
package main

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"

	"panda-cooking-go-api/api"
	"panda-cooking-go-api/internal/applog"
	"panda-cooking-go-api/internal/config"
	"panda-cooking-go-api/internal/database"
	"panda-cooking-go-api/internal/seed"
	"panda-cooking-go-api/internal/server"

	"github.com/joho/godotenv"
)

func main() {
	if len(os.Args) > 1 && os.Args[1] == "healthcheck" {
		os.Exit(healthcheck())
	}
	if err := run(); err != nil {
		slog.Error("a API parou", "erro", err)
		os.Exit(1)
	}
}

func run() error {
	// o .env é opcional: no Docker as variáveis já vêm do ambiente
	_ = godotenv.Load()

	cfg, err := config.Load()
	if err != nil {
		return fmt.Errorf("configuração inválida: %w", err)
	}

	log := applog.New(cfg.IsProduction(), os.Stdout)
	slog.SetDefault(log)

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	db, err := database.Connect(ctx, cfg.DB.URL(), log)
	if err != nil {
		return fmt.Errorf("conectar ao banco: %w", err)
	}
	if err := database.Migrate(ctx, db, cfg.DB.URL()); err != nil {
		return err
	}
	log.Info("banco conectado e migrations aplicadas")

	if cfg.SeedDemo {
		created, err := seed.Demo(ctx, db, false)
		if err != nil {
			return fmt.Errorf("criar dados de demonstração: %w", err)
		}
		if created {
			log.Info("dados de demonstração criados", "senha_dos_usuarios", seed.DemoPassword)
		}
	}

	router, err := server.New(cfg, server.FromDB(db, log, api.Spec))
	if err != nil {
		return fmt.Errorf("montar rotas: %w", err)
	}

	srv := &http.Server{
		Addr:    ":" + cfg.Port,
		Handler: router,
		// prazos da conexão: um cliente lento não segura o servidor
		ReadHeaderTimeout: 5 * time.Second,
		ReadTimeout:       15 * time.Second,
		WriteTimeout:      30 * time.Second,
		IdleTimeout:       2 * time.Minute,
		MaxHeaderBytes:    64 << 10,
		ErrorLog:          slog.NewLogLogger(log.Handler(), slog.LevelWarn),
	}

	serveErr := make(chan error, 1)
	go func() { serveErr <- srv.ListenAndServe() }()
	log.Info("API no ar", "porta", cfg.Port, "ambiente", cfg.Env)

	select {
	case err := <-serveErr:
		return err
	case <-ctx.Done():
	}

	// SIGTERM (docker stop): para de aceitar conexão e espera as requisições
	// em andamento terminarem
	log.Info("desligando a API")
	shutdownCtx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()
	if err := srv.Shutdown(shutdownCtx); err != nil {
		return fmt.Errorf("desligar: %w", err)
	}
	if err := <-serveErr; err != nil && !errors.Is(err, http.ErrServerClosed) {
		return err
	}
	if sqlDB, err := db.DB(); err == nil {
		_ = sqlDB.Close()
	}
	return nil
}

// healthcheck devolve 0 se a API local responde 200 no /health.
func healthcheck() int {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	if _, err := strconv.ParseUint(port, 10, 16); err != nil {
		fmt.Fprintln(os.Stderr, "healthcheck: PORT inválida")
		return 1
	}
	client := http.Client{Timeout: 3 * time.Second}
	resp, err := client.Get("http://127.0.0.1:" + port + "/health") //nolint:gosec // G704: endereço fixo da própria máquina; a porta é um número conferido acima
	if err != nil {
		fmt.Fprintln(os.Stderr, "healthcheck:", err)
		return 1
	}
	_ = resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		fmt.Fprintln(os.Stderr, "healthcheck: status", resp.StatusCode)
		return 1
	}
	return 0
}
