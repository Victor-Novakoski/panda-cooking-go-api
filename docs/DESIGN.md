# Convenções da API

- JSON com campos em `snake_case`.
- Erro: `{"error": "mensagem em português"}`.
- Status: 200 leitura/edição, 201 criação, 204 remoção, 400 corpo inválido, 401 sem token ou token inválido, 403 sem permissão, 404 não encontrado, 409 conflito, 500 erro inesperado.
- Token no cabeçalho `Authorization: Bearer <token>`, válido por 24h.
- Recurso de outra pessoa: 403 (ver [SECURITY.md](SECURITY.md#6-idor)).

Pendências de padronização estão em [TASKS.md](TASKS.md) (ex.: 422 com erro por campo, paginação).
