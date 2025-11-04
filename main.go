package main

import (

	"github.com/21RMT14Muthuram/my-new-app/controller"
	"github.com/gin-gonic/gin"
)

func main(){
	// fmt.Print("Hello, World!")
	// controller.NewPrint()
	// controller.AddSingleUser()
	r := gin.Default()

	//public routes
	r.POST("./signup", controller.SignUpHandler)
	r.POST("/login", controller.LoginHandler)


	r.Run(":8000")
}