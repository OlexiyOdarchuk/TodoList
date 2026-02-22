package service

import (
	"errors"
	"strings"
	"time"
	"todolist/internal/models"
	"todolist/internal/utils"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type UserRepository interface {
	GetByID(id string) (*models.User, error)
	GetByEmail(email string) (*models.User, error)
	GetByUsername(username string) (*models.User, error)
	GetByOAuth(provider, oauthID string) (*models.User, error)
	Create(u *models.User) error
	Update(u *models.User) error
	Delete(id string) error
}

type UserService struct {
	repo       UserRepository
	jwtManager *utils.JWTManager
}

type LoginResponse struct {
	Message string       `json:"message"`
	User    *models.User `json:"user"`
	Token   string       `json:"token"`
}

func NewUserService(repo UserRepository, jwtManager *utils.JWTManager) *UserService {
	return &UserService{repo: repo, jwtManager: jwtManager}
}

func (s *UserService) RegisterUser(username, email, password string) (*LoginResponse, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), 10)
	if err != nil {
		return nil, err
	}

	code, err := utils.GenerateVerificationCode()
	if err != nil {
		return nil, err
	}
	expires := time.Now().Add(10 * time.Minute)
	user := &models.User{
		Id:                      uuid.New().String(),
		Username:                username,
		Email:                   email,
		PasswordHash:            string(hashedPassword),
		VerificationCode:        code,
		VerificationCodeExpires: expires,
	}

	err = s.repo.Create(user)
	if err != nil {
		if strings.Contains(err.Error(), "unique constraint") {
			return nil, errors.New("username or email already exists")
		}
		return nil, err
	}

	go utils.SendVerificationEmail(user.Email, code)

	return &LoginResponse{Message: "User registered successfully. Check email for verification", Token: ""}, nil
}
