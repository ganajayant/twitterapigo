package controllers

import (
	"fmt"
	"io"
	"net/http"
	"os"

	"github.com/ganajayant/twitterapigo/initializers"
	"github.com/ganajayant/twitterapigo/models"
	"github.com/gin-gonic/gin"
)

func TweetCreation(ctx *gin.Context) {
	token := ctx.GetHeader("Authorization")
	if token == "" {
		ctx.JSON(http.StatusUnauthorized, gin.H{
			"message": "Token is not provided",
		})
		return
	}
	tokenString := token[len("Bearer "):]
	payload, err := initializers.VerifyToken(tokenString)
	if err != nil {
		fmt.Println(err)
		ctx.JSON(http.StatusUnauthorized, gin.H{
			"message": "Token is invalid",
		})
		return
	}
	user := models.User{}
	email := payload["email"].(string)
	initializers.Db.Where("email = ?", email).First(&user)
	if user.ID == "" {
		ctx.JSON(http.StatusNotFound, gin.H{
			"message": "User not found",
		})
		return
	}
	file, err := ctx.FormFile("tweetimage")
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": "can't get the file",
		})
		return
	}
	text := ctx.Request.FormValue("text")
	if text == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": "text is not provided",
		})
		return
	}
	imagepath := file.Filename
	Bucket := initializers.Bucket
	wc := Bucket.Object(imagepath).NewWriter(ctx)
	f, err := file.Open()
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": "can't open the file",
		})
		return
	}
	_, err = io.Copy(wc, f)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": "can't copy the file",
		})
		return
	}
	if err := wc.Close(); err != nil {
		fmt.Println(err)
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": "can't close the file",
		})
		return
	}
	imageurl := "https://storage.cloud.google.com/" + os.Getenv("Bucket") + "/" + imagepath

	tweet := models.Tweet{
		Text:     text,
		UserID:   user.ID,
		ImageUrl: imageurl,
	}
	initializers.Db.Create(&tweet)
	ctx.JSON(http.StatusOK, gin.H{
		"message": "Tweet is created",
	})
}

func TweetEdit(ctx *gin.Context) {
	token := ctx.GetHeader("Authorization")
	if token == "" {
		ctx.JSON(http.StatusUnauthorized, gin.H{
			"message": "Token is not provided",
		})
		return
	}
	tokenString := token[len("Bearer "):]
	payload, err := initializers.VerifyToken(tokenString)
	if err != nil {
		fmt.Println(err)
		ctx.JSON(http.StatusUnauthorized, gin.H{
			"message": "Token is invalid",
		})
		return
	}
	user := models.User{}
	email := payload["email"].(string)
	initializers.Db.Where("email = ?", email).First(&user)
	if user.ID == "" {
		ctx.JSON(http.StatusNotFound, gin.H{
			"message": "User not found",
		})
		return
	}
	var Body struct {
		Text string
	}
	err1 := ctx.BindJSON(&Body)
	if err1 != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": "can't bind the body",
		})
		return
	}
	id := ctx.Query("id")
	if id == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": "id is not provided",
		})
		return
	}
	tweet := models.Tweet{}
	initializers.Db.Where("id = ? AND user_id = ?", id, user.ID).First(&tweet)
	if tweet.ID == "" {
		ctx.JSON(http.StatusNotFound, gin.H{
			"message": "Tweet not found",
		})
		return
	}
	initializers.Db.Model(&tweet).Update("text", Body.Text)
	ctx.JSON(http.StatusOK, gin.H{
		"message": "Tweet is updated",
	})
}
