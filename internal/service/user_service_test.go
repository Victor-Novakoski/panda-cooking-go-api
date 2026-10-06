package service_test

import (
	"errors"
	"testing"

	"panda-cooking-go-api/internal/model"
	"panda-cooking-go-api/internal/service"
	"panda-cooking-go-api/internal/service/mocks"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

const testSecret = "test-secret"

func TestUserService_Create(t *testing.T) {
	t.Run("cria usuário com sucesso", func(t *testing.T) {
		repo := &mocks.UserRepoMock{
			CreateFn: func(user *model.User) error {
				user.ID = "uuid-gerado"
				return nil
			},
		}
		svc := service.NewUserService(repo, testSecret)

		resp, err := svc.Create(service.CreateUserInput{
			Name:     "Victor",
			Email:    "victor@email.com",
			Password: "123456",
		})

		require.NoError(t, err)
		assert.Equal(t, "Victor", resp.Name)
		assert.Equal(t, "victor@email.com", resp.Email)
		assert.Equal(t, "uuid-gerado", resp.ID)
	})

	t.Run("retorna erro se o banco falhar", func(t *testing.T) {
		repo := &mocks.UserRepoMock{
			CreateFn: func(user *model.User) error {
				return errors.New("duplicate key")
			},
		}
		svc := service.NewUserService(repo, testSecret)

		_, err := svc.Create(service.CreateUserInput{
			Name:     "Victor",
			Email:    "victor@email.com",
			Password: "123456",
		})

		require.Error(t, err)
	})

	t.Run("e-mail já cadastrado retorna conflito", func(t *testing.T) {
		repo := &mocks.UserRepoMock{
			CreateFn: func(user *model.User) error { return gorm.ErrDuplicatedKey },
		}
		svc := service.NewUserService(repo, testSecret)

		_, err := svc.Create(service.CreateUserInput{Name: "Victor", Email: "victor@email.com", Password: "123456"})

		assert.ErrorIs(t, err, service.ErrEmailTaken)
	})
}

func TestUserService_Login(t *testing.T) {
	hash, _ := bcrypt.GenerateFromPassword([]byte("senha-correta"), bcrypt.DefaultCost)

	t.Run("login com sucesso retorna token", func(t *testing.T) {
		repo := &mocks.UserRepoMock{
			FindByEmailFn: func(email string) (*model.User, error) {
				return &model.User{ID: "abc", Email: email, Password: string(hash)}, nil
			},
		}
		svc := service.NewUserService(repo, testSecret)

		resp, err := svc.Login(service.LoginInput{
			Email:    "victor@email.com",
			Password: "senha-correta",
		})

		require.NoError(t, err)
		assert.NotEmpty(t, resp.Token)
	})

	t.Run("senha errada retorna erro genérico", func(t *testing.T) {
		repo := &mocks.UserRepoMock{
			FindByEmailFn: func(email string) (*model.User, error) {
				return &model.User{Password: string(hash)}, nil
			},
		}
		svc := service.NewUserService(repo, testSecret)

		_, err := svc.Login(service.LoginInput{Email: "v@v.com", Password: "senha-errada"})

		require.Error(t, err)
		assert.Equal(t, "email ou senha inválidos", err.Error())
	})

	t.Run("email não encontrado retorna erro genérico", func(t *testing.T) {
		repo := &mocks.UserRepoMock{
			FindByEmailFn: func(email string) (*model.User, error) {
				return nil, gorm.ErrRecordNotFound
			},
		}
		svc := service.NewUserService(repo, testSecret)

		_, err := svc.Login(service.LoginInput{Email: "naoexiste@email.com", Password: "qualquer"})

		require.Error(t, err)
		// mesma mensagem — não revela se o email existe
		assert.Equal(t, "email ou senha inválidos", err.Error())
	})
}

func TestUserService_GetProfile(t *testing.T) {
	t.Run("retorna perfil do usuário", func(t *testing.T) {
		repo := &mocks.UserRepoMock{
			FindByIDFn: func(id string) (*model.User, error) {
				return &model.User{ID: id, Name: "Victor", Email: "v@v.com"}, nil
			},
		}
		svc := service.NewUserService(repo, testSecret)

		resp, err := svc.GetProfile("meu-id")

		require.NoError(t, err)
		assert.Equal(t, "meu-id", resp.ID)
		assert.Equal(t, "Victor", resp.Name)
	})

	t.Run("usuário não encontrado retorna erro", func(t *testing.T) {
		repo := &mocks.UserRepoMock{
			FindByIDFn: func(id string) (*model.User, error) {
				return nil, gorm.ErrRecordNotFound
			},
		}
		svc := service.NewUserService(repo, testSecret)

		_, err := svc.GetProfile("id-inexistente")

		require.Error(t, err)
		assert.Equal(t, "usuário não encontrado", err.Error())
	})
}

func TestUserService_Delete(t *testing.T) {
	t.Run("deleta usuário existente", func(t *testing.T) {
		repo := &mocks.UserRepoMock{
			FindByIDFn: func(id string) (*model.User, error) {
				return &model.User{ID: id}, nil
			},
			DeleteFn: func(id string) error { return nil },
		}
		svc := service.NewUserService(repo, testSecret)

		err := svc.Delete("meu-id")
		require.NoError(t, err)
	})

	t.Run("não encontrado retorna erro", func(t *testing.T) {
		repo := &mocks.UserRepoMock{
			FindByIDFn: func(id string) (*model.User, error) {
				return nil, gorm.ErrRecordNotFound
			},
		}
		svc := service.NewUserService(repo, testSecret)

		err := svc.Delete("id-inexistente")
		require.Error(t, err)
		assert.Equal(t, "usuário não encontrado", err.Error())
	})
}

func TestUserService_LoginLockout(t *testing.T) {
	hash, _ := bcrypt.GenerateFromPassword([]byte("senha-correta"), bcrypt.MinCost)
	newRepo := func() *mocks.UserRepoMock {
		return &mocks.UserRepoMock{
			FindByEmailFn: func(email string) (*model.User, error) {
				return &model.User{ID: "abc", Email: email, Password: string(hash)}, nil
			},
		}
	}
	fail := func(svc *service.UserService, email string) error {
		_, err := svc.Login(service.LoginInput{Email: email, Password: "senha-errada"})
		return err
	}

	t.Run("bloqueia o e-mail depois de várias senhas erradas, mesmo com a senha certa", func(t *testing.T) {
		svc := service.NewUserService(newRepo(), testSecret)
		for range service.MaxLoginFailures {
			assert.ErrorIs(t, fail(svc, "victor@email.com"), service.ErrInvalidCredentials)
		}

		_, err := svc.Login(service.LoginInput{Email: "victor@email.com", Password: "senha-correta"})

		assert.ErrorIs(t, err, service.ErrTooManyAttempts)
	})

	t.Run("maiúsculas e espaços no e-mail não burlam o bloqueio", func(t *testing.T) {
		svc := service.NewUserService(newRepo(), testSecret)
		for i := range service.MaxLoginFailures {
			if i%2 == 0 {
				_ = fail(svc, " Victor@Email.com ")
			} else {
				_ = fail(svc, "victor@email.com")
			}
		}

		assert.ErrorIs(t, fail(svc, "VICTOR@EMAIL.COM"), service.ErrTooManyAttempts)
	})

	t.Run("e-mail inexistente também conta tentativa", func(t *testing.T) {
		repo := &mocks.UserRepoMock{
			FindByEmailFn: func(string) (*model.User, error) { return nil, gorm.ErrRecordNotFound },
		}
		svc := service.NewUserService(repo, testSecret)
		for range service.MaxLoginFailures {
			_ = fail(svc, "naoexiste@email.com")
		}

		assert.ErrorIs(t, fail(svc, "naoexiste@email.com"), service.ErrTooManyAttempts)
	})

	t.Run("login certo zera as tentativas", func(t *testing.T) {
		svc := service.NewUserService(newRepo(), testSecret)
		for range service.MaxLoginFailures - 1 {
			_ = fail(svc, "victor@email.com")
		}
		_, err := svc.Login(service.LoginInput{Email: "victor@email.com", Password: "senha-correta"})
		require.NoError(t, err)

		assert.ErrorIs(t, fail(svc, "victor@email.com"), service.ErrInvalidCredentials)
	})

	t.Run("bloqueio de um e-mail não afeta outro", func(t *testing.T) {
		svc := service.NewUserService(newRepo(), testSecret)
		for range service.MaxLoginFailures {
			_ = fail(svc, "victor@email.com")
		}

		_, err := svc.Login(service.LoginInput{Email: "outra@email.com", Password: "senha-correta"})

		require.NoError(t, err)
	})
}
