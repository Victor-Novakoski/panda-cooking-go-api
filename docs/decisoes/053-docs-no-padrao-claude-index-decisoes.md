# 053. Docs no padrão CLAUDE.md + INDEX + decisões numeradas

- **Data:** 2026-11-12
- **Status:** ativa

## Contexto

Os docs estavam em sete arquivos (PRD, ARCHITECTURE, DESIGN, RULES, SECURITY, TASKS, MEMORY), e MEMORY e TASKS cresciam sem limite.

## Decisão

Adotar o padrão novo: `CLAUDE.md` na raiz e por módulo, `docs/INDEX.md` como roteador, uma decisão por arquivo em `docs/decisoes/` e tarefas como issues do GitHub. PRD, DESIGN, ARCHITECTURE curta e SECURITY ficam para pessoas.

## Consequência

RULES, MEMORY, TASKS e `cmd/seed/README.md` saíram. Decisão nunca é editada: uma nova substitui a antiga.
