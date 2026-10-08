package repository

// Page é uma página de resultados. Number começa em 1.
type Page struct {
	Number int
	Size   int
}

func (p Page) offset() int { return (p.Number - 1) * p.Size }

// RecipeFilter são os filtros da listagem de receitas; campo vazio não filtra.
type RecipeFilter struct {
	Search     string
	CategoryID uint
	UserID     string
}
