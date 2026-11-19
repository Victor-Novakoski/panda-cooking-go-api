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
| Autenticação | JWT HS256 de 15 min + refresh token em cookie `HttpOnly`, guardado como hash; senha com bcrypt |
| Logs | `log/slog` em JSON (texto em desenvolvimento), com id da requisição |
| Testes | `testing` + testify; mocks escritos à mão; testcontainers para os de integração |
| Contrato | OpenAPI 3 em `api/openapi.yaml`, servido em `/api/openapi.yaml` |
| Dev | Docker Compose (Postgres, API com air e front com `next dev`) |

## Camadas

```
cmd/main.go            config → banco (migrations) → seed → server; desliga devagar no SIGTERM
cmd/seed               recria os dados de demonstração (apaga o banco)
internal/config        lê e valida as variáveis de ambiente
internal/database      conexão, pool e migrations
internal/server        middlewares globais, rotas e ligação dos handlers
internal/handler       HTTP: lê a requisição, chama o service, devolve JSON
internal/service       regra de negócio, validação das entradas e DTOs de resposta
internal/repository    acesso ao banco, atrás de interfaces
internal/middleware    JWT, limite por IP, headers de segurança, id da requisição, log, prazo
internal/seed          usuários e receitas de demonstração
internal/integration   testes com a API e um Postgres de verdade
pkg/token              gerar e validar o JWT
```

O handler só fala HTTP, a regra de negócio (dono da receita, admin) fica no service, e o service depende das interfaces do repository e devolve DTOs, nunca o model.

## Modelo de dados

- `users` 1:N `recipes`, `comments`, `favorite_recipes`, `sessions`, `uploads`
- `sessions` 1:N `refresh_tokens` (só o hash SHA-256 é guardado)
- `categories` 1:N `recipes` (categorias fixas, criadas pela migration)
- `recipes` 1:N `image_recipes`, `preparations`, `comments`; N:N `ingredients` via `ingredient_recipes`
- `uploads`: fotos enviadas do computador, guardadas no próprio banco (até 900 KB)

IDs de usuário, receita, sessão e foto enviada são UUID; o resto é inteiro sequencial. Apagar usuário ou receita apaga em cascata o que depende deles. O schema de verdade está em `internal/database/migrations/`.
