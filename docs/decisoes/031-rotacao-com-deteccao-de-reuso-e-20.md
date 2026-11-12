# 031. Rotação com detecção de reúso e 20 s de tolerância

- **Data:** 2026-10
- **Status:** ativa

## Contexto

Etapa "Sessão em cookie, API completa e revisão geral". Registrada no antigo `docs/MEMORY.md`.

## Decisão

**Rotação com detecção de reúso e 20 s de tolerância:** token usado de novo indica cópia e derruba a sessão. A tolerância cobre duas abas renovando juntas e a resposta que não chegou ao navegador, que derrubariam a sessão de quem não fez nada de errado. A troca é atômica (`SELECT ... FOR UPDATE` na transação).

## Consequência

O porquê e o efeito foram registrados junto com a decisão, no texto acima.
