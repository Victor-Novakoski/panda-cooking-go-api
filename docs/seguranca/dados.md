# Segurança: dados

## Validação no front-end — ✅

Os formulários do front usam Zod com os mesmos limites da API (`internal/service/inputs.go` ↔ `src/lib/limits.ts` do front) e mostram no campo certo os erros 422 que a API devolve. Ajuda o usuário, mas não protege nada.

## Validação no back-end — ✅

Cada entrada tem limite de tamanho em todos os textos e listas (até 10 fotos, 50 ingredientes e 50 passos), trim e e-mail em minúsculas antes de validar, URLs só `https`, porções de 1 a 100. Campo desconhecido no JSON é recusado, e o corpo é limitado a 1 MB (413). Erro de validação responde 422 com a mensagem de cada campo, no caminho do JSON (`ingredients[2].amount`), sem nome interno de struct. Parâmetros de busca e paginação também são validados (`per_page` até 50, `search` até 100 caracteres).

## SQL injection — ✅

Todo acesso ao banco passa pelo GORM com parâmetros, inclusive o termo da busca, que entra no `LIKE` com `%` e `_` escapados. Nunca montar SQL concatenando string.

## IDOR — ✅

Editar ou apagar receita, comentário e itens da receita confere se a receita é de quem pede (403). Item de outra receita responde 404, como se não existisse, e nada é gravado (`TestRecipeService_ItemOfAnotherRecipe`). Favoritos e perfil usam sempre o usuário do token, nunca um id da URL. Os testes de integração repetem essas tentativas contra a API e o banco de verdade.

## Envio duplicado — ✅

O front desabilita o botão enquanto a requisição não volta. Na API, favoritar de novo responde 409 (índice único), e repetir a renovação de sessão dentro de 20 segundos é aceito sem derrubar a sessão (duas abas renovando ao mesmo tempo).

## Upload de fotos — ✅

`POST /api/uploads` exige login e aceita 30 fotos por minuto por IP. A foto tem até 900 KB e o tipo é conferido pelos primeiros bytes do arquivo (JPG, PNG, WebP ou GIF), não pelo nome nem pelo `Content-Type` que o navegador mandou. Ela fica no banco (`uploads`), com id gerado pelo servidor, e é servida em `GET /api/uploads/:id` com o tipo detectado e `nosniff`. Fotos por link continuam aceitas, só `https`.

## Vazamento de informação — ✅

O service devolve erros tipados (`service.Error` com um `Kind`) e o handler escolhe o status por eles (`respondError`). Qualquer outro erro vira 500 com "erro interno, tente novamente mais tarde" e o detalhe vai só para o log, com o `X-Request-Id`. E-mail duplicado responde 409; id que não é UUID responde 404 em vez do erro de sintaxe do Postgres. Gin em modo release com `APP_ENV=production`; panic vira 500 sem stack trace na resposta.

## Dados sensíveis expostos — ✅

Receitas, comentários e o perfil público mostram do autor só id, nome e foto (sem e-mail nem `is_adm`). O e-mail só aparece no perfil do próprio usuário. No log, a falha de login identifica o e-mail por um HMAC com a `SECRET_KEY`, nunca em claro.

## XSS — ✅

A API só devolve JSON, com `X-Content-Type-Options: nosniff` e `Content-Security-Policy: default-src 'none'`, então uma resposta aberta direto no navegador não executa nada. URLs de imagem só `https`. No front, o React escapa o texto e há CSP com nonce.
