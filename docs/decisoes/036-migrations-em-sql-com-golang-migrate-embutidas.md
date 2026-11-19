# 036. Migrations em SQL com golang-migrate, embutidas no binário

- **Data:** 2026-10
- **Status:** ativa

## Contexto

Etapa "Sessão em cookie, API completa e revisão geral". Registrada no antigo `docs/MEMORY.md`.

## Decisão

**Migrations em SQL com golang-migrate, embutidas no binário:** o `AutoMigrate` não apaga coluna, não cria extensão nem índice funcional (a busca precisa dos dois) e não tem histórico. Banco criado pela versão antiga é recusado com uma mensagem dizendo para recriar (`docker compose down -v`), em vez de tentar adivinhar o estado dele.

## Consequência

O porquê e o efeito foram registrados junto com a decisão, no texto acima.
