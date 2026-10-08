package service_test

import (
	"context"
	"testing"

	"panda-cooking-go-api/internal/model"
	"panda-cooking-go-api/internal/service"
	"panda-cooking-go-api/internal/service/mocks"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

func ptr[T any](v T) *T { return &v }

func TestUserService_Create(t *testing.T) {
	input := service.CreateUserInput{Name: "Maria", Email: "maria@email.com", Password: "panelas-de-barro"}

	t.Run("grava o hash da senha, nunca a senha", func(t *testing.T) {
		var saved *model.User
		repo := &mocks.UserRepoMock{CreateFn: func(_ context.Context, u *model.User) error {
			saved = u
			u.ID = userID
			return nil
		}}

		resp, err := service.NewUserService(repo).Create(ctx, input)

		require.NoError(t, err)
		assert.Equal(t, userID, resp.ID)
		assert.Equal(t, "maria@email.com", resp.Email)
		assert.False(t, resp.IsAdm)
		assert.NotEqual(t, input.Password, saved.PasswordHash)
		assert.NoError(t, bcrypt.CompareHashAndPassword([]byte(saved.PasswordHash), []byte(input.Password)))
	})

	t.Run("e-mail repetido é 409 no campo email", func(t *testing.T) {
		repo := &mocks.UserRepoMock{CreateFn: func(context.Context, *model.User) error { return gorm.ErrDuplicatedKey }}

		_, err := service.NewUserService(repo).Create(ctx, input)

		assert.ErrorIs(t, err, service.ErrEmailTaken)
		assert.Equal(t, "email", service.ErrEmailTaken.Field)
	})

	t.Run("senha comum é recusada antes de gravar", func(t *testing.T) {
		weak := input
		weak.Password = "1234567890"

		_, err := service.NewUserService(&mocks.UserRepoMock{}).Create(ctx, weak)

		var appErr *service.Error
		require.ErrorAs(t, err, &appErr)
		assert.Equal(t, service.KindInvalid, appErr.Kind)
		assert.Equal(t, "password", appErr.Field)
	})
}

func TestUserService_Update(t *testing.T) {
	current := func() *model.User {
		return &model.User{ID: userID, Name: "Maria", Email: "maria@email.com", ImageProfile: "https://img.exemplo.com/maria.jpg"}
	}

	t.Run("muda só o que veio", func(t *testing.T) {
		var saved *model.User
		repo := &mocks.UserRepoMock{
			FindByIDFn: func(context.Context, string) (*model.User, error) { return current(), nil },
			UpdateFn: func(_ context.Context, u *model.User) error {
				saved = u
				return nil
			},
		}

		resp, err := service.NewUserService(repo).Update(ctx, userID, service.UpdateUserInput{Name: ptr("Maria Silva")})

		require.NoError(t, err)
		assert.Equal(t, "Maria Silva", resp.Name)
		assert.Equal(t, "https://img.exemplo.com/maria.jpg", saved.ImageProfile)
	})

	t.Run("foto vazia remove a foto", func(t *testing.T) {
		repo := &mocks.UserRepoMock{
			FindByIDFn: func(context.Context, string) (*model.User, error) { return current(), nil },
			UpdateFn:   func(context.Context, *model.User) error { return nil },
		}

		resp, err := service.NewUserService(repo).Update(ctx, userID, service.UpdateUserInput{ImageProfile: ptr("")})

		require.NoError(t, err)
		assert.Empty(t, resp.ImageProfile)
		assert.Equal(t, "Maria", resp.Name)
	})

	t.Run("nome vazio é recusado", func(t *testing.T) {
		repo := &mocks.UserRepoMock{FindByIDFn: func(context.Context, string) (*model.User, error) { return current(), nil }}

		_, err := service.NewUserService(repo).Update(ctx, userID, service.UpdateUserInput{Name: ptr("")})

		assert.ErrorIs(t, err, service.ErrEmptyName)
	})

	t.Run("conta apagada é 404", func(t *testing.T) {
		repo := &mocks.UserRepoMock{FindByIDFn: func(context.Context, string) (*model.User, error) { return nil, gorm.ErrRecordNotFound }}

		_, err := service.NewUserService(repo).Update(ctx, userID, service.UpdateUserInput{Name: ptr("Maria")})

		assert.ErrorIs(t, err, service.ErrUserNotFound)
	})
}

func TestUserService_Delete(t *testing.T) {
	t.Run("apaga a conta", func(t *testing.T) {
		var deleted string
		repo := &mocks.UserRepoMock{
			FindByIDFn: func(context.Context, string) (*model.User, error) { return &model.User{ID: userID}, nil },
			DeleteFn: func(_ context.Context, id string) error {
				deleted = id
				return nil
			},
		}

		require.NoError(t, service.NewUserService(repo).Delete(ctx, userID))
		assert.Equal(t, userID, deleted)
	})

	t.Run("conta que não existe é 404", func(t *testing.T) {
		repo := &mocks.UserRepoMock{FindByIDFn: func(context.Context, string) (*model.User, error) { return nil, gorm.ErrRecordNotFound }}

		assert.ErrorIs(t, service.NewUserService(repo).Delete(ctx, userID), service.ErrUserNotFound)
	})
}

func TestInputs_Normalize(t *testing.T) {
	user := service.CreateUserInput{Name: "  Maria  ", Email: "  Maria@Email.COM ", ImageProfile: " https://img.exemplo.com/a.jpg "}
	user.Normalize()
	assert.Equal(t, service.CreateUserInput{Name: "Maria", Email: "maria@email.com", ImageProfile: "https://img.exemplo.com/a.jpg"}, user)

	login := service.LoginInput{Email: " MARIA@email.com"}
	login.Normalize()
	assert.Equal(t, "maria@email.com", login.Email)

	query := service.ListRecipesQuery{Search: "  pão   de   queijo ", UserID: " " + userID + " "}
	query.Normalize()
	assert.Equal(t, "pão de queijo", query.Search)
	assert.Equal(t, userID, query.UserID)

	recipe := service.RecipeInput{
		Name:         " Bolo ",
		Ingredients:  []service.IngredientRecipeInput{{Name: " Farinha ", Amount: " 2 xícaras "}},
		Preparations: []service.PreparationInput{{Description: " Misture. "}},
		Images:       []service.ImageRecipeInput{{URL: " https://img.exemplo.com/b.jpg "}},
	}
	recipe.Normalize()
	assert.Equal(t, "Bolo", recipe.Name)
	assert.Equal(t, service.IngredientRecipeInput{Name: "Farinha", Amount: "2 xícaras"}, recipe.Ingredients[0])
	assert.Equal(t, "Misture.", recipe.Preparations[0].Description)
	assert.Equal(t, "https://img.exemplo.com/b.jpg", recipe.Images[0].URL)
}
