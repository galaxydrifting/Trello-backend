package main

import (
	"log"

	"github.com/gin-gonic/gin"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	"trello-backend/docs"
	"trello-backend/internal/app"
	"trello-backend/internal/config"
	"trello-backend/internal/models"
	"trello-backend/internal/routes"
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

func setupSwagger() {
	docs.SwaggerInfo.Title = "Trello 後端 API"
	docs.SwaggerInfo.Description = "Trello 後端 API 文件"
	docs.SwaggerInfo.Version = "1.0"
	docs.SwaggerInfo.Host = "localhost:8080"
	docs.SwaggerInfo.BasePath = "/"
	docs.SwaggerInfo.Schemes = []string{"http"}
}

func main() {
	cfg := config.LoadConfig()
	db := initDB(cfg)

	// 初始化 Swagger 文件
	setupSwagger()

	// 使用 wire 進行相依性注入
	handlers, err := app.InitializeHandlers(db, cfg.JWTSecret)
	if err != nil {
		log.Fatal("無法初始化 Handlers:", err)
	}

	// 設定路由
	engine := gin.Default()
	router := routes.NewRouter(engine, handlers, cfg.JWTSecret)
	router.SetupRoutes()

	// 啟動伺服器
	if err := engine.Run(":8080"); err != nil {
		log.Fatal("伺服器啟動失敗:", err)
	}
}
