package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"

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
// @Success 201 {object} models.LoginResponse "註冊成功"
// @Failure 400 {object} models.APIResponse "無效的請求資料"
// @Failure 500 {object} models.APIResponse "內部伺服器錯誤"
// @Router /register [post]
func (h *AuthHandler) Register(c *gin.Context) {
	var req models.RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, models.APIResponse{Error: err.Error()})
		return
	}

	token, err := h.authSvc.Register(req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.APIResponse{Error: err.Error()})
		return
	}

	c.JSON(http.StatusCreated, models.LoginResponse{Token: token})
}

// Login godoc
// @Summary 使用者登入
// @Description 使用電子郵件和密碼登入並返回 JWT 令牌
// @Tags 認證
// @Accept json
// @Produce json
// @Param request body models.LoginRequest true "登入資訊"
// @Success 200 {object} models.LoginResponse "登入成功"
// @Failure 400 {object} models.APIResponse "無效的請求資料"
// @Failure 401 {object} models.APIResponse "帳號或密碼錯誤"
// @Router /login [post]
func (h *AuthHandler) Login(c *gin.Context) {
	var req models.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, models.APIResponse{Error: err.Error()})
		return
	}

	token, err := h.authSvc.Login(req)
	if err != nil {
		c.JSON(http.StatusUnauthorized, models.APIResponse{Error: "帳號或密碼錯誤"})
		return
	}

	c.JSON(http.StatusOK, models.LoginResponse{Token: token})
}
