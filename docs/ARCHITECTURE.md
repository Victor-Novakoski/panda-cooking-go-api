# Arquitetura

## Visão geral

```
Navegador ──▶ Front (Next.js) ──/api/*──▶ API Go (Gin) ──GORM──▶ PostgreSQL 16
```

O front fica no repositório [panda-cooking-front](https://github.com/Victor-Novakoski/panda-cooking-front). O navegador só fala com o front, que repassa `/api/*` para a API pelo servidor dele: para o navegador é tudo a mesma origem, e o cookie da sessão vai sem CORS. A API também aceita chamada direta de origens listadas em `CORS_ORIGINS`.

## Stack

| Parte | Escolha |
| --- | --- |
| Linguagem | Go 1.26 |
| HTTP | Gin, `http.Server` com timeouts, `http.CrossOriginProtection` e gin-contrib/cors |
| Banco | PostgreSQL 16 com GORM; migrations em SQL com golang-migrate, embutidas no binário |
| Busca | `unaccent` + `pg_trgm`, com índice GIN |
| Autenticação | JWT HS256 de 15 min (golang-jwt) + refresh token em cookie `HttpOnly`, guardado como hash; senha com bcrypt |
| Logs | `log/slog` em JSON (texto em desenvolvimento), com id da requisição |
| Configuração | variáveis de ambiente (`.env` via godotenv) |
| Testes | `testing` + testify; mocks escritos à mão; testcontainers para os de integração |
| Documentação | OpenAPI 3 em `api/openapi.yaml` |
| Dev | Docker Compose (Postgres, API com air e front com `next dev`) |

## Camadas

```
cmd/main.go            config → banco (migrations) → seed → server; desliga devagar no SIGTERM
cmd/seed               recria os dados de demonstração (apaga o banco)
api                    openapi.yaml, embutido no binário
internal/config        lê e valida as variáveis de ambiente
internal/database      conexão, pool e migrations (internal/database/migrations/*.sql)
internal/seed          usuários e receitas de demonstração (SEED_DEMO=true ao subir, ou cmd/seed)
internal/server        middlewares globais, tabela de rotas e ligação dos handlers
internal/model         structs GORM (tabelas)
internal/repository    acesso ao banco; interfaces em interfaces.go
internal/service       regra de negócio, validação das entradas (inputs.go) e DTOs de resposta
internal/handler       HTTP: lê a requisição, chama o service, devolve JSON; 422 por campo
internal/middleware    JWT, limite por IP, CSRF, headers de segurança, id da requisição, log, prazo
internal/ratelimit     contador de tentativas por chave (IP, e-mail) em memória
internal/integration   testes com a API e um Postgres de verdade (build tag integration)
pkg/token              gerar e validar o JWT
```

Regras:

- Handler só fala HTTP. Regra de negócio (dono da receita, admin) fica no service.
- O service depende das **interfaces** do repository, por isso os testes de unidade usam os mocks de `internal/service/mocks`.
- O service devolve DTOs (`UserResponse`, `RecipeResponse`...), nunca o model direto, para não vazar campos como `PasswordHash`.
- Erro esperado sai do service como `service.Error` (ex.: `ErrRecipeNotFound`), e o handler converte em status com `respondError`. Erro de outro tipo vira 500 genérico.
- O `context.Context` da requisição vai até o banco: cliente que desiste ou requisição que passa de 10 segundos cancela a consulta.

## Ordem dos middlewares

Recuperação de panic → id da requisição → log de acesso → headers de segurança → proteção contra CSRF → CORS → limite de corpo (1 MB) → prazo (10 s) → limite por IP (300/min em `/api`, mais os das rotas de login, renovação e cadastro) → JWT nas rotas com login.

## Modelo de dados

- `users` 1:N `recipes`, `comments`, `favorite_recipes`, `sessions`
- `sessions` 1:N `refresh_tokens` (um por renovação; só o hash SHA-256 é guardado)
- `categories` 1:N `recipes`
- `recipes` 1:N `image_recipes`, `preparations`, `comments`
- `recipes` N:N `ingredients` via `ingredient_recipes` (com a quantidade)

IDs de usuário, receita e sessão são UUID; o resto é inteiro sequencial. Apagar usuário ou receita apaga em cascata o que depende deles.

## Rotas

Todas sob `/api`, menos o `/health`. Detalhes de corpo e resposta em [`api/openapi.yaml`](../api/openapi.yaml).

| Método | Rota | Auth |
| --- | --- | --- |
| GET | `/health` | — |
| GET | `/api/openapi.yaml` | — |
| POST | `/api/auth/login`, `/api/auth/refresh`, `/api/auth/logout` | cookie (refresh e logout) |
| POST | `/api/users` (cadastro) | — |
| GET, PATCH, DELETE | `/api/users/profile` | ✅ |
| GET | `/api/users/profile/favorites` | ✅ |
| GET | `/api/categories` | — |
| GET | `/api/recipes` (busca, filtros e paginação), `/api/recipes/:id` | — |
| GET | `/api/recipes/:id/comments` (paginado) | — |
| POST | `/api/recipes` | ✅ |
| PUT | `/api/recipes/:id` (receita inteira, numa transação) | ✅ dono |
| PATCH, DELETE | `/api/recipes/:id` | ✅ dono |
| POST | `/api/recipes/:id/images`, `/ingredients`, `/preparations` | ✅ dono |
| PATCH, DELETE | `/api/recipes/:id/images/:imageID`, `/preparations/:prepID` | ✅ dono |
| DELETE | `/api/recipes/:id/ingredients/:ingredientID` | ✅ dono |
| POST | `/api/recipes/:id/comments` | ✅ |
| PATCH | `/api/comments/:id` | ✅ autor |
| DELETE | `/api/comments/:id` | ✅ autor ou admin |
| GET, POST, DELETE | `/api/favorites/:recipeID` | ✅ |
