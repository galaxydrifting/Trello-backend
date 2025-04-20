//go:build wireinject
// +build wireinject

package app

import (
	"github.com/google/wire"
	"gorm.io/gorm"

	"trello-backend/graph"
	"trello-backend/internal/handlers"
	"trello-backend/internal/repositories"
	"trello-backend/internal/services"
)

// Handler 介面定義所有 handler 必須實作的方法
type Handler interface{}

// API 包含所有的 handlers
type API struct {
	handlers map[string]Handler
	BoardSvc services.BoardService
	ListSvc  services.ListService
	CardSvc  services.CardService
}

func (a *API) BoardService() services.BoardService {
	return a.BoardSvc
}

func (a *API) ListService() services.ListService {
	return a.ListSvc
}

func (a *API) CardService() services.CardService {
	return a.CardSvc
}

// NewAPI 建立新的 API 實例
func NewAPI(authHandler *handlers.AuthHandler, boardService services.BoardService, listService services.ListService, cardService services.CardService) *API {
	api := &API{
		handlers: make(map[string]Handler),
		BoardSvc: boardService,
		ListSvc:  listService,
		CardSvc:  cardService,
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

// Board/List/Card Provider Set
var boardDomainSet = wire.NewSet(
	repositories.NewBoardRepository,
	services.NewBoardService,
)
var listDomainSet = wire.NewSet(
	repositories.NewListRepository,
	services.NewListService,
)
var cardDomainSet = wire.NewSet(
	repositories.NewCardRepository,
	services.NewCardService,
)

// GraphQL Resolver Provider
var resolverSet = wire.NewSet(
	boardDomainSet,
	listDomainSet,
	cardDomainSet,
	graph.NewResolver,
)

// API Provider Set
var apiSet = wire.NewSet(
	userDomainSet,
	resolverSet,
	NewAPI,
)

// InitializeAPI 初始化 API 相依性
func InitializeAPI(db *gorm.DB, jwtSecret string) (*API, error) {
	wire.Build(apiSet)
	return nil, nil
}
