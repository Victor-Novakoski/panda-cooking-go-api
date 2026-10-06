.PHONY: help setup dev down logs seed test lint build run clean deps

help: ## Mostra este help
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-20s\033[0m %s\n", $$1, $$2}'

setup: ## Primeira vez: cria o .env e sobe banco, API e front
	@test -f .env || cp .env.example .env
	docker compose up --build

dev: ## Sobe banco, API e front com hot reload
	docker compose up

down: ## Para os containers (os dados do banco ficam)
	docker compose down

logs: ## Mostra os logs dos containers
	docker compose logs -f

seed: ## Apaga o banco e recria os dados de demonstração
	docker compose exec api go run ./cmd/seed

test: ## Roda os testes (com detector de race, igual à CI)
	go test -race -count=1 ./...

lint: ## Roda o golangci-lint (mesma configuração da CI)
	golangci-lint run ./...

build: ## Compila a aplicação
	go build -o bin/api ./cmd

run: ## Executa a aplicação fora do Docker (banco: docker compose up -d postgres)
	go run ./cmd

clean: ## Remove arquivos temporários e build
	rm -rf bin/ .air/

deps: ## Baixa as dependências
	go mod download
	go mod tidy
