package routes

import (
	"trello-backend/internal/handlers"
	"trello-backend/internal/middleware"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

type Router struct {
	engine    *gin.Engine
	handlers  *handlers.Handlers
	jwtSecret string
}

func NewRouter(engine *gin.Engine, handlers *handlers.Handlers, jwtSecret string) *Router {
	return &Router{
		engine:    engine,
		handlers:  handlers,
		jwtSecret: jwtSecret,
	}
}

func (r *Router) SetupRoutes() {
	// Swagger 文件路由
	r.engine.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// 設定認證相關路由
	r.setupAuthRoutes()

	// API 路由群組
	api := r.engine.Group("/api")
	api.Use(middleware.AuthMiddleware(r.jwtSecret))

	// TODO: 在這裡加入其他 API 路由群組
	// r.setupBoardRoutes(api)
	// r.setupCardRoutes(api)
	// r.setupListRoutes(api)
	// r.setupCommentRoutes(api)
	// r.setupUserRoutes(api)
	// r.setupWorkspaceRoutes(api)
	// r.setupActivityRoutes(api)
	// r.setupLabelRoutes(api)
	// r.setupAttachmentRoutes(api)
}
