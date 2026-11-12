# 044. `healthcheck` no próprio binário

- **Data:** 2026-10
- **Status:** ativa

## Contexto

Etapa "Sessão em cookie, API completa e revisão geral". Registrada no antigo `docs/MEMORY.md`.

## Decisão

**`healthcheck` no próprio binário:** a imagem distroless não tem `curl` nem shell; `api healthcheck` faz a chamada ao `/health` e serve de `HEALTHCHECK` do Docker.

## Consequência

O porquê e o efeito foram registrados junto com a decisão, no texto acima.
