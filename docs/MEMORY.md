# Memória do projeto

Decisões e o porquê delas. Decisão nova entra aqui ([RULES.md](RULES.md#7-documentação)).

## 2026-10 — Setup

- **Mesmo fluxo do Rastreia:** `develop` + PRs, squash nas features, merge commit na versão, Conventional Commits em português, CI obrigatória antes do merge. Assim os dois projetos do portfólio contam a mesma história.
- **Mantida a stack que já existia (Gin + GORM):** o Rastreia mostra chi + sqlc; este mostra a outra abordagem comum no mercado. Trocar agora seria reescrever sem ganho.
- **Go 1.26.8:** a versão no `go.mod` define o Go da CI; com a 1.25.0 o govulncheck acusaria falhas já corrigidas na biblioteca padrão.
- **`misspell` fora do lint:** ele só conhece inglês e acusava os comentários em português.
- **Imagem distroless `nonroot`:** sem shell e sem root, menos superfície de ataque e o Trivy fica limpo.
- **Repositório público:** nada de código, nome de cliente ou padrão interno de outra empresa entra aqui.
