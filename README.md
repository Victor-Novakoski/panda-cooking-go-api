# Panda Cooking — API

API REST de uma rede de receitas: cadastro e login, receitas com fotos, ingredientes e modo de preparo, categorias, comentários e favoritos.

Front-end: [panda-cooking-front](https://github.com/Victor-Novakoski/panda-cooking-front).

[![CI](https://github.com/Victor-Novakoski/panda-cooking-go-api/actions/workflows/ci.yml/badge.svg?branch=develop)](https://github.com/Victor-Novakoski/panda-cooking-go-api/actions/workflows/ci.yml)

## Stack

Go 1.26 · Gin · GORM · PostgreSQL 16 · JWT · Docker · GitHub Actions

## Rodando local

Precisa de Go 1.26+, Docker e [air](https://github.com/air-verse/air).

```bash
cp .env.example .env
make setup      # sobe o Postgres, popula com receitas de exemplo e inicia com hot reload
```

A API fica em `http://localhost:8080` (`GET /health`). Depois da primeira vez, `make dev` basta.

| Comando | O que faz |
| --- | --- |
| `make dev` | API com hot reload |
| `make test` | testes com detector de race |
| `make lint` | golangci-lint |
| `make seed` | recria os dados de exemplo |

## Documentação

- [PRD](docs/PRD.md): o que o produto faz
- [Arquitetura](docs/ARCHITECTURE.md): camadas, modelo de dados e rotas
- [Convenções da API](docs/DESIGN.md)
- [Regras](docs/RULES.md): código, segurança e fluxo de git
- [Segurança](docs/SECURITY.md)
- [Tarefas](docs/TASKS.md) e [Memória](docs/MEMORY.md)

A coleção do Postman está em `panda-cooking.postman_collection.json`.
