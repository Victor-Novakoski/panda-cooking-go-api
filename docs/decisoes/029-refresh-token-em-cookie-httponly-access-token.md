# 029. Refresh token em cookie `HttpOnly` + access token de 15 min, não JWT de 24h no `localStorage`

- **Data:** 2026-10
- **Status:** ativa

## Contexto

Etapa "Sessão em cookie, API completa e revisão geral". Registrada no antigo `docs/MEMORY.md`.

## Decisão

**Refresh token em cookie `HttpOnly` + access token de 15 min, não JWT de 24h no `localStorage`:** com o token legível por JavaScript, um XSS levava a sessão por 24h sem como revogar. O JWT curto continua sem consulta ao banco em cada requisição; quem segura a sessão é o refresh token, que dá para revogar.

## Consequência

O porquê e o efeito foram registrados junto com a decisão, no texto acima.
