# Tarefas

Backlog em ordem. Só se trabalha na etapa atual; o que surgir no caminho entra aqui antes de ser feito ([RULES.md](RULES.md#1-escopo)). Os números `#N` apontam para os itens de [SECURITY.md](SECURITY.md).

**Etapa atual: 4 — Deploy** (as etapas 1 a 3 e 5 estão prontas)

## Etapa 0 — Base da API ✅

- [x] Gin + GORM + Postgres, camadas handler → service → repository
- [x] Cadastro, login JWT e perfil
- [x] Receitas com fotos, ingredientes e modo de preparo
- [x] Categorias, comentários e favoritos
- [x] Seed com receitas de exemplo
- [x] Testes de service (mocks) e de handler

## Etapa 1 — Setup do repositório

- [x] Docs em `docs/` (PRD, arquitetura, regras, convenções, tarefas, memória, segurança)
- [x] CI: lint, testes com race, imagem Docker + Trivy, govulncheck, gitleaks, título do PR
- [x] Dependabot para Go, Actions e Docker
- [x] golangci-lint configurado e código ajustado
- [x] Dockerfile multi-stage rodando sem root (distroless)
- [x] Branch `develop`, proteção da `main` e da `develop` com os checks obrigatórios (feito pelo victor no GitHub)

## Etapa 2 — Segurança da base

- [x] Item da receita (foto, passo, ingrediente) só pode ser alterado pela receita a que pertence, com teste (#6)
- [x] Erros tipados no service; 500 com mensagem genérica e detalhe no log; e-mail duplicado vira 409 (#12)
- [x] Id que não é UUID responde 404 em vez de 500 (#12)
- [x] `SECRET_KEY` obrigatória com 32+ caracteres; produção recusa a de exemplo (#1)
- [x] Tamanho máximo nos textos, trim, e-mail minúsculo, limite do corpo (#3)
- [x] Senha: mínimo 10, máximo 72 bytes, recusa de senhas comuns; tempo constante no login (#5)
- [x] Limite no login: 10 por minuto por IP e bloqueio de 15 min depois de 5 senhas erradas no mesmo e-mail (#8)
- [x] Rate limit global (#15)
- [x] Erro de validação sem nomes internos de struct: 422 com o erro de cada campo, ver [DESIGN.md](DESIGN.md) (#12)
- [x] `http.Server` com timeouts (#20)
- [x] CORS por `CORS_ORIGINS` (#17) e middleware de headers de segurança (#19)
- [x] Postgres do compose só em `127.0.0.1` (#21)
- [x] Logar falha de login, 401, 403 e 429 (#22)

## Etapa 2.1 — Tudo com um `docker compose up`

- [x] Compose com banco, API (air) e front (`next dev`), front clonado ao lado
- [x] Seed com usuários que conseguem logar (senha com bcrypt), criado sozinho ao subir com `SEED_DEMO=true`
- [x] `make seed` roda dentro do container, sem precisar de Go na máquina

## Etapa 3 — API completa

- [x] Paginação em receitas e comentários (#20)
- [x] Busca por nome e filtro por categoria na API
- [x] Comentários por receita (`GET /recipes/:id/comments`), autor sem e-mail
- [x] `PUT /recipes/:id` troca a receita inteira numa transação (tela de edição do front)
- [x] `PATCH /recipes/:id` volta a trocar a categoria; categoria inexistente responde 400
- [x] Perfil: foto vazia remove a foto, nome vazio é recusado
- [x] Migrations versionadas em SQL (golang-migrate) no lugar do `AutoMigrate`
- [x] Testes de integração com Postgres real (testcontainers, job `Testes de integração` na CI)
- [x] Documentação OpenAPI (`api/openapi.yaml`, servida em `/api/openapi.yaml`, com teste que confere todas as rotas)

## Etapa 4 — Deploy

- [ ] Publicar junto com o Rastreia, sem custo
  - `CORS_ORIGINS` com a origem `https` do front; `TRUSTED_PROXIES` com a rede do front e do proxy reverso
  - proxy reverso sobrescrevendo o `X-Forwarded-For` (#23) e servindo em HTTPS (cookie `Secure`)
  - banco sem porta publicada; `DB_SSLMODE` conforme o banco

## Etapa 5 — Sessão em cookie, API completa e revisão geral ✅

- [x] Rotas sob `/api` (o `/health` continua na raiz)
- [x] Access token de 15 min + refresh token em cookie `HttpOnly` com rotação, detecção de reúso e revogação no logout (#14)
- [x] Cookie `panda_session` para o front saber no servidor que existe sessão
- [x] `http.CrossOriginProtection` nas requisições que mudam dados (#10)
- [x] `context.Context` da requisição até o banco, com prazo de 10 s por requisição (#20)
- [x] Limites de tamanho em todos os campos e listas, iguais aos do front (#3)
- [x] Busca sem diferenciar acento (`unaccent` + `pg_trgm`), filtros por categoria e autor, paginação com total
- [x] `GET /api/favorites/:recipeID` para o front saber se a receita é favorita sem baixar a lista
- [x] Log estruturado com `slog` e id da requisição; eventos de segurança em `WARN` (#22)
- [x] Desligamento gracioso no `SIGTERM`; `healthcheck` no próprio binário para o Docker
- [x] Seed com comentários e favoritos; coleção do Postman trocada pelo OpenAPI
