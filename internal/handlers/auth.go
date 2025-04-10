package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"trello-backend/internal/models"
	"trello-backend/internal/services"
)

type AuthHandler struct {
	db      *gorm.DB
	authSvc *services.AuthService
}

func NewAuthHandler(db *gorm.DB, authSvc *services.AuthService) *AuthHandler {
	return &AuthHandler{
		db:      db,
		authSvc: authSvc,
	}
}

func (h *AuthHandler) Register(c *gin.Context) {
	var req models.RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	token, err := h.authSvc.Register(req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, models.LoginResponse{Token: token})
}

func (h *AuthHandler) Login(c *gin.Context) {
	var req models.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	token, err := h.authSvc.Login(req)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "帳號或密碼錯誤"})
		return
	}

	c.JSON(http.StatusOK, models.LoginResponse{Token: token})
}