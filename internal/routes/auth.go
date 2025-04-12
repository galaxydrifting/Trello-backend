package routes

import (
	"trello-backend/internal/handlers"
	"trello-backend/internal/middleware"
)

func (r *Router) setupAuthRoutes() {
	authHandler := r.handlers["auth"].(*handlers.AuthHandler)

	auth := r.engine.Group("/auth")
	{
		auth.POST("/register", authHandler.Register)
		auth.POST("/login", authHandler.Login)
	}

	// 需要認證的路由
	authProtected := r.engine.Group("/api/auth")
	authProtected.Use(middleware.AuthMiddleware(r.jwtSecret))
	{
		authProtected.POST("/change-password", authHandler.ChangePassword)
	}
}
