# Panda Cooking — API

API REST de uma rede de receitas: cadastro e login, receitas com fotos, ingredientes e modo de preparo, busca, categorias, comentários e favoritos.

Esta é a API, em Go. O front, em Next.js, está em [panda-cooking-front](https://github.com/Victor-Novakoski/panda-cooking-front).

[![CI](https://github.com/Victor-Novakoski/panda-cooking-go-api/actions/workflows/ci.yml/badge.svg?branch=develop)](https://github.com/Victor-Novakoski/panda-cooking-go-api/actions/workflows/ci.yml)

## Destaques

- **Sessão segura:** access token JWT de 15 minutos e refresh token em cookie `HttpOnly`, `SameSite=Strict`, guardado no banco só como hash. Cada renovação troca o token; token reutilizado derruba a sessão inteira (sinal de cópia), com tolerância para duas abas renovando juntas.
- **Proteções contra abuso:** limite por IP em toda a API e mais apertado no login, cadastro e renovação; bloqueio de 15 minutos por e-mail depois de 5 senhas erradas; login com o mesmo tempo de resposta para e-mail existente ou não; CSRF barrado com o `http.CrossOriginProtection` do Go.
- **Validação com erro por campo:** limites em todo texto e lista, e 422 com a mensagem de cada campo no caminho do JSON (`ingredients[2].amount`), que o front mostra no lugar certo.
- **Busca que ignora acento** ("pao" acha "Pão de Queijo") com `unaccent` e índice de trigramas, mais filtro por categoria e autor e paginação com total.
- **Migrations versionadas em SQL**, embutidas no binário e aplicadas ao subir.
- **Testes em três níveis:** unidade nos services (mocks), HTTP nos handlers e integração com Postgres de verdade via testcontainers, incluindo tentativas de mexer em dado de outra pessoa.
- **Contrato em OpenAPI** (`api/openapi.yaml`, servido em `/api/openapi.yaml`), com teste que falha se uma rota não estiver documentada.
- **Pronta para produção:** timeouts no servidor e prazo por requisição até o banco, limite de corpo, headers de segurança, log estruturado com id da requisição, desligamento gracioso e imagem distroless sem root.
- **CI completa:** golangci-lint, testes com race e cobertura, testes de integração, build da imagem com Trivy, govulncheck, gitleaks e título do PR em Conventional Commits.

Detalhes de cada risco e de como foi tratado em [docs/SECURITY.md](docs/SECURITY.md).

## Stack

Go 1.26 · Gin · GORM · PostgreSQL 16 · golang-migrate · JWT · slog · testify + testcontainers · Docker · GitHub Actions

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
| Site (front) | http://localhost:3000 |
| API | http://localhost:8080/health e http://localhost:8080/api/recipes |
| OpenAPI | http://localhost:8080/api/openapi.yaml |
| Postgres | `localhost:5433` (`postgres` / `postgres`) |

Na primeira subida a API aplica as migrations e cria dez receitas e três usuários. Todos entram com a senha `panda-cooking-demo`:

| Usuário | E-mail | Admin |
| --- | --- | --- |
| Chef Maria Silva | `maria@pandacooking.com` | sim |
| João Cozinheiro | `joao@pandacooking.com` | não |
| Ana Paula Gourmet | `ana@pandacooking.com` | não |

API e front recarregam sozinhos ao salvar um arquivo. Se o front estiver em outra pasta, aponte `FRONT_DIR` no `.env`; para subir só banco e API, `docker compose up api`.

> Já tinha subido uma versão antiga? Ela criava o banco sem migrations. Rode `docker compose down -v` uma vez para recriar o banco.

| Comando | O que faz |
| --- | --- |
| `make dev` | sobe banco, API e front |
| `make down` | para tudo (os dados ficam; `docker compose down -v` apaga) |
| `make seed` | apaga o banco e recria os dados de demonstração |
| `make test` | testes de unidade e HTTP com detector de race (precisa de Go 1.26+) |
| `make test-integration` | testes com a API e um Postgres de verdade (precisa de Docker) |
| `make lint` | golangci-lint |

## Exemplo

```bash
# entrar: o access token vem no corpo, o refresh token num cookie HttpOnly
curl -i -c cookies.txt http://localhost:8080/api/auth/login \
  -H 'Content-Type: application/json' \
  -d '{"email":"maria@pandacooking.com","password":"panda-cooking-demo"}'

# buscar receitas
curl 'http://localhost:8080/api/recipes?search=pao&per_page=5'

# renovar a sessão com o cookie
curl -b cookies.txt -c cookies.txt -X POST http://localhost:8080/api/auth/refresh
```

## Documentação

- [PRD](docs/PRD.md): o que o produto faz
- [Arquitetura](docs/ARCHITECTURE.md): camadas, middlewares, modelo de dados e rotas
- [Convenções da API](docs/DESIGN.md): erros, paginação e autenticação
- [Segurança](docs/SECURITY.md): cada risco e como foi tratado
- [Regras](docs/RULES.md): código, segurança, testes e fluxo de git
- [Tarefas](docs/TASKS.md) e [Memória](docs/MEMORY.md): o que foi feito e o porquê de cada decisão
- [Seed](cmd/seed/README.md): dados de demonstração
