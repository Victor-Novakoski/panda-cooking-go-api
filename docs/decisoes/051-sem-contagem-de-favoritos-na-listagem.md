# 051. Sem contagem de favoritos na listagem

- **Data:** 2026-10-13
- **Status:** ativa

## Contexto

Segunda leva de sugestões (etapa 7, 13/10/2026). Registrada no antigo `docs/TASKS.md`.

## Decisão

A listagem de receitas não traz a contagem de favoritos: vira placar e influencia quem vê, o mesmo motivo de não ter nota (049).

## Consequência

A contagem só é usada para ordenar (`sort=favorites`), nunca devolvida.
