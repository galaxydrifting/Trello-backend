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
