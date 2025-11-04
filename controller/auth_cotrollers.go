package controller

import (
	"net/http"
	"time"

	models "github.com/21RMT14Muthuram/my-new-app/model"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)



var jwtKey = []byte("my-secret-key")

type Claims struct{
	Username string `json: "username"`
	jwt.RegisteredClaims
}

func SignUpHandler(c *gin.Context){
	var newUser models.User

	if err := c.ShouldBindJSON(&newUser); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Bad request"})
		return
	}
	if _, exists := models.UserStore[newUser.Username]; exists{
		c.JSON(http.StatusBadRequest, gin.H{"error": "User Already exists"})
		return
	}
	models.UserStore[newUser.Username] = newUser.Password
	c.JSON(http.StatusOK, gin.H{"message": "User registered successfully"})

}


func LoginHandler(c *gin.Context) {
	var credentials models.User
	if err := c.ShouldBindJSON(&credentials); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	storedPassword, ok := models.UserStore[credentials.Username]
	if !ok || storedPassword != credentials.Password {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid username or password"})
		return
	}
	expirationTime := time.Now().Add(1 * time.Hour)
	claims := &Claims{
		Username: credentials.Username,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(jwtKey)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not generate token"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"token": tokenString})
}