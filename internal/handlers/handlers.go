package handlers

// HandlersOptions 用於設定 Handlers 的選項
type HandlersOptions struct {
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
func NewHandlers(opts HandlersOptions) *Handlers {
	return &Handlers{
		Auth: opts.Auth,
		// 在這裡初始化其他的 handlers
	}
}
