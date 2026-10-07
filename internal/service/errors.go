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
	KindInvalid
)

// Error é um erro esperado da regra de negócio. A mensagem pode ir para o
// cliente; qualquer outro erro é tratado como interno e não sai no corpo.
// Field, quando preenchido, é o campo do corpo da requisição que causou o
// erro, para o front mostrar a mensagem ao lado dele.
type Error struct {
	Kind    Kind
	Message string
	Field   string
}

func (e *Error) Error() string { return e.Message }

func newError(kind Kind, msg string) *Error { return &Error{Kind: kind, Message: msg} }

func fieldError(kind Kind, field, msg string) *Error {
	return &Error{Kind: kind, Message: msg, Field: field}
}

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
	ErrForbiddenDelRecipe   = newError(KindForbidden, "sem permissão para apagar esta receita")
	ErrForbiddenEditComment = newError(KindForbidden, "sem permissão para editar este comentário")
	ErrForbiddenDelComment  = newError(KindForbidden, "sem permissão para apagar este comentário")

	ErrAlreadyFavorite = newError(KindConflict, "receita já está nos favoritos")
	ErrEmailTaken      = fieldError(KindConflict, "email", "e-mail já cadastrado")

	ErrEmptyName        = fieldError(KindInvalid, "name", "o nome não pode ficar vazio")
	ErrCategoryNotFound = fieldError(KindInvalid, "category_id", "categoria não encontrada")
	ErrTooManyImages    = fieldError(KindInvalid, "images", "a receita pode ter no máximo 10 fotos")
	ErrTooManyItems     = newError(KindInvalid, "a receita já tem o máximo de itens permitido")

	ErrInvalidCredentials = newError(KindUnauthorized, "e-mail ou senha inválidos")
	ErrSessionInvalid     = newError(KindUnauthorized, "sessão expirada, entre de novo")
	ErrTooManyAttempts    = newError(KindTooManyRequests, "muitas tentativas de login, tente novamente mais tarde")
)
