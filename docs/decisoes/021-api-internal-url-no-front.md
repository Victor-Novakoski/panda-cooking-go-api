# 021. `API_INTERNAL_URL` no front

- **Data:** 2026-10
- **Status:** substituída pela decisão 027 do front (rota `/api/*` no Next, com `API_URL` só no servidor)

## Contexto

Etapa "Tudo com um `docker compose up`". Registrada no antigo `docs/MEMORY.md`.

## Decisão

**`API_INTERNAL_URL` no front:** a página da receita é renderizada no servidor do Next, que dentro do Docker não alcança `localhost:8080`; ele usa `http://api:8080` e o navegador continua com a URL pública. (Substituído na etapa 5 pelo `API_URL`, ver abaixo.)

## Consequência

Substituída na etapa 5: o front passou a usar `API_URL`, só no servidor.
