//go:build wireinject
// +build wireinject

package app

import (
	"github.com/google/wire"
	"gorm.io/gorm"

	"trello-backend/internal/handlers"
	"trello-backend/internal/repository/postgres"
	"trello-backend/internal/services"
)

var userRepositorySet = wire.NewSet(
	postgres.NewUserRepository,
)

var serviceSet = wire.NewSet(
	services.NewAuthService,
)

var handlerSet = wire.NewSet(
	handlers.NewAuthHandler,
)

// InitializeAPI 初始化 API 相依性
func InitializeAPI(db *gorm.DB, jwtSecret string) (*handlers.AuthHandler, error) {
	wire.Build(
		userRepositorySet,
		serviceSet,
		handlerSet,
	)
	return nil, nil
}

// 初始化所有 handlers
func InitializeHandlers(db *gorm.DB, jwtSecret string) (*handlers.Handlers, error) {
	wire.Build(
		// 使用者相關
		postgres.NewUserRepository,
		services.NewAuthService,
		handlers.NewAuthHandler,

		// TODO: 在這裡新增其他服務的相依性注入
		// services.NewBoardService,
		// handlers.NewBoardHandler,
		// ... 其他服務

		// 最後注入 Handlers 結構體
		handlers.NewHandlers,
	)
	return &handlers.Handlers{}, nil
}
