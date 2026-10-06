// Recria os dados de demonstração. Apaga tudo o que estiver no banco.
package main

import (
	"log"

	"panda-cooking-go-api/internal/config"
	"panda-cooking-go-api/internal/database"
	"panda-cooking-go-api/internal/seed"

	"github.com/joho/godotenv"
)

func main() {
	_ = godotenv.Load()

	db, err := database.Connect(config.LoadDB())
	if err != nil {
		log.Fatalf("erro ao conectar ao banco: %v", err)
	}

	if _, err := seed.Demo(db, true); err != nil {
		log.Fatalf("erro no seed: %v", err)
	}
	log.Printf("dados de demonstração recriados (senha dos usuários: %s)", seed.DemoPassword)
}
