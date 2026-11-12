# 022. `PUT /recipes/:id` com a receita inteira

- **Data:** 2026-10
- **Status:** ativa

## Contexto

Etapa "Telas de edição". Registrada no antigo `docs/MEMORY.md`.

## Decisão

**`PUT /recipes/:id` com a receita inteira:** a tela de edição manda tudo de uma vez e a API troca dados, fotos, ingredientes e passos numa transação. Com as rotas por item o front teria que fazer várias chamadas, e uma falha no meio deixaria a receita pela metade. As rotas por item continuam para quem quiser mexer em uma parte só.

## Consequência

O porquê e o efeito foram registrados junto com a decisão, no texto acima.
