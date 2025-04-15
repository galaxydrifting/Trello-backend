package services

import (
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"golang.org/x/crypto/bcrypt"

	"trello-backend/internal/models"
)

type MockUserRepository struct {
	mock.Mock
}

func (m *MockUserRepository) Create(user *models.User) error {
	args := m.Called(user)
	return args.Error(0)
}

func (m *MockUserRepository) FindByEmail(email string) (*models.User, error) {
	args := m.Called(email)
	return args.Get(0).(*models.User), args.Error(1)
}

func (m *MockUserRepository) FindByID(id uuid.UUID) (*models.User, error) {
	args := m.Called(id)
	return args.Get(0).(*models.User), args.Error(1)
}

func (m *MockUserRepository) UpdatePassword(id uuid.UUID, password string) error {
	args := m.Called(id, password)
	return args.Error(0)
}

func TestAuthService_Register(t *testing.T) {
	mockRepo := new(MockUserRepository)
	jwtSecret := "testsecret"
	authService := NewAuthService(mockRepo, jwtSecret)

	req := models.RegisterRequest{
		Email:    "test@example.com",
		Name:     "Test User",
		Password: "password123",
	}

	mockRepo.On("Create", mock.Anything).Return(nil)

	token, err := authService.Register(req)

	assert.NoError(t, err)
	assert.NotEmpty(t, token)
	mockRepo.AssertExpectations(t)
}

func TestAuthService_Login(t *testing.T) {
	mockRepo := new(MockUserRepository)
	jwtSecret := "testsecret"
	authService := NewAuthService(mockRepo, jwtSecret)

	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte("password123"), bcrypt.DefaultCost)
	user := &models.User{
		ID:           uuid.New(),
		Email:        "test@example.com",
		PasswordHash: string(hashedPassword),
	}

	mockRepo.On("FindByEmail", "test@example.com").Return(user, nil)

	req := models.LoginRequest{
		Email:    "test@example.com",
		Password: "password123",
	}

	token, err := authService.Login(req)

	assert.NoError(t, err)
	assert.NotEmpty(t, token)
	mockRepo.AssertExpectations(t)
}

func TestAuthService_ChangePassword(t *testing.T) {
	mockRepo := new(MockUserRepository)
	jwtSecret := "testsecret"
	authService := NewAuthService(mockRepo, jwtSecret)

	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte("oldpassword"), bcrypt.DefaultCost)
	user := &models.User{
		ID:           uuid.New(),
		PasswordHash: string(hashedPassword),
	}

	mockRepo.On("FindByID", user.ID).Return(user, nil)
	mockRepo.On("UpdatePassword", user.ID, mock.Anything).Return(nil)

	req := models.ChangePasswordRequest{
		OldPassword: "oldpassword",
		NewPassword: "newpassword",
	}

	err := authService.ChangePassword(user.ID, req)

	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
}
