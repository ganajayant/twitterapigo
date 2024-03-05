package main

import (
	"github.com/ganajayant/twitterapigo/controllers"
	"github.com/ganajayant/twitterapigo/initializers"
	"github.com/gin-gonic/gin"
)

func init() {
	initializers.LoadENV()
	initializers.Connect()
	initializers.GetJWTKey()
}

func main() {
	r := gin.Default()
	r.POST("/user", controllers.UserCreation)
	r.GET("/user", controllers.UserLogin)
	r.PUT("/user", controllers.UserGet)
	r.Run()
}
