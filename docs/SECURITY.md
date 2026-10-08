# Segurança

Como a API trata cada risco. As tarefas estão em [TASKS.md](TASKS.md) e o checklist rápido para cada mudança está em [RULES.md](RULES.md#3-segurança). Os riscos que dependem do front estão no [SECURITY.md do front](https://github.com/Victor-Novakoski/panda-cooking-front/blob/develop/docs/SECURITY.md).

Situação revisada em 07/10/2026, na etapa 5 (sessão em cookie, API completa e revisão geral).

**Legenda:** ✅ feito · 🟡 parcial · ⚪ ainda não se aplica

## Resumo

| # | Risco | Situação | O que é |
| --- | --- | --- | --- |
| 1 | Variáveis de ambiente expostas | ✅ | Segredo vazando por arquivo commitado, imagem Docker ou valor padrão usado em produção. |
| 2 | Validação no front-end | ✅ | Dar retorno rápido no formulário. Ajuda o usuário, mas não protege nada. |
| 3 | Validação no back-end | ✅ | A API conferir tipo, formato e tamanho de tudo que recebe. |
| 4 | SQL Injection | ✅ | Texto do usuário virando parte do comando SQL. |
| 5 | Autenticação fraca | ✅ | Senha fraca, login que dá pistas, credencial padrão. |
| 6 | IDOR | ✅ | Trocar um id na URL e mexer em dado de outra pessoa. |
| 7 | Senhas no banco | ✅ | Guardar a senha em texto puro. |
| 8 | Força bruta | ✅ | Tentar milhares de senhas seguidas. |
| 9 | Envio duplicado | ✅ | Clique duplo criando o mesmo registro duas vezes. |
| 10 | CSRF | ✅ | Outro site fazendo o navegador logado enviar uma ação. |
| 11 | Upload sem validação | ⚪ | Arquivo malicioso ou grande demais. Hoje imagem é só URL. |
| 12 | Vazamento de informação | ✅ | Erro interno ou detalhe do banco aparecendo na resposta. |
| 13 | Dependências vulneráveis | ✅ | Biblioteca com falha conhecida. |
| 14 | Tokens | ✅ | Token longo demais, sem revogação, guardado onde script lê. |
| 15 | Rate limit | ✅ | Limitar requisições por cliente. |
| 16 | Dados sensíveis expostos | ✅ | Resposta ou log mostrando dado pessoal além do necessário. |
| 17 | CORS | ✅ | Quais sites podem chamar a API. |
| 18 | XSS | ✅ | Script injetado num dado rodando no navegador de outra pessoa. |
| 19 | Headers de segurança | ✅ | Cabeçalhos que mandam o navegador se proteger. |
| 20 | Timeouts e negação de serviço | ✅ | Conexão lenta ou corpo enorme prendendo o servidor. |
| 21 | Banco exposto | ✅ | Banco acessível pela rede sem passar pela API. |
| 22 | Logs e auditoria | ✅ | Registrar eventos suspeitos para investigar. |
| 23 | IP do cliente falsificado | 🟡 | `X-Forwarded-For` falso burlando o limite por IP. Depende do deploy. |

---

## 1. Variáveis de ambiente expostas — ✅

`.env` no `.gitignore` e no `.dockerignore`; segredos só por variável de ambiente; gitleaks na CI varre todo o histórico a cada PR. A API não sobe sem `SECRET_KEY`, nem com chave de menos de 32 caracteres, e com `APP_ENV=production` recusa a chave do `.env.example` (`internal/config`, com teste).

## 2. Validação no front-end — ✅

Os formulários do front usam Zod com os mesmos limites da API (`internal/service/inputs.go` ↔ `src/lib/limits.ts` do front) e mostram no campo certo os erros 422 que a API devolve.

## 3. Validação no back-end — ✅

Cada entrada tem limite de tamanho em todos os textos e listas (até 10 fotos, 50 ingredientes e 50 passos), trim e e-mail em minúsculas antes de validar, URLs só `https`, porções de 1 a 100. Campo desconhecido no JSON é recusado, e o corpo é limitado a 1 MB (413). Erro de validação responde 422 com a mensagem de cada campo, no caminho do JSON (`ingredients[2].amount`), sem nome interno de struct. Parâmetros de busca e paginação também são validados (`per_page` até 50, `search` até 100 caracteres).

## 4. SQL Injection — ✅

Todo acesso ao banco passa pelo GORM com parâmetros, inclusive o termo da busca, que entra no `LIKE` com `%` e `_` escapados. Regra em [RULES.md](RULES.md): nunca montar SQL concatenando string.

## 5. Autenticação fraca — ✅

Senha com 10 a 72 bytes (limite do bcrypt, conferido antes de chegar nele), recusando as mais comuns, as de um caractere só e a igual ao e-mail. Mesma mensagem para e-mail inexistente e senha errada, e o mesmo tempo de resposta: e-mail inexistente também passa pelo bcrypt. Os usuários do seed têm senha pública de propósito (é um portfólio); com `SEED_DEMO=false` eles não são criados.

## 6. IDOR — ✅

Editar ou apagar receita, comentário e itens da receita confere se a receita é de quem pede (403). Item de outra receita responde 404, como se não existisse, e nada é gravado (`TestRecipeService_ItemOfAnotherRecipe`). Favoritos e perfil usam sempre o usuário do token, nunca um id da URL. Os testes de integração repetem essas tentativas contra a API e o banco de verdade.

## 7. Senhas no banco — ✅

Só o hash bcrypt é salvo (`password_hash`). As respostas usam `UserResponse`, que não tem o campo.

## 8. Força bruta — ✅

`POST /api/auth/login` aceita 10 requisições por minuto por IP (429 com `Retry-After`). Além disso, 5 senhas erradas para o mesmo e-mail bloqueiam o login dele por 15 minutos, mesmo com a senha certa e vindo de outro IP; e-mail inexistente conta igual, para o bloqueio não revelar quem tem conta. A tentativa é contada antes do bcrypt, então requisições simultâneas não passam do limite, e o 429 do bloqueio também traz `Retry-After`. Cadastro: 10 por hora por IP; renovação de sessão: 30 por minuto. Os contadores ficam em memória (`internal/ratelimit`), o que serve para uma instância só.

## 9. Envio duplicado — ✅

O front desabilita o botão enquanto a requisição não volta. Na API, favoritar de novo responde 409 (índice único), e repetir a renovação de sessão dentro de 20 segundos é aceito sem derrubar a sessão (duas abas renovando ao mesmo tempo).

## 10. CSRF — ✅

- Rotas de dados exigem `Authorization: Bearer`, que outro site não consegue mandar.
- As rotas que usam cookie (`/api/auth/refresh` e `/logout`) recebem o refresh token com `SameSite=Strict`, só em `/api/auth`.
- Toda requisição que muda dados passa pelo `http.CrossOriginProtection` do Go: vinda de outra origem (pelo `Sec-Fetch-Site` ou `Origin`), responde 403, a não ser que a origem esteja em `CORS_ORIGINS`.

## 11. Upload sem validação — ⚪

Hoje a imagem é uma URL `https` com até 2048 caracteres. Se virar upload: limite de tamanho, tipo conferido pelo conteúdo e nome gerado pelo servidor.

## 12. Vazamento de informação — ✅

O service devolve erros tipados (`service.Error` com um `Kind`) e o handler escolhe o status por eles (`respondError`). Qualquer outro erro vira 500 com "erro interno, tente novamente mais tarde" e o detalhe vai só para o log, com o `X-Request-Id`. E-mail duplicado responde 409; id que não é UUID responde 404 em vez do erro de sintaxe do Postgres. Gin em modo release com `APP_ENV=production`; panic vira 500 sem stack trace na resposta.

## 13. Dependências vulneráveis — ✅

govulncheck na CI a cada PR, Trivy na imagem Docker e Dependabot semanal para Go, Actions e Docker.

## 14. Tokens — ✅

- **Access token:** JWT HS256 de 15 minutos, algoritmo conferido na validação, mandado no corpo da resposta. O front guarda só em memória.
- **Refresh token:** 32 bytes aleatórios num cookie `HttpOnly` (`SameSite=Strict`, só em `/api/auth`, `Secure` em produção). No banco fica só o hash SHA-256: quem lê o banco não consegue usar o token.
- **Rotação:** cada renovação troca o refresh token. Token já usado, apresentado de novo depois de 20 segundos, é sinal de cópia: a sessão inteira é revogada e o evento vai para o log.
- **Validade:** cada refresh token vale 7 dias; a sessão, 30 dias no máximo, depois disso é preciso entrar de novo.
- **Revogação:** logout revoga a sessão; apagar a conta apaga as sessões em cascata. O access token já emitido vale até vencer (no máximo 15 minutos).

## 15. Rate limit — ✅

300 requisições por minuto por IP em todo o `/api`, além dos limites de login, cadastro e renovação (item 8). Responde 429 com `Retry-After`.

## 16. Dados sensíveis expostos — ✅

Receitas e comentários mostram do autor só id, nome e foto (sem e-mail nem `is_adm`). O e-mail só aparece no perfil do próprio usuário. No log, a falha de login identifica o e-mail por um HMAC com a `SECRET_KEY`, nunca em claro.

## 17. CORS — ✅

Origens lidas de `CORS_ORIGINS`; em produção só `https`, e a API não sobe com origem inválida. Com o front repassando `/api` pela mesma origem, o navegador nem precisa de CORS; a lista serve para chamar a API direto de outro front.

## 18. XSS — ✅

A API só devolve JSON, com `X-Content-Type-Options: nosniff` e `Content-Security-Policy: default-src 'none'`, então uma resposta aberta direto no navegador não executa nada. URLs de imagem só `https`. No front, o React escapa o texto e há CSP com nonce.

## 19. Headers de segurança — ✅

Em toda resposta: `X-Content-Type-Options`, `X-Frame-Options: DENY`, `Referrer-Policy: no-referrer`, `Content-Security-Policy`, `Cross-Origin-Resource-Policy: same-origin`, `Cache-Control: no-store` e, em produção, `Strict-Transport-Security`.

## 20. Timeouts e negação de serviço — ✅

`http.Server` com `ReadHeaderTimeout` (5 s), `ReadTimeout` (15 s), `WriteTimeout` (30 s), `IdleTimeout` e cabeçalhos até 64 KB. Cada requisição tem prazo de 10 segundos, que chega até a consulta no banco pelo `context`. Corpo até 1 MB, listas paginadas (máximo 50 por página), busca com índice e pool de conexões limitado. Ao receber `SIGTERM`, a API para de aceitar conexões e termina as requisições em andamento.

## 21. Banco exposto — ✅

O compose de dev publica Postgres, API e front só em `127.0.0.1`: ninguém na mesma rede (Wi-Fi do café, por exemplo) acessa o banco com a senha padrão. Em produção o banco não publica porta, e `DB_SSLMODE` permite exigir TLS num banco gerenciado.

## 22. Logs e auditoria — ✅

Log estruturado (JSON em produção) com id da requisição em toda linha. 401, 403 e 429 saem em nível `WARN`, 5xx em `ERROR`. Falha de login, bloqueio por e-mail e refresh token reutilizado têm linha própria, sem e-mail em claro.

## 23. IP do cliente falsificado — 🟡

O limite por IP usa o IP da conexão. Só os proxies em `TRUSTED_PROXIES` podem informar o IP do cliente pelo `X-Forwarded-For`. No deploy, quem fica na frente (o servidor do front e o proxy reverso) precisa estar nessa lista, e o proxy reverso precisa **sobrescrever** o `X-Forwarded-For` que vem da internet (o Caddy faz isso por padrão); senão, um cliente troca de "IP" a cada requisição.
