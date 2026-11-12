# 034. Rotas sob `/api`

- **Data:** 2026-10
- **Status:** ativa

## Contexto

Etapa "Sessão em cookie, API completa e revisão geral". Registrada no antigo `docs/MEMORY.md`.

## Decisão

**Rotas sob `/api`:** o front repassa `/api/*` para cá pela mesma origem; com o prefixo, o proxy reverso também consegue mandar `/api` direto para a API, sem conflito com as páginas.

## Consequência

O porquê e o efeito foram registrados junto com a decisão, no texto acima.
