# 033. Cookie `panda_session` separado

- **Data:** 2026-10
- **Status:** ativa

## Contexto

Etapa "Sessão em cookie, API completa e revisão geral". Registrada no antigo `docs/MEMORY.md`.

## Decisão

**Cookie `panda_session` separado:** o refresh token só vai para `/api/auth`, então o servidor do front não o vê nas páginas. O `panda_session` não tem segredo; só diz que existe sessão, para o front redirecionar rota protegida no servidor e não tentar renovar quem nunca entrou.

## Consequência

O porquê e o efeito foram registrados junto com a decisão, no texto acima.
