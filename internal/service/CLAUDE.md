# internal/service

Regra de negócio, validação das entradas (`inputs.go`) e DTOs de resposta (`responses.go`). O handler (`internal/handler`) só lê a requisição, chama o service e devolve JSON.

## Regras

- O service depende das **interfaces** de `internal/repository/interfaces.go`. Mudou a interface, atualize o mock em `internal/service/mocks` (escrito à mão).
- Devolva DTOs (`UserResponse`, `RecipeResponse`...), nunca o model: evita vazar `PasswordHash` e afins.
- Erro esperado sai como `service.Error` com um `Kind` (`errors.go`); o handler converte em status com `respondError`. Erro de outro tipo vira 500 genérico, com o detalhe só no log. Nunca decida status comparando o texto do erro.
- Erro inesperado sobe com contexto: `fmt.Errorf("contexto: %w", err)`.
- Toda função que chega ao banco recebe o `context.Context` da requisição.
- Dono da receita e admin são conferidos aqui, não no handler. Item de outra receita responde **404**, não 403 (403 confirmaria que o id existe).
- Limite de campo novo entra em `inputs.go` **e** no `src/lib/limits.ts` do front; cada arquivo aponta para o outro.
- Erro de regra (categoria inexistente) é `KindInvalid` (400); validação de campo é 422 com o caminho do JSON (`ingredients[0].amount`).
- `UpdateUserInput` usa ponteiros: `nil` é "não mandou", `""` é "apagar" (foto do perfil).

## Armadilhas conhecidas

- **Senha:** 10 a 72 bytes (`MaxPasswordBytes`, limite do bcrypt), conferido antes do bcrypt. E-mail inexistente também passa pelo bcrypt, para o tempo de resposta não revelar quem tem conta.
- **Bloqueio por e-mail:** a tentativa é contada **antes** do bcrypt; conferir antes e contar depois deixava passar várias senhas em paralelo. E-mail inexistente conta igual. Login certo zera a contagem.
- **Refresh token:** opaco, guardado como SHA-256. A troca é atômica (`SELECT ... FOR UPDATE`); reúso depois de 20 s derruba a sessão inteira. Renovação que falha no servidor (banco fora) **não** apaga os cookies; só sessão encerrada (401) apaga.
- **`PUT /recipes/:id`** troca a receita inteira numa transação: os itens antigos são apagados e recriados na ordem recebida (o id dos itens muda, e nada guarda esse id).
- **Upload:** até 900 KB, tipo pelos primeiros bytes (`http.DetectContentType`), endereço relativo `/api/uploads/{id}`. Fotos enviadas e nunca usadas ficam no banco.
- **Vídeo:** só YouTube, guardado sempre como `https://www.youtube.com/watch?v={id}` (`video.go`).
- Comentário e perfil público mostram do autor só id, nome e foto.

## Ver também

`docs/decisoes/` 008, 009, 011, 022, 023, 027, 028, 029–031, 038, 043, 045, 046.
