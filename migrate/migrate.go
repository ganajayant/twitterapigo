package main

import (
	"github.com/ganajayant/twitterapigo/initializers"
	"github.com/ganajayant/twitterapigo/models"
)

func init() {
	initializers.LoadENV()
	initializers.Connect()
}
func main() {
	initializers.Db.AutoMigrate(&models.User{})
}
