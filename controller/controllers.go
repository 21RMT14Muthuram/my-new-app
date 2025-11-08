package controller

import (
	"fmt"

	"github.com/21RMT14Muthuram/my-new-app/model"
)

var users []models.User



func NewPrint(){
	fmt.Println("this new one !!!")
}

func AddSingleUser() {
    user := models.User{
        Usermail: "Muthuram",
        Password: "12345",
    }

    users = append(users, user)

    fmt.Println("Added user:", user)
}


// func LoginHandler(c *gin.Context) {
// 	var credentials models.User
// 	if err := c.ShouldBindJSON(&credentials); err != nil {
// 		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
// 		return
// 	}

// 	storedPassword, ok := models.UserStore[credentials.Usermail]
// 	if !ok || storedPassword != credentials.Password {
// 		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid username or password"})
// 		return
// 	}
// 	expirationTime := time.Now().Add(1 * time.Hour)
// 	claims := &Claims{
// 		Usermail: credentials.Usermail,
// 		RegisteredClaims: jwt.RegisteredClaims{
// 			ExpiresAt: jwt.NewNumericDate(expirationTime),
// 		},
// 	}
// 	fmt.Print("=====", models.UserStore)
// 	fmt.Print("=====", claims)

// 	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
// 	tokenString, err := token.SignedString(jwtKey)
// 	if err != nil {
// 		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not generate token"})
// 		return
// 	}

// 	c.JSON(http.StatusOK, gin.H{"token": tokenString})
// }