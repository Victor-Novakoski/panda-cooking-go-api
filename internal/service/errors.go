package service

// Kind diz que tipo de falha um erro do service representa.
// O handler usa o tipo para escolher o status HTTP.
type Kind int

const (
	KindNotFound Kind = iota + 1
	KindForbidden
	KindConflict
	KindUnauthorized
	KindTooManyRequests
)

// Error é um erro esperado da regra de negócio. A mensagem pode ir para o
// cliente; qualquer outro erro é tratado como interno e não sai no corpo.
type Error struct {
	Kind    Kind
	Message string
}

func (e *Error) Error() string { return e.Message }

func newError(kind Kind, msg string) *Error { return &Error{Kind: kind, Message: msg} }

var (
	ErrRecipeNotFound      = newError(KindNotFound, "receita não encontrada")
	ErrImageNotFound       = newError(KindNotFound, "imagem não encontrada")
	ErrIngredientNotFound  = newError(KindNotFound, "ingrediente não encontrado nesta receita")
	ErrPreparationNotFound = newError(KindNotFound, "passo de preparo não encontrado")
	ErrCommentNotFound     = newError(KindNotFound, "comentário não encontrado")
	ErrUserNotFound        = newError(KindNotFound, "usuário não encontrado")
	ErrFavoriteNotFound    = newError(KindNotFound, "receita não está nos favoritos")

	ErrForbidden            = newError(KindForbidden, "sem permissão")
	ErrForbiddenEditRecipe  = newError(KindForbidden, "sem permissão para editar esta receita")
	ErrForbiddenDelRecipe   = newError(KindForbidden, "sem permissão para deletar esta receita")
	ErrForbiddenEditComment = newError(KindForbidden, "sem permissão para editar este comentário")
	ErrForbiddenDelComment  = newError(KindForbidden, "sem permissão para deletar este comentário")

	ErrAlreadyFavorite = newError(KindConflict, "receita já está nos favoritos")
	ErrEmailTaken      = newError(KindConflict, "e-mail já cadastrado")

	ErrInvalidCredentials = newError(KindUnauthorized, "email ou senha inválidos")
	ErrTooManyAttempts    = newError(KindTooManyRequests, "muitas tentativas de login, tente novamente mais tarde")
)
