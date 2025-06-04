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
// Deprecated: 請改用 AuthResponse
// type LoginResponse struct {
// 	Token string `json:"token" example:"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."`
// }

// AuthResponse 登入/註冊回應
// 回傳 token、name、email
// swagger:model
type AuthResponse struct {
	Token string `json:"token" example:"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."`
	Name  string `json:"name" example:"王小明"`
	Email string `json:"email" example:"user@example.com"`
}

// ChangePasswordRequest 變更密碼請求
type ChangePasswordRequest struct {
	OldPassword string `json:"oldPassword" binding:"required" example:"oldpass123"`
	NewPassword string `json:"newPassword" binding:"required" example:"newpass123"`
}

// UserProfileResponse 取得個人資訊回應
// swagger:model
// 用於 /auth/me
// 只回傳 name, email
type UserProfileResponse struct {
	Name  string `json:"name" example:"王小明"`
	Email string `json:"email" example:"user@example.com"`
}
