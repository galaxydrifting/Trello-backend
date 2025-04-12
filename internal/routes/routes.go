package routes

import (
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

type Router struct {
	engine    *gin.Engine
	handlers  map[string]interface{}
	jwtSecret string
}

func NewRouter(engine *gin.Engine, jwtSecret string) *Router {
	return &Router{
		engine:    engine,
		handlers:  make(map[string]interface{}),
		jwtSecret: jwtSecret,
	}
}

// RegisterHandler 註冊新的 handler
func (r *Router) RegisterHandler(name string, handler interface{}) {
	r.handlers[name] = handler
}

func (r *Router) SetupRoutes() {
	// Swagger 文件路由
	r.engine.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// API 路由群組
	api := r.engine.Group("/api")

	// 設定各個功能模組的路由
	r.setupAuthRoutes(api)
	// 之後可以輕鬆添加其他模組的路由
	// r.setupBoardRoutes(api)
	// r.setupCardRoutes(api)
	// r.setupListRoutes(api)
}
