# 018. Categorias por nome no seed

- **Data:** 2026-10
- **Status:** ativa

## Contexto

Etapa "Tudo com um `docker compose up`". Registrada no antigo `docs/MEMORY.md`.

## Decisão

**Categorias por nome no seed:** antes as receitas apontavam para `CategoryID: 1..12` e dependiam do `RESTART IDENTITY`; se a API tivesse criado as categorias padrão antes, as receitas iam para a categoria errada.

## Consequência

O porquê e o efeito foram registrados junto com a decisão, no texto acima.
