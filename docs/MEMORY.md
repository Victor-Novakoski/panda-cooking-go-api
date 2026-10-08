# Memória do projeto

Decisões e o porquê delas. Decisão nova entra aqui ([RULES.md](RULES.md#7-documentação)).

## 2026-10 — Setup

- **Mesmo fluxo do Rastreia:** `develop` + PRs, squash nas features, merge commit na versão, Conventional Commits em português, CI obrigatória antes do merge. Assim os dois projetos do portfólio contam a mesma história.
- **Mantida a stack que já existia (Gin + GORM):** o Rastreia mostra chi + sqlc; este mostra a outra abordagem comum no mercado. Trocar agora seria reescrever sem ganho.
- **Go 1.26.8:** a versão no `go.mod` define o Go da CI; com a 1.25.0 o govulncheck acusaria falhas já corrigidas na biblioteca padrão.
- **`misspell` fora do lint:** ele só conhece inglês e acusava os comentários em português.
- **Imagem distroless `nonroot`:** sem shell e sem root, menos superfície de ataque e o Trivy fica limpo.
- **Repositório público:** nada de código, nome de cliente ou padrão interno de outra empresa entra aqui.
- **Dependências atualizadas no setup:** pgx 5.6 tinha falha crítica (e risco de SQL injection) e x/crypto, x/net e x/text tinham falhas altas; o govulncheck e o Trivy barraram na primeira CI.

## 2026-10 — Segurança da base

- **Erros tipados em vez de comparar texto:** antes o handler decidia o status comparando `err.Error()` com frases, e o que não batia saía inteiro na resposta. Agora o service devolve `service.Error` com um tipo, e erro desconhecido nunca chega ao cliente. Mudar a frase de um erro não quebra mais o status.
- **Item de outra receita responde 404, não 403:** 403 confirmaria que o id existe em outra receita.
- **Limite de login em memória, sem Redis:** a API roda numa instância só; o Rastreia usa Redis porque tem WebSocket e várias réplicas. Se escalar, o contador vai para um armazenamento compartilhado.
- **Bloqueio por e-mail conta e-mail inexistente:** senão o bloqueio só aconteceria para quem tem conta e serviria para descobrir e-mails cadastrados.
- **`SetTrustedProxies(nil)`:** com o padrão do Gin (confia em todo proxy), o cliente trocaria de "IP" a cada requisição mandando `X-Forwarded-For`.
- **Chave de exemplo só é recusada em produção:** em desenvolvimento o `.env.example` funciona copiado, sem passo extra.
- **`TranslateError` do GORM ligado:** transforma a violação de unicidade do Postgres em `gorm.ErrDuplicatedKey`, que o service converte em 409 sem depender do texto do erro do banco.

## 2026-10 — Tudo com um `docker compose up`

- **Compose no repositório da API, front clonado ao lado:** banco, migração e seed são da API, então ela é quem sabe subir o ambiente. O front entra por bind mount de `../panda-cooking-front` (ou `FRONT_DIR`), sem imagem publicada nem submódulo. Um terceiro repositório só para o compose seria mais uma coisa para manter sincronizada.
- **Seed de demonstração dentro da API (`SEED_DEMO=true`):** ao subir, cria categorias, três usuários e dez receitas se não houver usuário nenhum. Não depende de um container a mais nem de Go na máquina, e não apaga o que já foi cadastrado. `make seed` recria tudo do zero.
- **Senha de demonstração pública (`panda-cooking-demo`):** quem visita o portfólio precisa conseguir entrar. Antes os usuários do seed tinham um texto no lugar do hash e ninguém logava.
- **Categorias por nome no seed:** antes as receitas apontavam para `CategoryID: 1..12` e dependiam do `RESTART IDENTITY`; se a API tivesse criado as categorias padrão antes, as receitas iam para a categoria errada.
- **Portas só em `127.0.0.1` (#21):** banco com senha padrão não fica exposto na rede local.
- **Polling no air e no Next:** no Docker do Windows e do Mac a pasta montada não avisa o container de arquivo alterado; sem polling o hot reload não funciona.
- **`API_INTERNAL_URL` no front:** a página da receita é renderizada no servidor do Next, que dentro do Docker não alcança `localhost:8080`; ele usa `http://api:8080` e o navegador continua com a URL pública. (Substituído na etapa 5 pelo `API_URL`, ver abaixo.)

## 2026-10 — Telas de edição

- **`PUT /recipes/:id` com a receita inteira:** a tela de edição manda tudo de uma vez e a API troca dados, fotos, ingredientes e passos numa transação. Com as rotas por item o front teria que fazer várias chamadas, e uma falha no meio deixaria a receita pela metade. As rotas por item continuam para quem quiser mexer em uma parte só.
- **Itens recriados, não comparados:** no `PUT` os itens antigos são apagados e os novos criados na ordem recebida. Comparar item a item não traz nada para quem usa a tela e complicaria o código; o custo é o id do item mudar, e nada guarda esse id.
- **`Omit(clause.Associations)` no `Save`:** com a `Category` carregada pelo `Preload`, o `Save` do GORM gravava a associação e voltava o `category_id` para o antigo, e trocar a categoria no `PATCH` não tinha efeito.
- **Ordem por id nos `Preload`:** sem `ORDER BY` o Postgres não garante a ordem dos passos do preparo.
- **Comentário mostra só nome e foto do autor:** a lista de comentários é pública; antes a resposta trazia o e-mail e o `is_adm` de quem comentou.
- **Ponteiro no `UpdateUserInput`:** com `string`, vazio e "não mandou" eram a mesma coisa e não havia como tirar a foto do perfil.
- **`KindInvalid` (400) para erro de regra:** categoria inexistente era 500 (chave estrangeira violada); agora é 400 com a mensagem.

## 2026-10 — Sessão em cookie, API completa e revisão geral

- **Refresh token em cookie `HttpOnly` + access token de 15 min, não JWT de 24h no `localStorage`:** com o token legível por JavaScript, um XSS levava a sessão por 24h sem como revogar. O JWT curto continua sem consulta ao banco em cada requisição; quem segura a sessão é o refresh token, que dá para revogar.
- **Refresh token opaco guardado como SHA-256, não JWT nem bcrypt:** precisa ser revogável (então vive no banco), e quem lê o banco não deve conseguir usá-lo. Como é aleatório e longo, SHA-256 basta; bcrypt existe para senha fraca.
- **Rotação com detecção de reúso e 20 s de tolerância:** token usado de novo indica cópia e derruba a sessão. A tolerância cobre duas abas renovando juntas e a resposta que não chegou ao navegador, que derrubariam a sessão de quem não fez nada de errado. A troca é atômica (`SELECT ... FOR UPDATE` na transação).
- **Duas tabelas (`sessions` e `refresh_tokens`):** revogar a sessão invalida todos os tokens dela de uma vez, e o limite absoluto de 30 dias fica na sessão, não em cada token.
- **Cookie `panda_session` separado:** o refresh token só vai para `/api/auth`, então o servidor do front não o vê nas páginas. O `panda_session` não tem segredo; só diz que existe sessão, para o front redirecionar rota protegida no servidor e não tentar renovar quem nunca entrou.
- **Rotas sob `/api`:** o front repassa `/api/*` para cá pela mesma origem; com o prefixo, o proxy reverso também consegue mandar `/api` direto para a API, sem conflito com as páginas.
- **`http.CrossOriginProtection` do Go 1.25+ em vez de token CSRF:** usa `Sec-Fetch-Site`/`Origin`, que todo navegador atual manda, sem estado e sem mudar o front. Fica antes do CORS para a recusa sair no formato de erro da API.
- **Migrations em SQL com golang-migrate, embutidas no binário:** o `AutoMigrate` não apaga coluna, não cria extensão nem índice funcional (a busca precisa dos dois) e não tem histórico. Banco criado pela versão antiga é recusado com uma mensagem dizendo para recriar (`docker compose down -v`), em vez de tentar adivinhar o estado dele.
- **Busca com `unaccent` + `pg_trgm`:** "pao" acha "Pão", e o `LIKE '%termo%'` usa índice GIN. Full-text do Postgres foi descartado: corta palavras pela raiz e não acha pedaço de palavra, o que confunde numa busca de receita.
- **422 com erro por campo no caminho do JSON (`ingredients[0].amount`):** o front mostra a mensagem no campo certo, inclusive dentro de listas. O 400 continua para JSON malformado e campo desconhecido.
- **Paginação com `total` e `total_pages`:** o front monta a navegação sem uma chamada a mais. Offset em vez de cursor: a lista é pequena e ordenada por data, e o usuário pula para uma página específica.
- **`context.Context` até o repository:** com o prazo de 10 s e o cliente que desiste, a consulta é cancelada no banco em vez de continuar ocupando conexão.
- **Testes de integração com testcontainers, num job próprio:** os testes de unidade usam mocks e não pegam erro de SQL, de transação ou de migration. Rodam com a build tag `integration`, para o `make test` continuar rápido e sem Docker.
- **OpenAPI escrito à mão em vez de gerado por anotação:** a especificação é o contrato com o front e fica legível em um arquivo só. Um teste confere que toda rota do roteador está documentada, para os dois não se separarem. Substitui a coleção do Postman, que já estava desatualizada.
- **Limites iguais aos do front:** `inputs.go` e o `src/lib/limits.ts` do front têm os mesmos números, e cada um aponta para o outro.
- **`healthcheck` no próprio binário:** a imagem distroless não tem `curl` nem shell; `api healthcheck` faz a chamada ao `/health` e serve de `HEALTHCHECK` do Docker.
- **Tentativa de login contada antes do bcrypt:** conferir o bloqueio antes e contar depois deixava passar várias senhas ao mesmo tempo (8 em vez de 5 num teste com 10 simultâneas). O login certo zera a contagem, então quem acerta não perde nada.
- **Renovação que falha no servidor não apaga os cookies:** com o banco fora do ar a troca é desfeita e a sessão continua valendo; apagar o cookie deslogaria todo mundo que renovasse durante a queda. Só a sessão encerrada (401) apaga.
- **Revisão adversarial antes do PR:** os achados (bloqueio burlável com requisições simultâneas, logout em queda do banco, id gigante dando 500, conta apagada aparecendo como "receita não encontrada") viraram correção com teste.

