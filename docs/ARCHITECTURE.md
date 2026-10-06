# Arquitetura

## Visão geral

```
Front (Next.js) ──HTTP/JSON──▶ API Go (Gin) ──GORM──▶ PostgreSQL 16
```

O front fica no repositório [panda-cooking-front](https://github.com/Victor-Novakoski/panda-cooking-front).

## Stack

| Parte | Escolha |
| --- | --- |
| Linguagem | Go 1.26 |
| HTTP | Gin + gin-contrib/cors |
| Banco | PostgreSQL 16 com GORM (`AutoMigrate` ao subir) |
| Autenticação | JWT HS256 (golang-jwt), senha com bcrypt |
| Configuração | variáveis de ambiente (`.env` via godotenv) |
| Testes | `testing` + testify, mocks escritos à mão |
| Dev | Docker Compose (Postgres) + air (hot reload) |

## Camadas

```
cmd/main.go            monta tudo: config → banco → repositories → services → handlers → rotas
cmd/seed               popula o banco com receitas de exemplo
internal/config        lê as variáveis de ambiente
internal/database      conexão, AutoMigrate e categorias padrão
internal/model         structs GORM (tabelas)
internal/repository    acesso ao banco; interfaces em interfaces.go
internal/service       regra de negócio, DTOs de entrada e resposta
internal/handler       HTTP: lê a requisição, chama o service, devolve JSON
internal/middleware    autenticação JWT
pkg/token              gerar e validar o JWT
```

Regras:

- Handler só fala HTTP. Regra de negócio (dono da receita, admin) fica no service.
- O service depende das **interfaces** do repository, por isso os testes usam os mocks de `internal/service/mocks`.
- O service devolve DTOs (`UserResponse`, `RecipeResponse`...), nunca o model direto, para não vazar campos como `Password`.

## Modelo de dados

- `users` 1:N `recipes`, `comments`, `favorite_recipes`
- `categories` 1:N `recipes`
- `recipes` 1:N `image_recipes`, `preparations`, `comments`
- `recipes` N:N `ingredients` via `ingredient_recipes` (com a quantidade)

IDs de usuário e receita são UUID; o resto é inteiro sequencial.

## Rotas

| Método | Rota | Auth |
| --- | --- | --- |
| GET | `/health` | — |
| POST | `/auth` | — |
| POST | `/users` | — |
| GET, PATCH, DELETE | `/users/profile` | ✅ |
| GET | `/users/profile/favorite-recipes` | ✅ |
| GET | `/categories` | — |
| GET | `/recipes`, `/recipes/:id` | — |
| POST | `/recipes` | ✅ |
| PATCH, DELETE | `/recipes/:id` | ✅ dono |
| POST | `/recipes/:id/images`, `/ingredients`, `/preparations` | ✅ dono |
| PATCH, DELETE | `/recipes/:id/images/:imageID`, `/preparations/:prepID` | ✅ dono |
| DELETE | `/recipes/:id/ingredients/:ingredientID` | ✅ dono |
| GET | `/comments` | — |
| POST | `/comments` | ✅ |
| PATCH | `/comments/:id` | ✅ dono |
| DELETE | `/comments/:id` | ✅ dono ou admin |
| POST, DELETE | `/favorites/:recipeID` | ✅ |

A coleção do Postman (`panda-cooking.postman_collection.json`) tem exemplos de todas.
