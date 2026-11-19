# Segurança: autenticação

## Autenticação fraca — ✅

Senha com 10 a 72 bytes (limite do bcrypt, conferido antes de chegar nele), recusando as mais comuns, as de um caractere só e a igual ao e-mail. Mesma mensagem para e-mail inexistente e senha errada, e o mesmo tempo de resposta: e-mail inexistente também passa pelo bcrypt. Os usuários do seed têm senha pública de propósito (é um portfólio), definida em `internal/seed/seed.go`; com `SEED_DEMO=false` eles não são criados.

## Senhas no banco — ✅

Só o hash bcrypt é salvo (`password_hash`). As respostas usam `UserResponse`, que não tem o campo.

## Força bruta — ✅

`POST /api/auth/login` aceita 10 requisições por minuto por IP (429 com `Retry-After`). Além disso, 5 senhas erradas para o mesmo e-mail bloqueiam o login dele por 15 minutos, mesmo com a senha certa e vindo de outro IP; e-mail inexistente conta igual, para o bloqueio não revelar quem tem conta. A tentativa é contada antes do bcrypt, então requisições simultâneas não passam do limite, e o 429 do bloqueio também traz `Retry-After`. Cadastro: 10 por hora por IP; renovação de sessão: 30 por minuto. Os contadores ficam em memória (`internal/ratelimit`), o que serve para uma instância só.

## Tokens — ✅

- **Access token:** JWT HS256 de 15 minutos, algoritmo conferido na validação, mandado no corpo da resposta. O front guarda só em memória.
- **Refresh token:** 32 bytes aleatórios num cookie `HttpOnly` (`SameSite=Strict`, só em `/api/auth`, `Secure` em produção). No banco fica só o hash SHA-256: quem lê o banco não consegue usar o token.
- **Rotação:** cada renovação troca o refresh token. Token já usado, apresentado de novo depois de 20 segundos, é sinal de cópia: a sessão inteira é revogada e o evento vai para o log.
- **Validade:** cada refresh token vale 7 dias; a sessão, 30 dias no máximo, depois disso é preciso entrar de novo.
- **Revogação:** logout revoga a sessão; apagar a conta apaga as sessões em cascata. O access token já emitido vale até vencer (no máximo 15 minutos).

## CSRF — ✅

- Rotas de dados exigem `Authorization: Bearer`, que outro site não consegue mandar.
- As rotas que usam cookie (`/api/auth/refresh` e `/logout`) recebem o refresh token com `SameSite=Strict`, só em `/api/auth`.
- Toda requisição que muda dados passa pelo `http.CrossOriginProtection` do Go: vinda de outra origem (pelo `Sec-Fetch-Site` ou `Origin`), responde 403, a não ser que a origem esteja em `CORS_ORIGINS`.
