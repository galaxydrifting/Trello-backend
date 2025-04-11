package handlers

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
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
