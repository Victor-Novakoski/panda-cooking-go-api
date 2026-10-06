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
- **`API_INTERNAL_URL` no front:** a página da receita é renderizada no servidor do Next, que dentro do Docker não alcança `localhost:8080`; ele usa `http://api:8080` e o navegador continua com a URL pública.
