# internal/seed

Dados de demonstração: três usuários e dezoito receitas brasileiras, com fotos, ingredientes, preparo, comentários e favoritos. As 12 categorias vêm da migration, não do seed.

## Regras

- **A senha dos usuários de demonstração é definida em `seed.go` (`DemoPassword`).** Ela é pública de propósito (quem visita o portfólio precisa entrar), mas não se escreve em doc nenhum: aponte para a constante.
- Com `SEED_DEMO=true` (padrão no compose), a API cria os dados ao subir **só se não houver usuário nenhum**; nada cadastrado é apagado. Com `false`, os usuários não existem.
- `make seed` (`cmd/seed`) **apaga todos os dados** antes de popular. Só em desenvolvimento.
- Receita nova entra em `demoRecipes` (`data.go`): o autor é a posição em `demoUsers`, a categoria é o **nome** (nunca o id: o id depende da ordem de criação).
- Os textos seguem os limites da API; `data_test.go` confere.

## Armadilhas conhecidas

- Antes de usar uma foto, abra o link e confira se ela mostra o prato: link do Unsplash some. A mesma foto não pode servir para duas receitas (o teste confere).
- Fotos de outro host funcionam, mas o front só otimiza as de `images.unsplash.com`.

## Ver também

`docs/decisoes/` 016, 017, 018.
