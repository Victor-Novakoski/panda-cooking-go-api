# 025. Ordem por id nos `Preload`

- **Data:** 2026-10
- **Status:** ativa

## Contexto

Etapa "Telas de edição". Registrada no antigo `docs/MEMORY.md`.

## Decisão

**Ordem por id nos `Preload`:** sem `ORDER BY` o Postgres não garante a ordem dos passos do preparo.

## Consequência

O porquê e o efeito foram registrados junto com a decisão, no texto acima.
