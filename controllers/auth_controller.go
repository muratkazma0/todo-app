package controllers

import (
	"net/http"
	"os"
	"time"

	"todoapp/entity"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/joho/godotenv"
)

var SecretKey = os.Getenv("JWT_SECRET")

func init() {
	err := godotenv.Load()
	if err != nil {
		panic("Error loading .env file")
	}

	SecretKey = os.Getenv("JWT_SECRET")
	if SecretKey == "" {
		panic("JWT_SECRET not set in .env file")
	}
}

type AuthController struct {
	userModel *entity.UserModel
}

func NewAuthController(userModel *entity.UserModel) *AuthController {
	return &AuthController{
		userModel: userModel,
	}
}

type LoginRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

func (c *AuthController) Login(ctx *gin.Context) {
	var req LoginRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user, err := c.userModel.GetByUsername(req.Username)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "invalid credentials"})
		return
	}

	if !user.CheckPassword(req.Password) {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "invalid credentials"})
		return
	}

	// Create token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": user.ID,
		"role":    user.Role,
		"exp":     time.Now().Add(time.Hour * 24).Unix(),
	})

	// Sign token
	tokenString, err := token.SignedString([]byte(SecretKey))
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "failed to generate token"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"token": tokenString,
		"user": gin.H{
			"id":       user.ID,
			"username": user.Username,
			"role":     user.Role,
		},
	})
}
