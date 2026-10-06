.PHONY: help setup dev test lint seed reseed build run clean deps docker-up docker-down docker-logs

help: ## Mostra este help
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-20s\033[0m %s\n", $$1, $$2}'

setup: ## Setup completo: Docker + Seed + Dev server (ideal para primeira vez)
	@echo "🐼 Configurando ambiente..."
	docker compose up -d
	@echo "⏳ Aguardando banco de dados..."
	@sleep 3
	@echo "🌱 Populando banco com receitas..."
	go run cmd/seed/main.go
	@echo "🚀 Iniciando servidor em modo desenvolvimento..."
	air

dev: ## Inicia o servidor em modo desenvolvimento com hot reload
	air

test: ## Roda os testes (com detector de race, igual à CI)
	go test -race -count=1 ./...

lint: ## Roda o golangci-lint (mesma configuração da CI)
	golangci-lint run ./...

seed: ## Popula o banco com dados de exemplo
	@echo "🐼 Populando banco de dados..."
	go run cmd/seed/main.go

reseed: docker-up seed ## Re-popula o banco (limpa e cria tudo de novo)

build: ## Compila a aplicação
	go build -o bin/api ./cmd

run: ## Executa a aplicação
	go run ./cmd

clean: ## Remove arquivos temporários e build
	rm -rf bin/ .temp/

deps: ## Baixa as dependências
	go mod download
	go mod tidy

docker-up: ## Sobe os containers Docker
	docker compose up -d

docker-down: ## Para os containers Docker
	docker compose down

docker-logs: ## Mostra os logs dos containers
	docker compose logs -f
