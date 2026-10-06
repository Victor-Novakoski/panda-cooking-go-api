# PRD — Panda Cooking

## Problema

Quem cozinha em casa guarda receitas em print, caderno e grupo de WhatsApp. Fica difícil achar, compartilhar e saber se a receita é boa.

## Solução

Uma rede de receitas: qualquer pessoa navega e lê receitas; quem cria conta publica as suas, comenta e salva as favoritas.

## Público

- **Visitante:** navega pelas receitas e categorias, lê comentários.
- **Usuário cadastrado:** publica e edita as próprias receitas (fotos, ingredientes, modo de preparo), comenta e favorita.
- **Admin:** modera, podendo apagar comentários de qualquer pessoa.

## Funcionalidades

| Funcionalidade | Situação |
| --- | --- |
| Cadastro e login (JWT) | ✅ |
| Perfil: ver, editar, apagar conta | ✅ |
| Listar e ver receitas, filtrar por categoria | ✅ (filtro no front) |
| Publicar, editar e apagar receita própria | ✅ |
| Fotos, ingredientes e passos do preparo | ✅ |
| Comentários (dono edita; dono ou admin apaga) | ✅ |
| Favoritos | ✅ |
| Busca e paginação na API | ⏳ etapa 3 |
| Upload de imagem (hoje é só URL) | 💭 a decidir |

## Fora do escopo por enquanto

Avaliação com estrelas, seguir outros usuários, notificações e app mobile.
