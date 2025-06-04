package services

import (
	"errors"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"

	"trello-backend/internal/models"
	"trello-backend/internal/repositories"
	"trello-backend/pkg/utils"
)

type AuthService interface {
	Register(req models.RegisterRequest) (models.AuthResponse, error)
	Login(req models.LoginRequest) (models.AuthResponse, error)
	ChangePassword(userID uuid.UUID, req models.ChangePasswordRequest) error
	GetProfile(userID uuid.UUID) (models.UserProfileResponse, error)
}

type authService struct {
	userRepo  repositories.UserRepository
	jwtSecret string
}

func NewAuthService(userRepo repositories.UserRepository, jwtSecret string) AuthService {
	return &authService{
		userRepo:  userRepo,
		jwtSecret: jwtSecret,
	}
}

func (s *authService) Register(req models.RegisterRequest) (models.AuthResponse, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return models.AuthResponse{}, errors.New("密碼加密失敗")
	}

	user := models.User{
		ID:           uuid.New(),
		Email:        req.Email,
		Name:         req.Name,
		PasswordHash: string(hashedPassword),
	}

	if err := s.userRepo.Create(&user); err != nil {
		return models.AuthResponse{}, errors.New("使用者建立失敗")
	}

	token, err := utils.GenerateToken(user.ID, s.jwtSecret)
	if err != nil {
		return models.AuthResponse{}, errors.New("Token 產生失敗")
	}

	return models.AuthResponse{
		Token: token,
		Name:  user.Name,
		Email: user.Email,
	}, nil
}

func (s *authService) Login(req models.LoginRequest) (models.AuthResponse, error) {
	user, err := s.userRepo.FindByEmail(req.Email)
	if err != nil {
		return models.AuthResponse{}, errors.New("帳號或密碼錯誤")
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(req.Password)); err != nil {
		return models.AuthResponse{}, errors.New("帳號或密碼錯誤")
	}

	token, err := utils.GenerateToken(user.ID, s.jwtSecret)
	if err != nil {
		return models.AuthResponse{}, errors.New("Token 產生失敗")
	}

	return models.AuthResponse{
		Token: token,
		Name:  user.Name,
		Email: user.Email,
	}, nil
}

func (s *authService) ChangePassword(userID uuid.UUID, req models.ChangePasswordRequest) error {
	user, err := s.userRepo.FindByID(userID)
	if err != nil {
		return errors.New("使用者不存在")
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(req.OldPassword)); err != nil {
		return errors.New("舊密碼錯誤")
	}

	newHashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.NewPassword), bcrypt.DefaultCost)
	if err != nil {
		return errors.New("密碼加密失敗")
	}

	return s.userRepo.UpdatePassword(userID, string(newHashedPassword))
}

func (s *authService) GetProfile(userID uuid.UUID) (models.UserProfileResponse, error) {
	user, err := s.userRepo.FindByID(userID)
	if err != nil {
		return models.UserProfileResponse{}, errors.New("使用者不存在")
	}
	return models.UserProfileResponse{
		Name:  user.Name,
		Email: user.Email,
	}, nil
}
