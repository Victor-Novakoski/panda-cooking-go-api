# 024. `Omit(clause.Associations)` no `Save`

- **Data:** 2026-10
- **Status:** ativa

## Contexto

Etapa "Telas de edição". Registrada no antigo `docs/MEMORY.md`.

## Decisão

**`Omit(clause.Associations)` no `Save`:** com a `Category` carregada pelo `Preload`, o `Save` do GORM gravava a associação e voltava o `category_id` para o antigo, e trocar a categoria no `PATCH` não tinha efeito.

## Consequência

O porquê e o efeito foram registrados junto com a decisão, no texto acima.
