// Command seed apaga todos os dados do banco e recria os de demonstração.
// Só para desenvolvimento.
package main

import (
	"context"
	"log/slog"
	"os"

	"panda-cooking-go-api/internal/applog"
	"panda-cooking-go-api/internal/config"
	"panda-cooking-go-api/internal/database"
	"panda-cooking-go-api/internal/seed"

	"github.com/joho/godotenv"
)

func main() {
	_ = godotenv.Load()
	log := applog.New(false, os.Stdout)
	ctx := context.Background()
	url := config.LoadDB().URL()

	db, err := database.Connect(ctx, url, log)
	if err != nil {
		fail(log, "conectar ao banco", err)
	}
	if err := database.Migrate(ctx, db, url); err != nil {
		fail(log, "aplicar migrations", err)
	}
	if _, err := seed.Demo(ctx, db, true); err != nil {
		fail(log, "criar dados de demonstração", err)
	}
	log.Info("dados de demonstração recriados", "senha_dos_usuarios", seed.DemoPassword)
}

func fail(log *slog.Logger, step string, err error) {
	log.Error("seed falhou", "etapa", step, "erro", err)
	os.Exit(1)
}
