# 046. Renovação que falha no servidor não apaga os cookies

- **Data:** 2026-10
- **Status:** ativa

## Contexto

Etapa "Sessão em cookie, API completa e revisão geral". Registrada no antigo `docs/MEMORY.md`.

## Decisão

**Renovação que falha no servidor não apaga os cookies:** com o banco fora do ar a troca é desfeita e a sessão continua valendo; apagar o cookie deslogaria todo mundo que renovasse durante a queda. Só a sessão encerrada (401) apaga.

## Consequência

O porquê e o efeito foram registrados junto com a decisão, no texto acima.
