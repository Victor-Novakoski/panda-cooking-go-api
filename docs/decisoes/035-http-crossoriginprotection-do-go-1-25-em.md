# 035. `http.CrossOriginProtection` do Go 1.25+ em vez de token CSRF

- **Data:** 2026-10
- **Status:** ativa

## Contexto

Etapa "Sessão em cookie, API completa e revisão geral". Registrada no antigo `docs/MEMORY.md`.

## Decisão

**`http.CrossOriginProtection` do Go 1.25+ em vez de token CSRF:** usa `Sec-Fetch-Site`/`Origin`, que todo navegador atual manda, sem estado e sem mudar o front. Fica antes do CORS para a recusa sair no formato de erro da API.

## Consequência

O porquê e o efeito foram registrados junto com a decisão, no texto acima.
