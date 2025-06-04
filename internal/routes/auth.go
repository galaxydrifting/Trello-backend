package routes

import (
	"trello-backend/internal/handlers"
	"trello-backend/internal/middlewares"

	"github.com/gin-gonic/gin"
)

func (r *Router) setupAuthRoutes(api *gin.RouterGroup) {
	authHandler := r.handlers["auth"].(*handlers.AuthHandler)

	// 認證相關路由群組
	auth := api.Group("/auth")
	{
		// 公開的認證端點
		public := auth.Group("")
		{
			public.POST("/register", authHandler.Register)
			public.POST("/login", authHandler.Login)
			public.GET("/ping", authHandler.Ping)
		}

		// 需要認證的端點
		protected := auth.Group("")
		protected.Use(middlewares.AuthMiddleware(r.jwtSecret))
		{
			protected.POST("/change-password", authHandler.ChangePassword)
			protected.GET("/me", authHandler.GetProfile)
		}
	}
}
