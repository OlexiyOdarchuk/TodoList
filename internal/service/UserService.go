package service

import (
	"context"
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

func (s *UserService) VerifyEmail(email, code string) (*LoginResponse, error) {
	user, err := s.repo.GetByEmail(email)
	if err != nil {
		return nil, errors.New("invalid verification code or email")
	}

	if user.IsVerified {
		return nil, errors.New("email already verified")
	}

	if user.VerificationCode != code {
		return nil, errors.New("invalid verification code or email")
	}

	if time.Now().After(user.VerificationCodeExpires) {
		return nil, errors.New("verification code expired")
	}

	user.VerificationCode = ""
	user.IsVerified = true
	user.VerificationCodeExpires = time.Time{}
	err = s.repo.Update(user)
	if err != nil {
		return nil, errors.New("failed to update user")
	}

	token, err := s.jwtManager.Generate(user.Id)
	if err != nil {
		return nil, err
	}

	return &LoginResponse{
		Token: token,
		User:  user,
	}, nil
}

func (s *UserService) LoginUser(userdata, password string) (*LoginResponse, error) {
	user := &models.User{}
	var err error

	hasAt := strings.Contains(userdata, "@")
	if hasAt {
		user, err = s.repo.GetByEmail(userdata)
		if err != nil {
			return nil, err
		}
	} else {
		user, err = s.repo.GetByUsername(userdata)
		if err != nil {
			return nil, err
		}
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password))
	if err != nil {
		return nil, errors.New("invalid credentials")
	}

	if !user.IsVerified {
		return nil, errors.New("user is not verified")
	}

	token, err := s.jwtManager.Generate(user.Id)
	if err != nil {
		return nil, err
	}
	return &LoginResponse{User: user, Token: token}, nil
}

func (s *UserService) LoginWithOAuth(ctx context.Context, data models.OAuthUser) (*LoginResponse, error) {
	user, err := s.repo.GetByOAuth(data.Provider, data.ID)
	if err != nil {
		user, err = s.repo.GetByEmail(data.Email)

		if err != nil {
			user = &models.User{
				Id:            uuid.New().String(),
				Email:         data.Email,
				Username:      data.Name,
				OauthProvider: data.Provider,
				OauthId:       data.ID,
				IsVerified:    true,
			}
			err = s.repo.Create(user)
			if err != nil {
				if strings.Contains(err.Error(), "unique constraint") {
					return nil, errors.New("username or email already exists")
				}
				return nil, err
			}
		} else {
			updateUser := &models.User{
				Id:                      user.Id,
				Username:                user.Username,
				Email:                   data.Email,
				PasswordHash:            user.PasswordHash,
				IsVerified:              true,
				VerificationCode:        "",
				VerificationCodeExpires: time.Time{},
				OauthProvider:           data.Provider,
				OauthId:                 data.ID,
			}
			s.repo.Update(updateUser)
		}
	}
	token, _ := s.jwtManager.Generate(user.Id)
	return &LoginResponse{User: user, Token: token}, nil
}
