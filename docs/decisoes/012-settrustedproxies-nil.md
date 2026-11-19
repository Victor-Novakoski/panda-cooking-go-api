# 012. `SetTrustedProxies(nil)`

- **Data:** 2026-10
- **Status:** ativa

## Contexto

Etapa "Segurança da base". Registrada no antigo `docs/MEMORY.md`.

## Decisão

**`SetTrustedProxies(nil)`:** com o padrão do Gin (confia em todo proxy), o cliente trocaria de "IP" a cada requisição mandando `X-Forwarded-For`.

## Consequência

O porquê e o efeito foram registrados junto com a decisão, no texto acima.
