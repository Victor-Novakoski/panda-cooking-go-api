# Segurança

Como a API trata cada risco, agrupado por tema. Os riscos que dependem do front estão no [SECURITY.md do front](https://github.com/Victor-Novakoski/panda-cooking-front/blob/develop/docs/SECURITY.md). O checklist para cada mudança está no [`CLAUDE.md`](../CLAUDE.md).

Situação revisada em 07/10/2026 (etapa 5) e conferida com o código em 12/11/2026.

**Legenda:** ✅ feito · 🟡 parcial

| Tema | Riscos | Situação |
| --- | --- | --- |
| [Autenticação](seguranca/autenticacao.md) | autenticação fraca, senhas no banco, força bruta, tokens, CSRF | ✅ |
| [Dados](seguranca/dados.md) | validação no front e no back-end, SQL injection, IDOR, envio duplicado, upload, vazamento de informação, dados sensíveis, XSS | ✅ |
| [Infraestrutura](seguranca/infraestrutura.md) | segredos, dependências, rate limit, CORS, headers, timeouts, banco exposto, logs, IP do cliente | 🟡 o IP do cliente depende do deploy |

O que ainda falta para produção está na issue de publicação.
