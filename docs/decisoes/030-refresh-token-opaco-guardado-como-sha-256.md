# 030. Refresh token opaco guardado como SHA-256, não JWT nem bcrypt

- **Data:** 2026-10
- **Status:** ativa

## Contexto

Etapa "Sessão em cookie, API completa e revisão geral". Registrada no antigo `docs/MEMORY.md`.

## Decisão

**Refresh token opaco guardado como SHA-256, não JWT nem bcrypt:** precisa ser revogável (então vive no banco), e quem lê o banco não deve conseguir usá-lo. Como é aleatório e longo, SHA-256 basta; bcrypt existe para senha fraca.

## Consequência

O porquê e o efeito foram registrados junto com a decisão, no texto acima.
