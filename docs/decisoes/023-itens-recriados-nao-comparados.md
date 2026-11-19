# 023. Itens recriados, não comparados

- **Data:** 2026-10
- **Status:** ativa

## Contexto

Etapa "Telas de edição". Registrada no antigo `docs/MEMORY.md`.

## Decisão

**Itens recriados, não comparados:** no `PUT` os itens antigos são apagados e os novos criados na ordem recebida. Comparar item a item não traz nada para quem usa a tela e complicaria o código; o custo é o id do item mudar, e nada guarda esse id.

## Consequência

O porquê e o efeito foram registrados junto com a decisão, no texto acima.
