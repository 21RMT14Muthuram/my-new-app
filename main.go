package main

import (
    "fmt"
    "github.com/21RMT14Muthuram/my-new-app/controller"
    "github.com/21RMT14Muthuram/my-new-app/database"
    "github.com/gin-gonic/gin"
)

func main() {
    // Initialize database connection
    if err := Config.Connect(); err != nil {
        fmt.Printf("Failed to connect to database: %v\n", err)
        return
    }
    defer Config.DB.Close()
    

    r := gin.Default()

    // Public routes
	r.GET("/get-users", controller.GetUsers)
    r.POST("/signup", controller.SignUpHandler)
    r.POST("/login", controller.LoginHandler)
	r.DELETE("/delete/:id", controller.DeleteUser)
    fmt.Println("Server starting on :9000")
    r.Run(":9000")
}