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

// 使用者領域的 Provider Set
var userDomainSet = wire.NewSet(
	postgres.NewUserRepository,
	services.NewAuthService,
	handlers.NewAuthHandler,
)

// 所有領域的 Provider Set 組合
var allDomainsSet = wire.NewSet(
	userDomainSet,
)

// InitializeAPI 初始化 API 相依性
func InitializeAPI(db *gorm.DB, jwtSecret string) (*handlers.AuthHandler, error) {
	wire.Build(userDomainSet)
	return nil, nil
}
