package handlers

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"trello-backend/internal/models"
)

type MockAuthService struct {
	mock.Mock
}

func (m *MockAuthService) Register(req models.RegisterRequest) (string, error) {
	args := m.Called(req)
	return args.String(0), args.Error(1)
}

func (m *MockAuthService) Login(req models.LoginRequest) (string, error) {
	args := m.Called(req)
	return args.String(0), args.Error(1)
}

func (m *MockAuthService) ChangePassword(userID uuid.UUID, req models.ChangePasswordRequest) error {
	args := m.Called(userID, req)
	return args.Error(0)
}

func TestAuthHandler_Register(t *testing.T) {
	gin.SetMode(gin.TestMode)

	mockService := new(MockAuthService)
	handler := NewAuthHandler(mockService)

	tests := []struct {
		name          string
		input         models.RegisterRequest
		mockResponse  string
		mockError     error
		expectedCode  int
		expectedToken string
	}{
		{
			name: "成功註冊",
			input: models.RegisterRequest{
				Email:    "test@example.com",
				Password: "password123",
				Name:     "Test User",
			},
			mockResponse:  "test-token",
			mockError:     nil,
			expectedCode:  http.StatusCreated,
			expectedToken: "test-token",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockService.On("Register", tt.input).Return(tt.mockResponse, tt.mockError)

			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)

			body, _ := json.Marshal(tt.input)
			c.Request = httptest.NewRequest(http.MethodPost, "/register", bytes.NewBuffer(body))
			c.Request.Header.Set("Content-Type", "application/json")

			handler.Register(c)

			assert.Equal(t, tt.expectedCode, w.Code)
			if tt.expectedCode == http.StatusCreated {
				var response models.LoginResponse
				err := json.Unmarshal(w.Body.Bytes(), &response)
				assert.NoError(t, err)
				assert.Equal(t, tt.expectedToken, response.Token)
			}

			mockService.AssertExpectations(t)
		})
	}
}

func TestAuthHandler_Login(t *testing.T) {
	gin.SetMode(gin.TestMode)

	mockService := new(MockAuthService)
	handler := NewAuthHandler(mockService)

	tests := []struct {
		name          string
		input         models.LoginRequest
		mockResponse  string
		mockError     error
		expectedCode  int
		expectedToken string
	}{
		{
			name: "成功登入",
			input: models.LoginRequest{
				Email:    "test@example.com",
				Password: "password123",
			},
			mockResponse:  "test-token",
			mockError:     nil,
			expectedCode:  http.StatusOK,
			expectedToken: "test-token",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockService.On("Login", tt.input).Return(tt.mockResponse, tt.mockError)

			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)

			body, _ := json.Marshal(tt.input)
			c.Request = httptest.NewRequest(http.MethodPost, "/login", bytes.NewBuffer(body))
			c.Request.Header.Set("Content-Type", "application/json")

			handler.Login(c)

			assert.Equal(t, tt.expectedCode, w.Code)
			if tt.expectedCode == http.StatusOK {
				var response models.LoginResponse
				err := json.Unmarshal(w.Body.Bytes(), &response)
				assert.NoError(t, err)
				assert.Equal(t, tt.expectedToken, response.Token)
			}

			mockService.AssertExpectations(t)
		})
	}
}

func TestAuthHandler_ChangePassword(t *testing.T) {
	gin.SetMode(gin.TestMode)

	mockService := new(MockAuthService)
	handler := NewAuthHandler(mockService)
	userID := uuid.New()

	tests := []struct {
		name         string
		input        models.ChangePasswordRequest
		setupAuth    bool
		mockError    error
		expectedCode int
	}{
		{
			name: "成功變更密碼",
			input: models.ChangePasswordRequest{
				OldPassword: "oldpass123",
				NewPassword: "newpass123",
			},
			setupAuth:    true,
			mockError:    nil,
			expectedCode: http.StatusOK,
		},
		{
			name: "未認證",
			input: models.ChangePasswordRequest{
				OldPassword: "oldpass123",
				NewPassword: "newpass123",
			},
			setupAuth:    false,
			mockError:    nil,
			expectedCode: http.StatusUnauthorized,
		},
		{
			name: "舊密碼錯誤",
			input: models.ChangePasswordRequest{
				OldPassword: "wrongpass",
				NewPassword: "newpass123",
			},
			setupAuth:    true,
			mockError:    errors.New("舊密碼錯誤"),
			expectedCode: http.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.setupAuth {
				mockService.On("ChangePassword", userID, tt.input).Return(tt.mockError)
			}

			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)

			if tt.setupAuth {
				c.Set("userID", userID)
			}

			body, _ := json.Marshal(tt.input)
			c.Request = httptest.NewRequest(http.MethodPost, "/api/auth/change-password", bytes.NewBuffer(body))
			c.Request.Header.Set("Content-Type", "application/json")

			handler.ChangePassword(c)

			assert.Equal(t, tt.expectedCode, w.Code)
			mockService.AssertExpectations(t)
		})
	}
}
