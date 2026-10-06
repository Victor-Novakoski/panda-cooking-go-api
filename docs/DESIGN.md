# Convenções da API

- JSON com campos em `snake_case`.
- Erro: `{"error": "mensagem em português"}`.
- Status: 200 leitura/edição, 201 criação, 204 remoção, 400 corpo inválido, 401 sem token ou token inválido, 403 sem permissão, 404 não encontrado, 409 conflito, 429 muitas tentativas (com `Retry-After` quando é limite por IP), 500 erro inesperado.
- 500 sempre responde `{"error": "erro interno, tente novamente mais tarde"}`; o detalhe fica só no log.
- Item de uma receita acessado pela rota de outra receita: 404, como se não existisse.
- Token no cabeçalho `Authorization: Bearer <token>`, válido por 24h.
- Recurso de outra pessoa: 403 (ver [SECURITY.md](SECURITY.md#6-idor)).

Pendências de padronização estão em [TASKS.md](TASKS.md) (ex.: 422 com erro por campo, paginação).
