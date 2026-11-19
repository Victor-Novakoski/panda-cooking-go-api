# 016. Seed de demonstração dentro da API (`SEED_DEMO=true`)

- **Data:** 2026-10
- **Status:** ativa

## Contexto

Etapa "Tudo com um `docker compose up`". Registrada no antigo `docs/MEMORY.md`.

## Decisão

**Seed de demonstração dentro da API (`SEED_DEMO=true`):** ao subir, cria categorias, três usuários e dezoito receitas se não houver usuário nenhum. Não depende de um container a mais nem de Go na máquina, e não apaga o que já foi cadastrado. `make seed` recria tudo do zero.

## Consequência

O porquê e o efeito foram registrados junto com a decisão, no texto acima.
