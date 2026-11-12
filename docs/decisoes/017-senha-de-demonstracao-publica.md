# 017. Senha de demonstração pública

- **Data:** 2026-10
- **Status:** ativa

## Contexto

Etapa "Tudo com um `docker compose up`". Registrada no antigo `docs/MEMORY.md`.

## Decisão

**Senha de demonstração pública:** quem visita o portfólio precisa conseguir entrar. Antes os usuários do seed tinham um texto no lugar do hash e ninguém logava.

## Consequência

A senha é definida em `internal/seed/seed.go` (`DemoPassword`) e não é escrita nos docs.
