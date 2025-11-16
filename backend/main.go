package main

import (
	"fmt"
	"github.com/21RMT14Muthuram/my-new-app/controller"
	Config "github.com/21RMT14Muthuram/my-new-app/database"
	"github.com/gin-gonic/gin"
	"github.com/gin-contrib/cors"
	"time"
)

func main() {
	// Initialize database connection
	if err := Config.Connect(); err != nil {
		fmt.Printf("Failed to connect to database: %v\n", err)
		return
	}
	defer Config.DB.Close()

	// Run migrations
	if err := Config.Migrate(); err != nil {
		fmt.Printf("Migration failed: %v\n", err)
		return
	}

	// Initialize email configuration
	controller.InitEmailConfig()
	r := gin.Default()
r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:5173"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	// Public routes
	r.GET("/get-users", controller.GetUsers)
	r.POST("/signup", controller.SignUpHandler)
	r.POST("/login", controller.LoginHandler)
	r.POST("/verify-otp", controller.VerifyOTPHandler)     // New
	r.POST("/resend-otp", controller.ResendOTPHandler)     // New
	r.DELETE("/delete/:id", controller.DeleteUser)

	fmt.Println("Server starting on :9000")
	
	// Check email configuration
	if !controller.IsEmailConfigured() {
		fmt.Println("   Email configuration not set. OTPs will be returned in API responses.")
		fmt.Println("   Set SMTP_USERNAME, SMTP_PASSWORD, and FROM_EMAIL environment variables")
	}
	
	r.Run(":9000")
}