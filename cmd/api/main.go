package main

import (
	"log"

	"github.com/gin-gonic/gin"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	"trello-backend/docs"
	"trello-backend/graph"
	"trello-backend/internal/app"
	"trello-backend/internal/config"
	"trello-backend/internal/routes"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/handler/extension"
	"github.com/99designs/gqlgen/graphql/handler/lru"
	"github.com/99designs/gqlgen/graphql/handler/transport"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/vektah/gqlparser/v2/ast"
)

// @title Trello 後端 API
// @version 1.0
// @description Trello 後端 API 文件
// @host localhost:8080
// @BasePath /api
// @schemes http
// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
func initDB(cfg *config.Config) *gorm.DB {
	db, err := gorm.Open(postgres.Open(cfg.GetDBConnString()), &gorm.Config{})
	if err != nil {
		log.Fatal("無法連線到資料庫:", err)
	}

	// Run migrations
	app.Migrate(db)

	log.Println("資料庫遷移完成")

	return db
}

func setupSwagger() {
	docs.SwaggerInfo.Title = "Trello 後端 API"
	docs.SwaggerInfo.Description = "Trello 後端 API 文件"
	docs.SwaggerInfo.Version = "1.0"
	docs.SwaggerInfo.Host = "localhost:8080"
	docs.SwaggerInfo.BasePath = "/api"
	docs.SwaggerInfo.Schemes = []string{"http"}
}

func main() {
	cfg := config.LoadConfig()
	db := initDB(cfg)

	// 初始化 Swagger 文件
	setupSwagger()

	// 使用 wire 進行相依性注入
	api, err := app.InitializeAPI(db, cfg.JWTSecret)
	if err != nil {
		log.Fatal("無法初始化 API:", err)
	}

	// 設定路由
	engine := gin.Default()

	// GraphQL 設定
	gqlSrv := handler.New(graph.NewExecutableSchema(graph.Config{
		Resolvers: graph.NewResolverFromAPI(api),
	}))
	gqlSrv.AddTransport(transport.Options{})
	gqlSrv.AddTransport(transport.GET{})
	gqlSrv.AddTransport(transport.POST{})
	gqlSrv.SetQueryCache(lru.New[*ast.QueryDocument](1000))
	gqlSrv.Use(extension.Introspection{})
	gqlSrv.Use(extension.AutomaticPersistedQuery{Cache: lru.New[string](100)})

	// GraphQL Playground 路由
	engine.GET("/graphql/playground", gin.WrapH(playground.Handler("GraphQL playground", "/graphql/query")))
	// GraphQL 查詢路由
	engine.POST("/graphql/query", func(c *gin.Context) {
		gqlSrv.ServeHTTP(c.Writer, c.Request)
	})

	router := routes.NewRouter(engine, cfg.JWTSecret, cfg)

	// 註冊所有 handlers
	for name, handler := range api.GetHandlers() {
		router.RegisterHandler(name, handler)
	}

	// 設定所有路由
	router.SetupRoutes()

	// 啟動伺服器
	if err := engine.Run(":8080"); err != nil {
		log.Fatal("伺服器啟動失敗:", err)
	}
}
