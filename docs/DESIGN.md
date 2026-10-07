# Convenções da API

Especificação completa (rotas, corpos, respostas e erros) em [`api/openapi.yaml`](../api/openapi.yaml), servida em `GET /api/openapi.yaml`.

## Formato

- Tudo sob `/api`, menos o `/health`.
- JSON com campos em `snake_case`. Texto recebido passa por trim; e-mail vira minúsculo.
- Corpo com campo desconhecido é recusado (400), para erro de digitação no cliente não passar em silêncio. Corpo maior que 1 MB: 413.

## Erros

- Sempre `{"error": "mensagem em português"}`, sem detalhe interno.
- Validação: **422** com o erro de cada campo, no caminho do campo no JSON:

  ```json
  {"error": "dados inválidos", "fields": {"name": "campo obrigatório", "ingredients[0].amount": "use no máximo 60 caracteres"}}
  ```

- E-mail já cadastrado: **409** no mesmo formato, com `fields.email`, para o front mostrar no campo.
- 500 sempre responde `{"error": "erro interno, tente novamente mais tarde"}`; o detalhe fica só no log, com o `X-Request-Id` da requisição.

| Status | Quando |
| --- | --- |
| 200 | leitura e edição |
| 201 | criação |
| 204 | remoção, logout |
| 400 | JSON malformado, campo desconhecido, regra de negócio (ex.: categoria inexistente) |
| 401 | sem token, token inválido ou vencido, sessão encerrada |
| 403 | recurso de outra pessoa; requisição de outra origem |
| 404 | não existe (inclusive id que não é UUID e item de outra receita) |
| 409 | conflito (e-mail já cadastrado) |
| 413 | corpo grande demais |
| 422 | validação, com `fields` |
| 429 | limite de requisições, com `Retry-After` |
| 503 | requisição que passou de 10 segundos; `/health` com o banco fora do ar |

## Listas

- Paginadas com `page` (padrão 1) e `per_page` (padrão 12 nas receitas e 10 nos comentários, máximo 50):

  ```json
  {"items": [...], "page": 1, "per_page": 12, "total": 37, "total_pages": 4}
  ```

- Mais novos primeiro. `GET /api/recipes` aceita `search` (nome e descrição, sem diferenciar acento nem maiúscula), `category_id` e `user_id`, combináveis.
- Página depois da última volta `items` vazio com o `total` certo.

## Autenticação

- `POST /api/auth/login` devolve o access token (JWT, 15 minutos) no corpo e grava dois cookies `HttpOnly`:
  - `panda_refresh`: refresh token, `SameSite=Strict`, só enviado para `/api/auth`;
  - `panda_session`: só indica que existe sessão (sem segredo), `SameSite=Lax`, para o front saber disso no servidor.
- As outras rotas recebem o access token em `Authorization: Bearer <token>`.
- `POST /api/auth/refresh` troca o refresh token por um novo e devolve um access token novo. `POST /api/auth/logout` encerra a sessão e apaga os cookies.
- Em produção os cookies têm `Secure` (só HTTPS).
- Recurso de outra pessoa: 403 (ver [SECURITY.md](SECURITY.md#6-idor)).
