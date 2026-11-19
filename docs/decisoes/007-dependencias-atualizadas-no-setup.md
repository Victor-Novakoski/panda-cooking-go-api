# 007. Dependências atualizadas no setup

- **Data:** 2026-10
- **Status:** ativa

## Contexto

Etapa "Setup". Registrada no antigo `docs/MEMORY.md`.

## Decisão

**Dependências atualizadas no setup:** pgx 5.6 tinha falha crítica (e risco de SQL injection) e x/crypto, x/net e x/text tinham falhas altas; o govulncheck e o Trivy barraram na primeira CI.

## Consequência

O porquê e o efeito foram registrados junto com a decisão, no texto acima.
