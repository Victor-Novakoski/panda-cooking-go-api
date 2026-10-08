package service

import (
	"context"
	"errors"
	"fmt"

	"panda-cooking-go-api/internal/model"
	"panda-cooking-go-api/internal/repository"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

// UserService cuida do cadastro e do perfil. Login e sessão ficam no
// AuthService.
type UserService struct {
	repo repository.UserRepo
}

func NewUserService(repo repository.UserRepo) *UserService {
	return &UserService{repo: repo}
}

func (s *UserService) Create(ctx context.Context, input CreateUserInput) (*UserResponse, error) {
	if err := checkPassword(input.Password, input.Email); err != nil {
		return nil, err
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, fmt.Errorf("gerar hash da senha: %w", err)
	}

	user := &model.User{
		Name:         input.Name,
		Email:        input.Email,
		PasswordHash: string(hash),
		ImageProfile: input.ImageProfile,
	}
	if err := s.repo.Create(ctx, user); err != nil {
		if errors.Is(err, gorm.ErrDuplicatedKey) {
			return nil, ErrEmailTaken
		}
		return nil, err
	}

	resp := toUserResponse(user)
	return &resp, nil
}

func (s *UserService) GetProfile(ctx context.Context, userID string) (*UserResponse, error) {
	user, err := s.find(ctx, userID)
	if err != nil {
		return nil, err
	}
	resp := toUserResponse(user)
	return &resp, nil
}

func (s *UserService) Update(ctx context.Context, userID string, input UpdateUserInput) (*UserResponse, error) {
	user, err := s.find(ctx, userID)
	if err != nil {
		return nil, err
	}

	if input.Name != nil {
		if *input.Name == "" {
			return nil, ErrEmptyName
		}
		user.Name = *input.Name
	}
	if input.ImageProfile != nil {
		user.ImageProfile = *input.ImageProfile
	}

	if err := s.repo.Update(ctx, user); err != nil {
		return nil, err
	}

	resp := toUserResponse(user)
	return &resp, nil
}

// Delete apaga a conta com tudo o que ela criou (receitas, comentários,
// favoritos e sessões).
func (s *UserService) Delete(ctx context.Context, userID string) error {
	if _, err := s.find(ctx, userID); err != nil {
		return err
	}
	return s.repo.Delete(ctx, userID)
}

func (s *UserService) find(ctx context.Context, userID string) (*model.User, error) {
	user, err := s.repo.FindByID(ctx, userID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrUserNotFound
		}
		return nil, err
	}
	return user, nil
}
