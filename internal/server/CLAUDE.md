# internal/server

Middlewares globais, a tabela de rotas e a ligação dos handlers (`server.go`).

## Regras

- **Ordem dos middlewares:** recuperação de panic → id da requisição → log de acesso → headers de segurança → `http.CrossOriginProtection` → CORS → limite de corpo (1 MB) → prazo (10 s) → limite por IP → JWT nas rotas com login. O `CrossOriginProtection` fica antes do CORS para a recusa sair no formato de erro da API.
- Tudo sob `/api`, menos o `/health`.
- Rota nova ou alterada entra no `api/openapi.yaml` no mesmo commit: um teste falha se uma rota do roteador não estiver documentada.
- Rota nova declara a autenticação explicitamente; rota pública é exceção justificada.
- Limites por IP ficam nas constantes do topo de `server.go` (global 300/min, login 10/min, renovação 30/min, cadastro 10/h, fotos 30/min). Os contadores são em memória (`internal/ratelimit`): servem para uma instância só.

## Armadilhas conhecidas

- **`SetTrustedProxies`:** o padrão do Gin confia em todo proxy, e o cliente trocaria de "IP" a cada requisição mandando `X-Forwarded-For`. Só os endereços de `TRUSTED_PROXIES` são aceitos; em produção, o proxy reverso precisa sobrescrever o header.
- `CORS_ORIGINS` em produção só aceita `https`; sem ele, a API cai em `http://localhost:3000` e não sobe com `APP_ENV=production`.
- A imagem é distroless (sem shell nem `curl`): o `HEALTHCHECK` do Docker é o subcomando `api healthcheck`.
- O `/health` responde 503 com o banco fora do ar.

## Ver também

`docs/decisoes/` 010, 012, 034, 035, 042, 044.
