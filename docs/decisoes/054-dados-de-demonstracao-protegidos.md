# 054. Dados de demonstração protegidos com `is_demo`

- **Data:** 2026-11-19
- **Status:** ativa

## Contexto

A senha das contas do seed é pública (017). Com o site no ar, qualquer visitante entrava como a admin de demonstração e podia apagar as receitas dela, mudar o perfil ou apagar a conta (a cascata levava receitas e comentários). O próximo visitante encontraria o site vazio. Issue #30.

## Decisão

Coluna `is_demo` em `users` e em `recipes` (migration 000006), marcada pelo seed. O service recusa com 403 (`ErrDemoRecipe`, `ErrDemoUser`):

- mudar ou apagar receita de demonstração, inclusive fotos, ingredientes e passos (a regra fica em `findOwned`, por onde todo caminho de escrita passa);
- mudar o perfil ou apagar a conta de demonstração.

Continua liberado: comentar, favoritar, publicar receita nova (nasce com `is_demo = false`, então quem publicou pode editar e apagar) e a admin de demonstração apagar comentários. `is_demo` sai no perfil e na receita para o front esconder os botões.

Marcar a receita e não só o usuário separa o que o seed criou do que a visita publicou com a mesma conta.

## Consequência

- A migration marca os bancos que já tinham o seed: usuários pelos e-mails de `internal/seed/data.go` e **todas** as receitas deles. Receita que uma visita tenha publicado com essas contas antes da migration também fica protegida; é aceito.
- Comentários de demonstração podem ser apagados pela admin; o seed recriado (`make seed`) os traz de volta.
- Mudou o e-mail de um usuário do seed? A migration não muda (já aplicada); o seed novo marca sozinho.
