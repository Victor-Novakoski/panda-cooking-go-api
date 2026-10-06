# Segurança

Como a API trata cada risco, o que já está feito e o que falta. As tarefas estão em [TASKS.md](TASKS.md) e o checklist rápido para cada mudança está em [RULES.md](RULES.md#3-segurança).

Situação revisada em 06/10/2026, na etapa 1 (setup).

**Legenda:** ✅ feito · 🟡 parcial · 🔴 pendente · ⚪ ainda não se aplica

## Resumo

| # | Risco | Situação | Prioridade | O que é |
| --- | --- | --- | --- | --- |
| 1 | Variáveis de ambiente expostas | 🟡 | Alta | Segredo vazando por arquivo commitado, imagem Docker ou valor padrão usado em produção. |
| 2 | Validação no front-end | 🟡 | Média | Dar retorno rápido no formulário. Ajuda o usuário, mas não protege nada. |
| 3 | Validação no back-end | 🟡 | Alta | A API conferir tipo, formato e tamanho de tudo que recebe. |
| 4 | SQL Injection | ✅ | — | Texto do usuário virando parte do comando SQL. |
| 5 | Autenticação fraca | 🟡 | Média | Senha fraca, login que dá pistas, credencial padrão. |
| 6 | IDOR | 🔴 | **Alta** | Trocar um id na URL e mexer em dado de outra pessoa. |
| 7 | Senhas no banco | ✅ | — | Guardar a senha em texto puro. |
| 8 | Força bruta | 🔴 | Alta | Tentar milhares de senhas seguidas. |
| 9 | Envio duplicado | 🔴 | Baixa | Clique duplo criando o mesmo registro duas vezes. |
| 10 | CSRF | ✅ | — | Outro site fazendo o navegador logado enviar uma ação. |
| 11 | Upload sem validação | ⚪ | — | Arquivo malicioso ou grande demais. Hoje imagem é só URL. |
| 12 | Vazamento de informação | 🔴 | Alta | Erro interno ou detalhe do banco aparecendo na resposta. |
| 13 | Dependências vulneráveis | ✅ | — | Biblioteca com falha conhecida. |
| 14 | Tokens | 🟡 | Média | Token longo demais, sem revogação, guardado onde script lê. |
| 15 | Rate limit | 🔴 | Média | Limitar requisições por cliente. |
| 16 | Dados sensíveis expostos | 🟡 | Média | Resposta ou log mostrando dado pessoal além do necessário. |
| 17 | CORS | 🟡 | Média | Quais sites podem chamar a API. |
| 18 | XSS | 🟡 | Média | Script injetado num dado rodando no navegador de outra pessoa. |
| 19 | Headers de segurança | 🔴 | Baixa | Cabeçalhos que mandam o navegador se proteger. |
| 20 | Timeouts e negação de serviço | 🔴 | Média | Conexão lenta ou corpo enorme prendendo o servidor. |
| 21 | Banco exposto | 🟡 | Média | Banco acessível pela rede sem passar pela API. |
| 22 | Logs e auditoria | 🔴 | Baixa | Registrar eventos suspeitos para investigar. |

---

## 1. Variáveis de ambiente expostas — 🟡

**Feito:** `.env` no `.gitignore` e no `.dockerignore`; segredos só por variável de ambiente; gitleaks na CI varre todo o histórico a cada PR.

**Falta:** a API sobe com `SECRET_KEY=changeme` se a variável não existir. Exigir chave com 32+ caracteres e recusar a de exemplo em produção (`APP_ENV=production`).

## 2. Validação no front-end — 🟡

Os formulários usam zod. Falta alinhar os limites com os da API quando a etapa 2 definir os tamanhos máximos.

## 3. Validação no back-end — 🟡

**Feito:** `binding` do Gin nos DTOs (obrigatório, e-mail, URL, porções ≥ 1).

**Falta:** tamanho máximo em todos os textos; trim e e-mail em minúsculas; limite de tamanho do corpo; recusar campos desconhecidos; senha com mais de 72 bytes (limite do bcrypt) devolver erro de validação em vez de 500.

## 4. SQL Injection — ✅

Todo acesso ao banco passa pelo GORM com parâmetros. Regra em [RULES.md](RULES.md): nunca montar SQL concatenando string, nem em `Where`.

## 5. Autenticação fraca — 🟡

**Feito:** bcrypt; mesma mensagem para e-mail inexistente e senha errada.

**Falta:** senha mínima sobe de 6 para 10 caracteres e recusa senhas comuns; tempo de resposta igual quando o e-mail não existe (hoje responde mais rápido, o que entrega quem tem conta).

## 6. IDOR — 🔴

**Feito:** editar ou apagar receita, comentário e itens da receita confere se a receita é de quem pede.

**Bug encontrado:** nas rotas de itens (`/recipes/:id/images/:imageID`, `/preparations/:prepID`, `/ingredients/:ingredientID`) o service confere o dono da receita `:id`, mas **não confere se o item pertence a essa receita**. Dono da receita A consegue apagar ou editar a foto, o passo ou o ingrediente da receita B passando o id do item. Correção na etapa 2, com teste que falha antes.

## 7. Senhas no banco — ✅

Só o hash bcrypt é salvo. As respostas usam `UserResponse`, que não tem o campo da senha.

## 8. Força bruta — 🔴

Sem limite de tentativas no `POST /auth`. Etapa 2: limite por IP e bloqueio progressivo por e-mail.

## 9. Envio duplicado — 🔴

Clique duplo no front cria duas receitas ou comentários. Desabilitar o botão durante o envio resolve a maior parte.

## 10. CSRF — ✅

A API só aceita o token pelo cabeçalho `Authorization`, que outro site não consegue enviar. Se o token for para cookie (item 14), entra `SameSite` e proteção de CSRF junto.

## 11. Upload sem validação — ⚪

Hoje a imagem é uma URL validada com `binding:"url"`. Se virar upload: limite de tamanho, tipo conferido pelo conteúdo e nome gerado pelo servidor.

## 12. Vazamento de informação — 🔴

Vários handlers devolvem `err.Error()` direto, inclusive em 500, o que mostra mensagem do Postgres (ex.: e-mail duplicado ao cadastrar). Etapa 2: erros tipados no service, 500 com mensagem genérica e detalhe só no log; Gin em modo release em produção.

## 13. Dependências vulneráveis — ✅

govulncheck na CI a cada PR, Trivy na imagem Docker e Dependabot semanal para Go, Actions e Docker.

## 14. Tokens — 🟡

**Feito:** HS256 com algoritmo conferido na validação; validade de 24h.

**Falta:** o front guarda o token em `localStorage` e num cookie que o JavaScript lê, então um XSS rouba a sessão. Discutir na etapa do front: token curto + refresh em cookie `HttpOnly`, como no Rastreia.

## 15. Rate limit — 🔴

Sem limite. Entra junto com o item 8.

## 16. Dados sensíveis expostos — 🟡

Respostas não têm senha. Conferir na etapa 2 se as rotas públicas de receita e comentário não devolvem o e-mail do autor.

## 17. CORS — 🟡

Lista fixa com `localhost:3000` e `3001`. Falta ler as origens de `CORS_ORIGINS` e exigir `https` em produção.

## 18. XSS — 🟡

O React escapa o texto por padrão. Regra: nunca usar `dangerouslySetInnerHTML` com dado do usuário; URLs de imagem só `http`/`https`.

## 19. Headers de segurança — 🔴

Middleware com `X-Content-Type-Options`, `X-Frame-Options`, `Referrer-Policy` e, com HTTPS, `Strict-Transport-Security`.

## 20. Timeouts e negação de serviço — 🔴

`r.Run` não tem timeout. Trocar por `http.Server` com `ReadHeaderTimeout`, `ReadTimeout`, `WriteTimeout` e `IdleTimeout`, e limitar o tamanho do corpo. Listagens precisam de paginação.

## 21. Banco exposto — 🟡

O compose de dev publica o Postgres na porta 5433 de todas as interfaces. Publicar só em `127.0.0.1`. Em produção o banco não publica porta.

## 22. Logs e auditoria — 🔴

Registrar falhas de login (sem o e-mail em claro), 401, 403 e 429.
