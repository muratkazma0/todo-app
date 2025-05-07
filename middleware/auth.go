package middleware

import (
	"net/http"
	"os"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/joho/godotenv"
)

var SecretKey string

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

func AuthMiddleware() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		authHeader := ctx.GetHeader("Authorization")
		if authHeader == "" {
			ctx.JSON(http.StatusUnauthorized, gin.H{"error": "authorization header is required"})
			ctx.Abort()
			return
		}

		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			ctx.JSON(http.StatusUnauthorized, gin.H{"error": "invalid authorization header format"})
			ctx.Abort()
			return
		}

		tokenString := parts[1]
		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, jwt.ErrSignatureInvalid
			}
			return []byte(SecretKey), nil
		})

		if err != nil {
			ctx.JSON(http.StatusUnauthorized, gin.H{"error": "invalid token"})
			ctx.Abort()
			return
		}

		if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
			userID := int(claims["user_id"].(float64))
			userRole := claims["role"].(string)

			ctx.Set("user_id", userID)
			ctx.Set("user_role", userRole)
			ctx.Next()
		} else {
			ctx.JSON(http.StatusUnauthorized, gin.H{"error": "invalid token claims"})
			ctx.Abort()
			return
		}
	}
}

func AdminOnly() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		userRole, exists := ctx.Get("user_role")
		if !exists || userRole != "admin" {
			ctx.JSON(http.StatusForbidden, gin.H{"error": "admin access required"})
			ctx.Abort()
			return
		}
		ctx.Next()
	}
}
