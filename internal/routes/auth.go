package routes

import "trello-backend/internal/handlers"

func (r *Router) setupAuthRoutes() {
	handler := r.handlers["auth"].(*handlers.AuthHandler)

	// 公開路由
	r.engine.POST("/register", handler.Register)
	r.engine.POST("/login", handler.Login)
}
