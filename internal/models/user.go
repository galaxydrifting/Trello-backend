package models

import (
	"time"

	"github.com/google/uuid"
)

type User struct {
	ID           uuid.UUID `gorm:"type:uuid;primary_key"`
	Email        string    `gorm:"unique;not null"`
	Name         string    `gorm:"not null"`
	PasswordHash string    `gorm:"not null"`
	CreatedAt    time.Time `gorm:"default:CURRENT_TIMESTAMP"`
}

// APIResponse 定義通用的 API 回應格式
type APIResponse struct {
	Error string `json:"error,omitempty" example:"錯誤訊息"`
}

// RegisterRequest 註冊請求
type RegisterRequest struct {
	Email    string `json:"email" binding:"required,email" example:"user@example.com"`
	Name     string `json:"name" binding:"required" example:"王小明"`
	Password string `json:"password" binding:"required,min=6" example:"password123"`
}

// LoginRequest 登入請求
type LoginRequest struct {
	Email    string `json:"email" binding:"required,email" example:"user@example.com"`
	Password string `json:"password" binding:"required" example:"password123"`
}

// LoginResponse 登入回應
type LoginResponse struct {
	Token string `json:"token" example:"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."`
}
