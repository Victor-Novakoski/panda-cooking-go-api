# 032. Duas tabelas (`sessions` e `refresh_tokens`)

- **Data:** 2026-10
- **Status:** ativa

## Contexto

Etapa "Sessão em cookie, API completa e revisão geral". Registrada no antigo `docs/MEMORY.md`.

## Decisão

**Duas tabelas (`sessions` e `refresh_tokens`):** revogar a sessão invalida todos os tokens dela de uma vez, e o limite absoluto de 30 dias fica na sessão, não em cada token.

## Consequência

O porquê e o efeito foram registrados junto com a decisão, no texto acima.
