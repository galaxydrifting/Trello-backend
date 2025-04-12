package routes

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"trello-backend/internal/config"
)

type Router struct {
	engine    *gin.Engine
	handlers  map[string]interface{}
	jwtSecret string
	config    *config.Config
}

func NewRouter(engine *gin.Engine, jwtSecret string, cfg *config.Config) *Router {
	// 設定 CORS
	corsConfig := cors.DefaultConfig()
	corsConfig.AllowOrigins = cfg.CORSAllowOrigins
	corsConfig.AllowMethods = []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"}
	corsConfig.AllowHeaders = []string{"Origin", "Content-Type", "Accept", "Authorization"}
	
	engine.Use(cors.New(corsConfig))
	
	return &Router{
		engine:    engine,
		handlers:  make(map[string]interface{}),
		jwtSecret: jwtSecret,
		config:    cfg,
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
