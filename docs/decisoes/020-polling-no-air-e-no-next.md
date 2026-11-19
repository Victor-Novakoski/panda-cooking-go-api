# 020. Polling no air e no Next

- **Data:** 2026-10
- **Status:** ativa

## Contexto

Etapa "Tudo com um `docker compose up`". Registrada no antigo `docs/MEMORY.md`.

## Decisão

**Polling no air e no Next:** no Docker do Windows e do Mac a pasta montada não avisa o container de arquivo alterado; sem polling o hot reload não funciona.

## Consequência

O porquê e o efeito foram registrados junto com a decisão, no texto acima.
