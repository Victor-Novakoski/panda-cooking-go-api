# 037. Busca com `unaccent` + `pg_trgm`

- **Data:** 2026-10
- **Status:** ativa

## Contexto

Etapa "Sessão em cookie, API completa e revisão geral". Registrada no antigo `docs/MEMORY.md`.

## Decisão

**Busca com `unaccent` + `pg_trgm`:** "pao" acha "Pão", e o `LIKE '%termo%'` usa índice GIN. Full-text do Postgres foi descartado: corta palavras pela raiz e não acha pedaço de palavra, o que confunde numa busca de receita.

## Consequência

O porquê e o efeito foram registrados junto com a decisão, no texto acima.
