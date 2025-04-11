package main

import (
	"log"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	"trello-backend/docs"
	"trello-backend/internal/app"
	"trello-backend/internal/config"
	"trello-backend/internal/middleware"
	"trello-backend/internal/models"
)

// @title Trello 後端 API
// @version 1.0
// @description Trello 後端 API 文件
// @host localhost:8080
// @BasePath /
// @schemes http
func initDB(cfg *config.Config) *gorm.DB {
	db, err := gorm.Open(postgres.Open(cfg.GetDBConnString()), &gorm.Config{})
	if err != nil {
		log.Fatal("無法連線到資料庫:", err)
	}

	err = db.AutoMigrate(&models.User{})
	if err != nil {
		log.Fatal("資料庫遷移失敗:", err)
	}
	log.Println("資料庫遷移完成")

	return db
}

func main() {
	cfg := config.LoadConfig()
	db := initDB(cfg)

	// 初始化 Swagger 文件
	docs.SwaggerInfo.Title = "Trello 後端 API"
	docs.SwaggerInfo.Description = "Trello 後端 API 文件"
	docs.SwaggerInfo.Version = "1.0"
	docs.SwaggerInfo.Host = "localhost:8080"
	docs.SwaggerInfo.BasePath = "/"
	docs.SwaggerInfo.Schemes = []string{"http"}

	// 使用 wire 進行相依性注入
	authHandler, err := app.InitializeAPI(db, cfg.JWTSecret)
	if err != nil {
		log.Fatal("無法初始化 API:", err)
	}

	// 設定路由
	r := gin.Default()

	// Swagger 文件路由
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// 公開路由
	r.POST("/register", authHandler.Register)
	r.POST("/login", authHandler.Login)

	// 需要認證的路由
	auth := r.Group("/api")
	auth.Use(middleware.AuthMiddleware(cfg.JWTSecret))

	// 啟動伺服器
	if err := r.Run(":8080"); err != nil {
		log.Fatal("伺服器啟動失敗:", err)
	}
}
