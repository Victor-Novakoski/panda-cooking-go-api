# internal/database (e internal/repository)

Conexão, pool e migrations (`migrations/*.sql`, embutidas no binário e aplicadas quando a API sobe). O acesso ao banco mora em `internal/repository`, atrás das interfaces de `interfaces.go`.

## Regras

- Mudança no banco é uma migration nova: `NNNNNN_nome.up.sql` e `.down.sql`. Migration já aplicada na `develop` nunca é editada.
- Nada de `AutoMigrate`: ele não cria extensão nem índice funcional (a busca precisa dos dois) e não tem histórico.
- Banco só pelo GORM com parâmetros (`Where("email = ?", email)`). Nunca montar SQL concatenando string. Termo de busca entra no `LIKE` com `%` e `_` escapados.
- Toda consulta usa o `context.Context` recebido (`db.WithContext(ctx)`), para o prazo de 10 s cancelar no banco.
- Consulta nova no repository ganha teste em `internal/integration`.

## Armadilhas conhecidas

- **Banco antigo** (criado pelo `AutoMigrate`, antes das migrations) é recusado com `ErrLegacySchema`. Solução: `docker compose down -v` uma vez.
- **`Omit(clause.Associations)` no `Save`:** com a `Category` carregada por `Preload`, o `Save` gravava a associação e voltava o `category_id` para o antigo.
- **`Preload` sem ordem:** o Postgres não garante a ordem dos passos do preparo. Ordene por id.
- **`TranslateError: true`** transforma violação de unicidade em `gorm.ErrDuplicatedKey` (vira 409 no service) e de chave estrangeira em `gorm.ErrForeignKeyViolated`. Não compare o texto do erro do Postgres.
- Id que não é UUID quebraria a consulta com erro de sintaxe; o handler responde 404 antes.
- A busca usa `unaccent` + `pg_trgm` com índice GIN (migration 000003). Full-text do Postgres foi descartado: corta palavras pela raiz e não acha pedaço de palavra.

## Ver também

`docs/decisoes/` 014, 024, 025, 036, 037, 040.
