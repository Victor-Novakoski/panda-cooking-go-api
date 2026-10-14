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
	// MaxTimeMinutes: só as receitas que ficam prontas em até esse tempo.
	MaxTimeMinutes int
	// Sort vazio é SortRecent.
	Sort RecipeSort
}

// RecipeSort é a ordem da listagem de receitas.
type RecipeSort string

const (
	SortRecent    RecipeSort = "recent"    // mais novas primeiro
	SortQuickest  RecipeSort = "quickest"  // menor tempo de preparo primeiro
	SortFavorites RecipeSort = "favorites" // mais favoritadas primeiro
)
