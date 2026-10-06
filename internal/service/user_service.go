package service

import (
	"errors"
	"panda-cooking-go-api/internal/model"
	"panda-cooking-go-api/internal/repository"
	"panda-cooking-go-api/pkg/token"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type UserService struct {
	repo      repository.UserRepo
	secretKey string
}

func NewUserService(repo repository.UserRepo, secretKey string) *UserService {
	return &UserService{repo: repo, secretKey: secretKey}
}

// --- DTOs ---

type CreateUserInput struct {
	Name         string `json:"name" binding:"required"`
	Email        string `json:"email" binding:"required,email"`
	Password     string `json:"password" binding:"required,min=6"`
	ImageProfile string `json:"image_profile"`
}

type UpdateUserInput struct {
	Name         string `json:"name"`
	ImageProfile string `json:"image_profile"`
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
		return nil, err
	}

	return toUserResponse(user), nil
}

func (s *UserService) Login(input LoginInput) (*LoginResponse, error) {
	user, err := s.repo.FindByEmail(input.Email)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("email ou senha inválidos")
		}
		return nil, err
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(input.Password)); err != nil {
		return nil, errors.New("email ou senha inválidos")
	}

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
			return nil, errors.New("usuário não encontrado")
		}
		return nil, err
	}
	return toUserResponse(user), nil
}

func (s *UserService) Update(userID string, input UpdateUserInput) (*UserResponse, error) {
	user, err := s.repo.FindByID(userID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("usuário não encontrado")
		}
		return nil, err
	}

	if input.Name != "" {
		user.Name = input.Name
	}
	if input.ImageProfile != "" {
		user.ImageProfile = input.ImageProfile
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
			return errors.New("usuário não encontrado")
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
