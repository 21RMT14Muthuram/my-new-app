package middleware

import (
	"net/http"
	"strings"

	// "github.com/21RMT14Muthuram/my-new-app/controller"
	"github.com/gin-gonic/gin"
	// "github.com/golang-jwt/jwt/v5"
)



func AuthMiddleware() gin.HandlerFunc {
	return func (c *gin.Context)  {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Missing Toker"})
			return
		}

		// support Bearer & raw tokens
		tokenString := strings.TrimPrefix(authHeader, "Bearer")
		tokenString = strings.TrimSpace(tokenString)

		// claims := &controller.Claims{}
		// token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		// 	return []byte("my_secret_key"), nil
		// })
		// if err != nil || !token.Valid {
		// 	c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Invalid or expired token"})
		// 	return
		// }

		// c.Set("username", claims.Username)
		// c.Next()
	}
}