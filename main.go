package main

import (
	"fmt"
	"time"

	"github.com/21RMT14Muthuram/my-new-app/controller"
	Config "github.com/21RMT14Muthuram/my-new-app/database"
	"github.com/21RMT14Muthuram/my-new-app/middleware"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func main() {
	// Initialize database connection
	if err := Config.InitDB(); err != nil { // Changed from Connect() to InitDB()
		fmt.Printf("Failed to connect to database: %v\n", err)
		return
	}
	defer Config.CloseDB() // Changed from DB.Close() to CloseDB()

	// Run migrations (you'll need to update your migration function too)

	

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
	r.POST("/verify-otp", controller.VerifyOTPHandler)
	r.POST("/resend-otp", controller.ResendOTPHandler)
	r.DELETE("/delete/:id", controller.DeleteUser)
	
	// Google OAuth
	r.GET("/login/google", controller.HandleGoogleLogin)
	r.GET("/auth/google/callback", controller.HandleGoogleCallback)


	auth := r.Group("/auth")
	auth.Use(middleware.AuthMiddleware())
		
	auth.GET("/content", controller.Content)

	fmt.Println("Server starting on :9000")
	
	// Check email configuration
	if !controller.IsEmailConfigured() {
		fmt.Println("   Email configuration not set. OTPs will be returned in API responses.")
		fmt.Println("   Set SMTP_USERNAME, SMTP_PASSWORD, and FROM_EMAIL environment variables")
	}
	
	r.Run(":9000")
}