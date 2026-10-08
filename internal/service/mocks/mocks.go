// Package mocks tem implementações das interfaces de repositório feitas à
// mão para os testes de service e de handler, sem banco.
package mocks

import "panda-cooking-go-api/internal/repository"

var (
	_ repository.UserRepo     = (*UserRepoMock)(nil)
	_ repository.RecipeRepo   = (*RecipeRepoMock)(nil)
	_ repository.CategoryRepo = (*CategoryRepoMock)(nil)
	_ repository.CommentRepo  = (*CommentRepoMock)(nil)
	_ repository.FavoriteRepo = (*FavoriteRepoMock)(nil)
	_ repository.SessionRepo  = (*SessionRepoMock)(nil)
)
