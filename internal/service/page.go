package service

import "panda-cooking-go-api/internal/repository"

// Tamanhos de página: o padrão quando o cliente não manda per_page e o
// máximo aceito (o handler recusa mais que isso com 422).
const (
	DefaultRecipesPerPage  = 12
	DefaultCommentsPerPage = 10
	MaxPerPage             = 50
)

// PageQuery são os parâmetros de paginação da query string.
type PageQuery struct {
	Page    int `form:"page" binding:"omitempty,min=1,max=10000"`
	PerPage int `form:"per_page" binding:"omitempty,min=1,max=50"`
}

// Page é uma página de resultados, com o total para o front montar a
// navegação.
type Page[T any] struct {
	Items      []T   `json:"items"`
	Page       int   `json:"page"`
	PerPage    int   `json:"per_page"`
	Total      int64 `json:"total"`
	TotalPages int   `json:"total_pages"`
}

func (q PageQuery) repo(defaultSize int) repository.Page {
	p := repository.Page{Number: q.Page, Size: q.PerPage}
	if p.Number < 1 {
		p.Number = 1
	}
	if p.Size < 1 {
		p.Size = defaultSize
	}
	if p.Size > MaxPerPage {
		p.Size = MaxPerPage
	}
	return p
}

// newPage monta a resposta; items nunca sai como null no JSON.
func newPage[T any](items []T, p repository.Page, total int64) Page[T] {
	if items == nil {
		items = []T{}
	}
	return Page[T]{
		Items:      items,
		Page:       p.Number,
		PerPage:    p.Size,
		Total:      total,
		TotalPages: int((total + int64(p.Size) - 1) / int64(p.Size)),
	}
}
