package handlers

// Handlers 用於管理所有的 HTTP handlers
type Handlers struct {
	Auth *AuthHandler
	// 在這裡新增其他的 handlers
	// Board       *BoardHandler
	// Card        *CardHandler
	// List        *ListHandler
	// Comment     *CommentHandler
	// User        *UserHandler
	// Workspace   *WorkspaceHandler
	// Activity    *ActivityHandler
	// Label       *LabelHandler
	// Attachment  *AttachmentHandler
}

// NewHandlers 建立新的 Handlers 實例
func NewHandlers(
	authHandler *AuthHandler,
	// 在這裡新增其他的 handler 參數
) *Handlers {
	return &Handlers{
		Auth: authHandler,
		// 在這裡初始化其他的 handlers
	}
}
