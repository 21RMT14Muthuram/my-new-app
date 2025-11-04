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
        Username: "Muthuram",
        Password: "12345",
    }

    users = append(users, user)

    fmt.Println("Added user:", user)
}
