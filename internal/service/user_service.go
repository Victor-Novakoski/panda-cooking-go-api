package service

import (
	"errors"
	"strings"
	"time"

	"panda-cooking-go-api/internal/model"
	"panda-cooking-go-api/internal/ratelimit"
	"panda-cooking-go-api/internal/repository"
	"panda-cooking-go-api/pkg/token"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

// Depois de MaxLoginFailures senhas erradas para o mesmo e-mail, o login
// dele fica bloqueado até a janela de LoginLockout acabar.
const (
	MaxLoginFailures = 5
	LoginLockout     = 15 * time.Minute
)

type UserService struct {
	repo          repository.UserRepo
	secretKey     string
	loginFailures *ratelimit.Limiter
}

func NewUserService(repo repository.UserRepo, secretKey string) *UserService {
	return &UserService{
		repo:          repo,
		secretKey:     secretKey,
		loginFailures: ratelimit.New(MaxLoginFailures, LoginLockout),
	}
}

// --- DTOs ---

type CreateUserInput struct {
	Name         string `json:"name" binding:"required"`
	Email        string `json:"email" binding:"required,email"`
	Password     string `json:"password" binding:"required,min=6"`
	ImageProfile string `json:"image_profile"`
}

// UpdateUserInput usa ponteiro para separar "não mandou" (nil, fica como
// está) de "mandou vazio": foto vazia remove a foto; nome vazio é recusado.
type UpdateUserInput struct {
	Name         *string `json:"name" binding:"omitempty,min=1"`
	ImageProfile *string `json:"image_profile" binding:"omitempty,len=0|url"`
}

type UserResponse struct {
	ID           string `json:"id"`
	Name         string `json:"name"`
	Email        string `json:"email"`
	ImageProfile string `json:"image_profile"`
	IsAdm        bool   `json:"is_adm"`
}

type LoginInput struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

type LoginResponse struct {
	Token string `json:"token"`
}

// --- Métodos ---

func (s *UserService) Create(input CreateUserInput) (*UserResponse, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	user := &model.User{
		Name:         input.Name,
		Email:        input.Email,
		Password:     string(hash),
		ImageProfile: input.ImageProfile,
	}

	if err := s.repo.Create(user); err != nil {
		if errors.Is(err, gorm.ErrDuplicatedKey) {
			return nil, ErrEmailTaken
		}
		return nil, err
	}

	return toUserResponse(user), nil
}

func (s *UserService) Login(input LoginInput) (*LoginResponse, error) {
	key := strings.ToLower(strings.TrimSpace(input.Email))
	if blocked, _ := s.loginFailures.Blocked(key); blocked {
		return nil, ErrTooManyAttempts
	}

	user, err := s.repo.FindByEmail(input.Email)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			s.loginFailures.Allow(key)
			return nil, ErrInvalidCredentials
		}
		return nil, err
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(input.Password)); err != nil {
		s.loginFailures.Allow(key)
		return nil, ErrInvalidCredentials
	}
	s.loginFailures.Reset(key)

	t, err := token.Generate(user.ID, user.IsAdm, s.secretKey)
	if err != nil {
		return nil, err
	}

	return &LoginResponse{Token: t}, nil
}

func (s *UserService) GetProfile(userID string) (*UserResponse, error) {
	user, err := s.repo.FindByID(userID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrUserNotFound
		}
		return nil, err
	}
	return toUserResponse(user), nil
}

func (s *UserService) Update(userID string, input UpdateUserInput) (*UserResponse, error) {
	user, err := s.repo.FindByID(userID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrUserNotFound
		}
		return nil, err
	}

	if input.Name != nil {
		name := strings.TrimSpace(*input.Name)
		if name == "" {
			return nil, ErrEmptyName
		}
		user.Name = name
	}
	if input.ImageProfile != nil {
		user.ImageProfile = strings.TrimSpace(*input.ImageProfile)
	}

	if err := s.repo.Update(user); err != nil {
		return nil, err
	}

	return toUserResponse(user), nil
}

func (s *UserService) Delete(userID string) error {
	_, err := s.repo.FindByID(userID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return ErrUserNotFound
		}
		return err
	}
	return s.repo.Delete(userID)
}

func (s *UserService) GetFavoriteRecipes(userID string) ([]FavoriteRecipeSummary, error) {
	favorites, err := s.repo.FindFavoriteRecipes(userID)
	if err != nil {
		return nil, err
	}

	result := make([]FavoriteRecipeSummary, len(favorites))
	for i, f := range favorites {
		result[i] = FavoriteRecipeSummary{
			ID:       f.ID,
			RecipeID: f.RecipeID,
			Name:     f.Recipe.Name,
			Time:     f.Recipe.Time,
			Portions: f.Recipe.Portions,
		}
	}
	return result, nil
}

type FavoriteRecipeSummary struct {
	ID       uint   `json:"id"`
	RecipeID string `json:"recipe_id"`
	Name     string `json:"name"`
	Time     string `json:"time"`
	Portions int    `json:"portions"`
}

// toUserResponse converte model → DTO de resposta, nunca expondo a senha
func toUserResponse(u *model.User) *UserResponse {
	return &UserResponse{
		ID:           u.ID,
		Name:         u.Name,
		Email:        u.Email,
		ImageProfile: u.ImageProfile,
		IsAdm:        u.IsAdm,
	}
}
