package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"trello-backend/internal/models"
	"trello-backend/internal/services"
)

// AuthHandler 處理認證相關的請求
type AuthHandler struct {
	authSvc services.AuthService
}

func NewAuthHandler(authSvc services.AuthService) *AuthHandler {
	return &AuthHandler{
		authSvc: authSvc,
	}
}

// Register godoc
// @Summary 使用者註冊
// @Description 註冊新使用者並返回 JWT 令牌
// @Tags 認證
// @Accept json
// @Produce json
// @Param request body models.RegisterRequest true "註冊資訊"
// @Success 201 {object} models.AuthResponse "註冊成功"
// @Failure 400 {object} models.APIResponse "無效的請求資料"
// @Failure 500 {object} models.APIResponse "內部伺服器錯誤"
// @Router /auth/register [post]
func (h *AuthHandler) Register(c *gin.Context) {
	var req models.RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, models.APIResponse{Error: err.Error()})
		return
	}

	resp, err := h.authSvc.Register(req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.APIResponse{Error: err.Error()})
		return
	}

	c.JSON(http.StatusCreated, resp)
}

// Login godoc
// @Summary 使用者登入
// @Description 使用電子郵件和密碼登入並返回 JWT 令牌
// @Tags 認證
// @Accept json
// @Produce json
// @Param request body models.LoginRequest true "登入資訊"
// @Success 200 {object} models.AuthResponse "登入成功"
// @Failure 400 {object} models.APIResponse "無效的請求資料"
// @Failure 401 {object} models.APIResponse "帳號或密碼錯誤"
// @Router /auth/login [post]
func (h *AuthHandler) Login(c *gin.Context) {
	var req models.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, models.APIResponse{Error: err.Error()})
		return
	}

	resp, err := h.authSvc.Login(req)
	if err != nil {
		c.JSON(http.StatusUnauthorized, models.APIResponse{Error: "帳號或密碼錯誤"})
		return
	}

	c.JSON(http.StatusOK, resp)
}

// ChangePassword godoc
// @Summary 變更使用者密碼
// @Description 使用者變更密碼功能
// @Tags 認證
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body models.ChangePasswordRequest true "變更密碼資訊"
// @Success 200 {object} models.APIResponse "密碼變更成功"
// @Failure 400 {object} models.APIResponse "無效的請求資料"
// @Failure 401 {object} models.APIResponse "認證失敗"
// @Router /auth/change-password [post]
func (h *AuthHandler) ChangePassword(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, models.APIResponse{Error: "未認證"})
		return
	}

	var req models.ChangePasswordRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, models.APIResponse{Error: err.Error()})
		return
	}

	if err := h.authSvc.ChangePassword(userID.(uuid.UUID), req); err != nil {
		c.JSON(http.StatusBadRequest, models.APIResponse{Error: err.Error()})
		return
	}

	c.JSON(http.StatusOK, models.APIResponse{})
}

// GetProfile godoc
// @Summary 取得目前使用者資訊
// @Description 取得目前登入使用者的 name、email
// @Tags 認證
// @Produce json
// @Security BearerAuth
// @Success 200 {object} models.UserProfileResponse "使用者資訊"
// @Failure 401 {object} models.APIResponse "認證失敗"
// @Router /auth/me [get]
func (h *AuthHandler) GetProfile(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, models.APIResponse{Error: "未認證"})
		return
	}
	resp, err := h.authSvc.GetProfile(userID.(uuid.UUID))
	if err != nil {
		c.JSON(http.StatusUnauthorized, models.APIResponse{Error: err.Error()})
		return
	}
	c.JSON(http.StatusOK, resp)
}

// Ping godoc
// @Summary 測試 API 是否正常運作
// @Description 回傳簡單的 pong 回應
// @Tags 系統
// @Produce json
// @Success 200 {object} models.APIResponse "成功回應"
// @Router /auth/ping [get]
func (h *AuthHandler) Ping(c *gin.Context) {
	c.JSON(http.StatusOK, models.APIResponse{})
}
