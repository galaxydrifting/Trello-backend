package routes

func (r *Router) setupAuthRoutes() {
	// 公開路由
	r.engine.POST("/register", r.handlers.Auth.Register)
	r.engine.POST("/login", r.handlers.Auth.Login)
}
