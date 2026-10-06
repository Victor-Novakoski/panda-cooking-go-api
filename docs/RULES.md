# Regras do projeto

Regras que valem para qualquer mudança. Se uma regra atrapalhar, ela é discutida e alterada aqui, não ignorada em silêncio.

## 1. Escopo

- Toda funcionalidade nova precisa estar no [PRD](PRD.md) e em [TASKS.md](TASKS.md) antes de ser implementada.
- Uma tarefa por vez. Melhorias que aparecerem no caminho viram item novo em TASKS, não entram de carona.
- Nada de dependência, serviço ou camada nova "para o futuro". Entra quando uma tarefa precisar.

## 2. Código Go

- Respeitar as camadas de [ARCHITECTURE.md](ARCHITECTURE.md): handler só fala HTTP; regra de negócio fica no service; banco só no repository.
- Banco só pelo GORM com parâmetros (`Where("email = ?", email)`). Nunca montar SQL concatenando strings.
- O service depende das interfaces de `internal/repository/interfaces.go`; mudou a interface, atualiza o mock em `internal/service/mocks`.
- Services devolvem DTOs, nunca o model direto para o handler (evita vazar campos como `Password`).
- Erro inesperado sobe com `fmt.Errorf("contexto: %w", err)` e vira 500 com mensagem genérica; o detalhe vai só para o log.
- Identificadores em inglês, comentários, docs e mensagens de commit em português, como o código atual.
- `make lint` limpo (golangci-lint, mesma configuração da CI).

## 3. Segurança

Checklist para toda mudança (detalhes em [SECURITY.md](SECURITY.md)):

- [ ] Entrada validada no back-end, com tamanho máximo para textos, mesmo que o front já valide.
- [ ] Rota nova tem autenticação e papel definidos explicitamente; rota pública é exceção justificada.
- [ ] Recurso acessado por id confere se pertence a quem pede (IDOR).
- [ ] Resposta não expõe hash, token, dados pessoais desnecessários nem detalhes internos de erro.
- [ ] Nenhum segredo no código, em log ou em commit. Valores novos vão para `.env.example` com placeholder.
- [ ] Dependência nova é necessária, mantida e passa no `govulncheck`.
- [ ] Nada de código, nome de cliente ou padrão interno de outra empresa no repositório.

## 4. Testes

- Regra de negócio nova tem teste no service. Rota nova tem teste do handler (status e corpo).
- Bug corrigido ganha teste que falhava antes da correção.
- `make test` passando antes de qualquer commit.

## 5. API

- Seguir as convenções de [DESIGN.md](DESIGN.md) (formato de erro, status, nomes em snake_case).
- Toda rota nova ou alterada é atualizada na tabela de [ARCHITECTURE.md](ARCHITECTURE.md#rotas) e na coleção do Postman no mesmo commit.

## 6. Git

- **Branches:** `main` é o que está publicado; `develop` junta o trabalho pronto para a próxima versão. Ninguém faz push direto nas duas: tudo entra por pull request.
- **Fluxo:** criar a branch a partir da `develop` com o mesmo tipo do commit (`feat/busca-de-receitas`, `fix/...`, `docs/...`, `chore/...`, `ci/...`), abrir PR para a `develop` e fazer merge com os checks verdes. A branch é apagada automaticamente depois do merge.
- **Versão:** quando a `develop` fecha uma etapa, abrir PR da `develop` para a `main`.
- **Como fazer o merge:** feature → `develop` com squash (um commit por PR); `develop` → `main` sempre com merge commit, nunca rebase ou squash, para as duas branches não divergirem. Para atualizar a branch de feature, `git pull --rebase origin develop`.
- `main` e `develop` sempre funcionando: sobem com `make setup` e passam nos testes.
- **Conventional Commits:** mensagem no formato `tipo: assunto em português, minúsculo`. Tipos: `feat` (funcionalidade), `fix` (correção), `docs`, `test`, `refactor`, `perf`, `style`, `build`, `ci`, `chore`, `revert`. Escopo opcional, ex.: `feat(api): busca de receitas`. Mudança que quebra compatibilidade leva `!`: `feat!: ...`.
- O título do PR segue o mesmo formato, porque o squash usa o título como commit; a CI confere.
- Commits pequenos, um assunto por commit.
- Nunca commitar `.env` nem binários.

## 7. Documentação

- Mudou comportamento, rota, variável de ambiente ou decisão de arquitetura: atualizar o doc correspondente no mesmo commit.
- Decisão relevante (escolha de lib, trade-off, algo que foi descartado) vai para [MEMORY.md](MEMORY.md).
- Tarefa concluída é marcada em [TASKS.md](TASKS.md).

## Definição de pronto

Uma tarefa só está pronta quando: funciona com `docker compose up`, tem testes, passa no lint, o checklist de segurança foi revisado, OpenAPI e docs estão atualizados e o item está marcado em TASKS.
