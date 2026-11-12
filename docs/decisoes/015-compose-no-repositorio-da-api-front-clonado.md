# 015. Compose no repositório da API, front clonado ao lado

- **Data:** 2026-10
- **Status:** ativa

## Contexto

Etapa "Tudo com um `docker compose up`". Registrada no antigo `docs/MEMORY.md`.

## Decisão

**Compose no repositório da API, front clonado ao lado:** banco, migração e seed são da API, então ela é quem sabe subir o ambiente. O front entra por bind mount de `../panda-cooking-front` (ou `FRONT_DIR`), sem imagem publicada nem submódulo. Um terceiro repositório só para o compose seria mais uma coisa para manter sincronizada.

## Consequência

O porquê e o efeito foram registrados junto com a decisão, no texto acima.
