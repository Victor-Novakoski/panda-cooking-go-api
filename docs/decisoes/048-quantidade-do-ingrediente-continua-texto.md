# 048. Quantidade do ingrediente continua texto

- **Data:** 2026-10-13
- **Status:** ativa

## Contexto

Lista de melhorias levantada navegando pelo site (etapa 6, 13/10/2026). Registrada no antigo `docs/TASKS.md`.

## Decisão

Não mexer no banco para guardar a quantidade como número + unidade. O front lê o número do começo do texto ("500g", "3 colheres de sopa", "1/2 xícara") e recalcula; "a gosto" fica como está.

## Consequência

Ajustar porções funciona sem migration. Quantidade que não começa com número não é recalculada.
