//go:build wireinject
// +build wireinject

package app

import (
	"github.com/google/wire"
	"gorm.io/gorm"

	"trello-backend/internal/handlers"
	"trello-backend/internal/repositories"
	"trello-backend/internal/services"
)

// Handler 介面定義所有 handler 必須實作的方法
type Handler interface{}

// API 包含所有的 handlers
type API struct {
	handlers map[string]Handler
}

// NewAPI 建立新的 API 實例
func NewAPI(authHandler *handlers.AuthHandler) *API {
	api := &API{
		handlers: make(map[string]Handler),
	}
	api.RegisterHandler("auth", authHandler)
	return api
}

// RegisterHandler 註冊新的 handler
func (a *API) RegisterHandler(name string, handler Handler) {
	a.handlers[name] = handler
}

// GetHandlers 取得所有已註冊的 handlers
func (a *API) GetHandlers() map[string]Handler {
	return a.handlers
}

// 使用者領域的 Provider Set
var userDomainSet = wire.NewSet(
	repositories.NewUserRepository,
	services.NewAuthService,
	handlers.NewAuthHandler,
)

// API Provider Set
var apiSet = wire.NewSet(
	userDomainSet,
	NewAPI,
)

// InitializeAPI 初始化 API 相依性
func InitializeAPI(db *gorm.DB, jwtSecret string) (*API, error) {
	wire.Build(apiSet)
	return nil, nil
}
