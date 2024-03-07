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
	initializers.ConnectFireBase()
}

func main() {
	r := gin.Default()
	r.MaxMultipartMemory = 8 << 20
	userRoutes := r.Group("/user")
	{
		userRoutes.POST("/create", controllers.UserCreation)
		userRoutes.GET("/login", controllers.UserLogin)
		userRoutes.GET("/info", controllers.UserGet)
	}

	tweetRoutes := r.Group("/tweet")
	{
		tweetRoutes.POST("/create", controllers.TweetCreation)
		tweetRoutes.PUT("/edit", controllers.TweetEdit)
	}

	r.Run()
}
