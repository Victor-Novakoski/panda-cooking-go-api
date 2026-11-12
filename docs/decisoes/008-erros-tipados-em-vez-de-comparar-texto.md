# 008. Erros tipados em vez de comparar texto

- **Data:** 2026-10
- **Status:** ativa

## Contexto

Etapa "Segurança da base". Registrada no antigo `docs/MEMORY.md`.

## Decisão

**Erros tipados em vez de comparar texto:** antes o handler decidia o status comparando `err.Error()` com frases, e o que não batia saía inteiro na resposta. Agora o service devolve `service.Error` com um tipo, e erro desconhecido nunca chega ao cliente. Mudar a frase de um erro não quebra mais o status.

## Consequência

O porquê e o efeito foram registrados junto com a decisão, no texto acima.
