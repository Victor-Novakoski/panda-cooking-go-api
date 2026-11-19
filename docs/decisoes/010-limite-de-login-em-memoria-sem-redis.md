# 010. Limite de login em memória, sem Redis

- **Data:** 2026-10
- **Status:** ativa

## Contexto

Etapa "Segurança da base". Registrada no antigo `docs/MEMORY.md`.

## Decisão

**Limite de login em memória, sem Redis:** a API roda numa instância só; o Rastreia usa Redis porque tem WebSocket e várias réplicas. Se escalar, o contador vai para um armazenamento compartilhado.

## Consequência

O porquê e o efeito foram registrados junto com a decisão, no texto acima.
