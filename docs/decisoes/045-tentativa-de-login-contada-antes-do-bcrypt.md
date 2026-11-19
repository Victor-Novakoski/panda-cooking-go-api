# 045. Tentativa de login contada antes do bcrypt

- **Data:** 2026-10
- **Status:** ativa

## Contexto

Etapa "Sessão em cookie, API completa e revisão geral". Registrada no antigo `docs/MEMORY.md`.

## Decisão

**Tentativa de login contada antes do bcrypt:** conferir o bloqueio antes e contar depois deixava passar várias senhas ao mesmo tempo (8 em vez de 5 num teste com 10 simultâneas). O login certo zera a contagem, então quem acerta não perde nada.

## Consequência

O porquê e o efeito foram registrados junto com a decisão, no texto acima.
