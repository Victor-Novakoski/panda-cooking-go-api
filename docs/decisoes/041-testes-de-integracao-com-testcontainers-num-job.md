# 041. Testes de integração com testcontainers, num job próprio

- **Data:** 2026-10
- **Status:** ativa

## Contexto

Etapa "Sessão em cookie, API completa e revisão geral". Registrada no antigo `docs/MEMORY.md`.

## Decisão

**Testes de integração com testcontainers, num job próprio:** os testes de unidade usam mocks e não pegam erro de SQL, de transação ou de migration. Rodam com a build tag `integration`, para o `make test` continuar rápido e sem Docker.

## Consequência

O porquê e o efeito foram registrados junto com a decisão, no texto acima.
