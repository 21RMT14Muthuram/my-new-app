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


func GetUsers(c *gin.Context) {
    var userList []models.User

    rows, err := Config.DB.Query(`SELECT email, password_hash FROM users`)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch users"})
        fmt.Printf("Query error: %v\n", err)
        return
    }
    defer rows.Close()

    for rows.Next() {
        var user models.User
        err := rows.Scan(&user.Usermail, &user.Password)
        if err != nil {
            c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to scan user data"})
            fmt.Printf("Scan error: %v\n", err)
            return
        }
        userList = append(userList, user)
    }

    // Check for errors during iteration
    if err = rows.Err(); err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Error iterating through users"})
        fmt.Printf("Rows error: %v\n", err)
        return
    }

    c.JSON(http.StatusOK, gin.H{"users": userList})
}



func SignUpHandler(c *gin.Context){
	var newUser models.User
	if err := c.ShouldBindJSON(&newUser); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Bad request"})
		return
	}

	// Validate input
    if !IsValidEmail(newUser.Usermail) && !IsValidPassword(newUser.Password) {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Mismatch you Email or Password"})
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

	newUser.Password , _ = Hashing(newUser.Password)
	_, err = Config.DB.Exec(
		`insert into users (email, password_hash) values (?, ?)`,
		newUser.Usermail, newUser.Password ,
	)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to insert !"})
		fmt.Print("Insert error :", err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "User registered successfully"})

}

func LoginHandler(c *gin.Context){
	var lguser models.User
	if err := c.ShouldBindJSON(&lguser); err != nil{
		c.JSON(http.StatusBadRequest, gin.H{"message": "Ivaild Request"})
		return
	}

	 if lguser.Usermail == "" || lguser.Password == "" {
        c.JSON(http.StatusBadRequest, gin.H{"message": "Email and password are required"})
        return
    }

	var checkpass string
	err := Config.DB.QueryRow(`SELECT password_hash FROM users where email = ? `, lguser.Usermail).Scan(&checkpass)

	if err != nil {
        if err == sql.ErrNoRows {
            c.JSON(http.StatusUnauthorized, gin.H{"message": "Invalid credentials"})
        } else {
            c.JSON(http.StatusInternalServerError, gin.H{"message": "Database error"})
            fmt.Printf("Database error: %v\n", err)
        }
        return
    }

	if CheckHashPass(lguser.Password, checkpass){
		c.JSON(http.StatusOK, gin.H{"message": "Login Successfully"})
	} else {
		c.JSON(http.StatusUnauthorized, gin.H{"message": "It's very wrong bro!"})
	}

}

func DeleteUser(c *gin.Context){
	id := c.Param("id")
	result, err := Config.DB.Exec(`DELETE FROM users WHERE id = ?`, id)

	if err != nil{
		c.JSON(http.StatusInternalServerError, gin.H{"message":"DataBase Error"})
		return
	}
	row_Affected, _ := result.RowsAffected()
	if row_Affected == 0{
		c.JSON(http.StatusNotFound, gin.H{"message":"User not Found"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "User deleted successfully !"})
}
