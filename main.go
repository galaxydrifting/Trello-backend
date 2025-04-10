package main

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type User struct {
	ID           uuid.UUID `gorm:"type:uuid;primary_key"`
	Email        string    `gorm:"unique;not null"`
	Name         string    `gorm:"not null"`
	PasswordHash string    `gorm:"not null"`
	CreatedAt    time.Time `gorm:"default:CURRENT_TIMESTAMP"`
}

type RegisterRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Name     string `json:"name" binding:"required"`
	Password string `json:"password" binding:"required,min=6"`
}

type LoginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

type LoginResponse struct {
	Token string `json:"token"`
}

var db *gorm.DB

func initDB() {
	dsn := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%s sslmode=disable",
		os.Getenv("POSTGRES_HOSTNAME"),
		os.Getenv("POSTGRES_USER"),
		os.Getenv("POSTGRES_PASSWORD"),
		os.Getenv("POSTGRES_DB"),
		os.Getenv("POSTGRES_PORT"),
	)

	var err error
	db, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("無法連線到資料庫:", err)
	}

	err = db.AutoMigrate(&User{})
	if err != nil {
		log.Fatal("資料庫遷移失敗:", err)
	}
	fmt.Println("資料庫遷移完成")
}

func generateToken(userID uuid.UUID) (string, error) {
	token := jwt.New(jwt.SigningMethodHS256)
	claims := token.Claims.(jwt.MapClaims)
	claims["user_id"] = userID
	claims["exp"] = time.Now().Add(24 * time.Hour).Unix()

	return token.SignedString([]byte(os.Getenv("JWT_SECRET")))
}

func main() {
	initDB()
	r := gin.Default()

	r.POST("/register", func(c *gin.Context) {
		var req RegisterRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(400, gin.H{"error": err.Error()})
			return
		}

		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
		if err != nil {
			c.JSON(500, gin.H{"error": "密碼加密失敗"})
			return
		}

		user := User{
			ID:           uuid.New(),
			Email:        req.Email,
			Name:         req.Name,
			PasswordHash: string(hashedPassword),
		}

		if err := db.Create(&user).Error; err != nil {
			c.JSON(500, gin.H{"error": "使用者建立失敗"})
			return
		}

		token, err := generateToken(user.ID)
		if err != nil {
			c.JSON(500, gin.H{"error": "Token 產生失敗"})
			return
		}

		c.JSON(201, LoginResponse{Token: token})
	})

	r.POST("/login", func(c *gin.Context) {
		var req LoginRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(400, gin.H{"error": err.Error()})
			return
		}

		var user User
		if err := db.Where("email = ?", req.Email).First(&user).Error; err != nil {
			c.JSON(401, gin.H{"error": "帳號或密碼錯誤"})
			return
		}

		if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(req.Password)); err != nil {
			c.JSON(401, gin.H{"error": "帳號或密碼錯誤"})
			return
		}

		token, err := generateToken(user.ID)
		if err != nil {
			c.JSON(500, gin.H{"error": "Token 產生失敗"})
			return
		}

		c.JSON(200, LoginResponse{Token: token})
	})

	r.Run(":8080")
}
