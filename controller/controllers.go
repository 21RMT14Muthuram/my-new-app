package controller

import (
	"fmt"
	"net/http"

	"github.com/21RMT14Muthuram/my-new-app/model"
	"github.com/gin-gonic/gin"
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

func Content(c *gin.Context){
    c.JSON(http.StatusOK, gin.H{"message": "you find it !"})
}