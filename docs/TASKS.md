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

## Etapa 6 — Melhorias de uso

Lista levantada navegando pelo site (13/10/2026). Aqui fica a parte da API; as telas e o que é só do front estão em [TASKS.md do front](https://github.com/Victor-Novakoski/panda-cooking-front/blob/develop/docs/TASKS.md).

- [x] Tempo da receita em minutos (`time_minutes`, 1 a 2880) no lugar do texto livre; a migration 000005 converte o que já existe ("40", "20 minutos", "1h 30min", "3 horas") e o que não dá para entender vira 30
- [x] ~~Quantidade do ingrediente como número + unidade~~: não precisa mexer no banco. O front lê o número do começo do texto ("500g", "3 colheres de sopa", "1/2 xícara") e recalcula; "a gosto" fica como está
- [x] Ordenação em `GET /api/recipes` (`sort=recent`, `quickest` ou `favorites`) e filtro por tempo máximo (`max_time`)
- [x] ~~Avaliação com nota~~: decidido não fazer. Numa comunidade, nota baixa influencia quem vê a receita e desanima quem publicou; ficam os favoritos e os comentários
- [x] Vídeo por link do YouTube: `video_url` opcional na receita, aceitando só endereço do YouTube (enviar arquivo de vídeo fica de fora: os uploads moram no banco, com limite de 900 KB)
- [x] ~~Categorias pelo admin~~: decidido não fazer. As categorias continuam fixas, criadas pela migration
- [x] Perfil público do autor: `GET /api/users/{id}` com nome, foto e data de entrada, sem e-mail (as receitas já saem por `GET /api/recipes?user_id=`)
- [ ] Recuperar senha: token de uso único com validade curta, enviado por e-mail; depende de escolher um serviço de envio e fica para depois do deploy

## Etapa 7 — Próximas melhorias

Segunda leva de sugestões (13/10/2026), ainda sem nada começado. Aqui fica a parte da API; o que é só do front está em [TASKS.md do front](https://github.com/Victor-Novakoski/panda-cooking-front/blob/develop/docs/TASKS.md). Continuam valendo o deploy (etapa 4) e o recuperar senha (etapa 6).

- [ ] Buscar por ingrediente: filtro em `GET /api/recipes` pelas receitas que têm os ingredientes pedidos (os nomes já são guardados normalizados na tabela `ingredients`)
- [ ] Moderação: denúncia de receita e de comentário (quem denunciou, motivo, quando) e o admin podendo apagar receita imprópria; hoje o admin só apaga comentário
- [ ] Limpeza de fotos sem uso: as fotos enviadas ficam no banco (`uploads`), e a que foi enviada e nunca usada numa receita ou num perfil fica lá para sempre; apagar as órfãs com mais de um dia

**Decidido não fazer**

- Contagem de favoritos na listagem de receitas: vira placar e influencia quem vê, o mesmo motivo de não ter nota
