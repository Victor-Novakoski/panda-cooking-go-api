# Panda Cooking — API

API REST de uma rede de receitas: cadastro e login, receitas com fotos, ingredientes e modo de preparo, categorias, comentários e favoritos.

Front-end: [panda-cooking-front](https://github.com/Victor-Novakoski/panda-cooking-front).

[![CI](https://github.com/Victor-Novakoski/panda-cooking-go-api/actions/workflows/ci.yml/badge.svg?branch=develop)](https://github.com/Victor-Novakoski/panda-cooking-go-api/actions/workflows/ci.yml)

## Stack

Go 1.26 · Gin · GORM · PostgreSQL 16 · JWT · Docker · GitHub Actions

## Rodando local

Precisa só de Docker com Compose. Clone a API e o front lado a lado:

```bash
git clone https://github.com/Victor-Novakoski/panda-cooking-go-api.git
git clone https://github.com/Victor-Novakoski/panda-cooking-front.git
cd panda-cooking-go-api
cp .env.example .env
docker compose up --build
```

| O quê | Endereço |
| --- | --- |
| Front | http://localhost:3000 |
| API | http://localhost:8080/health |
| Postgres | `localhost:5433` (`postgres` / `postgres`) |

Na primeira subida a API cria dez receitas e três usuários. Todos entram com a senha `panda-cooking-demo`:

| Usuário | E-mail | Admin |
| --- | --- | --- |
| Chef Maria Silva | `maria@pandacooking.com` | sim |
| João Cozinheiro | `joao@pandacooking.com` | não |
| Ana Paula Gourmet | `ana@pandacooking.com` | não |

API e front recarregam sozinhos ao salvar um arquivo. Se o front estiver em outra pasta, aponte `FRONT_DIR` no `.env`; para subir só banco e API, `docker compose up api`.

| Comando | O que faz |
| --- | --- |
| `make dev` | sobe banco, API e front |
| `make down` | para tudo (os dados ficam; `docker compose down -v` apaga) |
| `make seed` | apaga o banco e recria os dados de demonstração |
| `make test` | testes com detector de race (precisa de Go 1.26+) |
| `make lint` | golangci-lint |

## Documentação

- [PRD](docs/PRD.md): o que o produto faz
- [Arquitetura](docs/ARCHITECTURE.md): camadas, modelo de dados e rotas
- [Convenções da API](docs/DESIGN.md)
- [Regras](docs/RULES.md): código, segurança e fluxo de git
- [Segurança](docs/SECURITY.md)
- [Tarefas](docs/TASKS.md) e [Memória](docs/MEMORY.md)

A coleção do Postman está em `panda-cooking.postman_collection.json`.
