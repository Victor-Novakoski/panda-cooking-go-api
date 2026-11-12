# internal/integration

Testes com a API, os services e os repositórios contra um Postgres de verdade, com as migrations aplicadas. Pegam o que os mocks não pegam: SQL, transação, busca, cascata e as tentativas de mexer em dado de outra pessoa.

## Regras

- Build tag `integration`: `make test` não roda estes testes; `make test-integration` roda. Na CI é um job próprio ("Testes de integração").
- O banco é um container novo (testcontainers) ou o de `TEST_DATABASE_URL`, quando definido.
- Os testes **apagam todos os dados**: o nome do banco em `TEST_DATABASE_URL` precisa conter `test`.

## Armadilhas conhecidas

- **Na nuvem (Claude Code) não há Docker**, então o testcontainers não sobe. Há Postgres 16 em `/usr/lib/postgresql/16`: crie o cluster com `initdb` como usuário `postgres` em `/var/tmp`, suba o servidor, crie um banco com `test` no nome e rode `TEST_DATABASE_URL=... make test-integration`.
- As migrations precisam das extensões `unaccent` e `pg_trgm` no Postgres usado.

## Ver também

`docs/decisoes/` 041.
