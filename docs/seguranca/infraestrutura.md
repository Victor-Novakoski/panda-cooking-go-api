# Segurança: infraestrutura

## Variáveis de ambiente e segredos — ✅

`.env` no `.gitignore` e no `.dockerignore`; segredos só por variável de ambiente; gitleaks na CI varre todo o histórico a cada PR. A API não sobe sem `SECRET_KEY`, nem com chave de menos de 32 caracteres, e com `APP_ENV=production` recusa a chave do `.env.example` (`internal/config`, com teste).

## Dependências vulneráveis — ✅

govulncheck na CI a cada PR, Trivy na imagem Docker e Dependabot semanal para Go, Actions e Docker.

## Rate limit — ✅

300 requisições por minuto por IP em todo o `/api`, além dos limites de login, cadastro e renovação (ver [autenticação](autenticacao.md)) e de envio de fotos. Responde 429 com `Retry-After`.

## CORS — ✅

Origens lidas de `CORS_ORIGINS`; em produção só `https`, e a API não sobe com origem inválida. Com o front repassando `/api` pela mesma origem, o navegador nem precisa de CORS; a lista serve para chamar a API direto de outro front.

## Headers de segurança — ✅

Em toda resposta: `X-Content-Type-Options`, `X-Frame-Options: DENY`, `Referrer-Policy: no-referrer`, `Content-Security-Policy`, `Cross-Origin-Resource-Policy: same-origin`, `Cache-Control: no-store` (as fotos enviadas usam cache longo, porque o conteúdo de um endereço nunca muda) e, em produção, `Strict-Transport-Security`.

## Timeouts e negação de serviço — ✅

`http.Server` com `ReadHeaderTimeout` (5 s), `ReadTimeout` (15 s), `WriteTimeout` (30 s), `IdleTimeout` e cabeçalhos até 64 KB. Cada requisição tem prazo de 10 segundos, que chega até a consulta no banco pelo `context`. Corpo até 1 MB, listas paginadas (máximo 50 por página), busca com índice e pool de conexões limitado. Ao receber `SIGTERM`, a API para de aceitar conexões e termina as requisições em andamento.

## Banco exposto — ✅

O compose de dev publica Postgres, API e front só em `127.0.0.1`: ninguém na mesma rede (Wi-Fi do café, por exemplo) acessa o banco com a senha padrão. Em produção o banco não publica porta, e `DB_SSLMODE` permite exigir TLS num banco gerenciado.

## Logs e auditoria — ✅

Log estruturado (JSON em produção) com id da requisição em toda linha. 401, 403 e 429 saem em nível `WARN`, 5xx em `ERROR`. Falha de login, bloqueio por e-mail e refresh token reutilizado têm linha própria, sem e-mail em claro.

## IP do cliente falsificado — 🟡

O limite por IP usa o IP da conexão. Só os proxies em `TRUSTED_PROXIES` podem informar o IP do cliente pelo `X-Forwarded-For`. No deploy, quem fica na frente (o servidor do front e o proxy reverso) precisa estar nessa lista, e o proxy reverso precisa **sobrescrever** o `X-Forwarded-For` que vem da internet (o Caddy faz isso por padrão); senão, um cliente troca de "IP" a cada requisição.
