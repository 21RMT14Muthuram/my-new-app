package controller

import (
	"database/sql"
	"fmt"
	"net/http"
	"github.com/21RMT14Muthuram/my-new-app/database"
	models "github.com/21RMT14Muthuram/my-new-app/model"
	"github.com/gin-gonic/gin"
	// "github.com/golang-jwt/jwt/v5"
)



var jwtKey = []byte("my-secret-key")



func SignUpHandler(c *gin.Context){
	var newUser models.User

	if err := c.ShouldBindJSON(&newUser); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Bad request"})
		return
	}

	// Validate input
    if newUser.Usermail == "" || newUser.Password == "" {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Email and password are required"})
        return
    }

	//check email already exists
	var existingEmail string
	err := Config.DB.QueryRow(`SELECT * FROM users WHERE email = ?`, newUser.Usermail).Scan(&existingEmail)
	if  err == nil {
		c.JSON(http.StatusConflict, gin.H{"message": "User email already exist !"})
		fmt.Print("Email :", existingEmail)
		return
	} else if err != sql.ErrNoRows {
        // Some other error occurred
        c.JSON(http.StatusInternalServerError, gin.H{"message": "Database error"})
        fmt.Printf("Database error: %v\n", err)
        return
    }

	fmt.Print("Email - :", existingEmail)

	_, err = Config.DB.Exec(
		`insert into users (email, password_hash) values (?, ?)`,
		newUser.Usermail, newUser.Password,
	)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to insert !"})
		fmt.Print("Insert error :", err)
		return
	} 



	c.JSON(http.StatusOK, gin.H{"message": "User registered successfully"})

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