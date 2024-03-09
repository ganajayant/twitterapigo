package controllers

import (
	"fmt"
	"mime/multipart"
	"net/http"

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
	var Body struct {
		Text       string
		TweetImage *multipart.FileHeader
	}
	err1 := ctx.Bind(&Body)
	if err1 != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": "can't bind the body",
		})
		return
	}

	imageurl, err := initializers.UploadFile(Body.TweetImage)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"message": "can't upload the file",
		})
		return
	}

	tweet := models.Tweet{
		Text:     Body.Text,
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
