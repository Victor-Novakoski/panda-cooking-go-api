package handler

import "net/http"

// statusFromErr mapeia mensagens de erro conhecidas para status HTTP.
func statusFromErr(err error) int {
	switch err.Error() {
	case "receita não encontrada",
		"imagem não encontrada",
		"ingrediente não encontrado nesta receita",
		"passo de preparo não encontrado",
		"comentário não encontrado",
		"receita não está nos favoritos":
		return http.StatusNotFound
	case "sem permissão",
		"sem permissão para editar esta receita",
		"sem permissão para deletar esta receita",
		"sem permissão para editar este comentário",
		"sem permissão para deletar este comentário":
		return http.StatusForbidden
	case "receita já está nos favoritos":
		return http.StatusConflict
	default:
		return http.StatusInternalServerError
	}
}
