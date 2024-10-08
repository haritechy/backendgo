package middleware

import (
	"net/http"
	"os"
	"slack-chatbot/database"
	"slack-chatbot/models"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
)

var jwtKey = []byte(os.Getenv("JWT_KEY"))

func AuthorizeJWT() gin.HandlerFunc {
	return func(c *gin.Context) {
		const BearerSchema = "Bearer "
		authHeader := c.GetHeader("Authorization")

		// Check if the authorization header is missing
		if authHeader == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Authorization header is missing"})
			return
		}

		// Ensure the token format is valid (Bearer Token)
		if !strings.HasPrefix(authHeader, BearerSchema) {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Invalid token format"})
			return
		}

		tokenString := strings.TrimSpace(authHeader[len(BearerSchema):])

		claims := &models.Claims{}

		// Parse the JWT token
		token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
			return jwtKey, nil
		})

		// Check for parsing errors
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized: " + err.Error()})
			return
		}

		// Validate the token
		if !token.Valid {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
			return
		}

		// Use the email from claims to find the user
		var user models.User
		if err := database.DB.Where("email = ?", claims.Email).First(&user).Error; err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "User not found"})
			return
		}

		// Set the user ID in the context
		c.Set("id", user.ID)
		c.Set("claims", claims)

		// Continue with the request
		c.Next()
	}
}
