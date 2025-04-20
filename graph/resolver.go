package graph

import (
	"trello-backend/internal/services"
)

// This file will not be regenerated automatically.
//
// It serves as dependency injection for your app, add any dependencies you require here.

type Resolver struct {
	BoardService *services.BoardService
	ListService  *services.ListService
	CardService  *services.CardService
}
