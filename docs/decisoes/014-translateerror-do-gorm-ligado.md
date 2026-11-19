# 014. `TranslateError` do GORM ligado

- **Data:** 2026-10
- **Status:** ativa

## Contexto

Etapa "Segurança da base". Registrada no antigo `docs/MEMORY.md`.

## Decisão

**`TranslateError` do GORM ligado:** transforma a violação de unicidade do Postgres em `gorm.ErrDuplicatedKey`, que o service converte em 409 sem depender do texto do erro do banco.

## Consequência

O porquê e o efeito foram registrados junto com a decisão, no texto acima.
