package controllers

import (
	"fmt"
	"mime/multipart"
	"net/http"
	"time"

	"github.com/ganajayant/twitterapigo/initializers"
	"github.com/ganajayant/twitterapigo/models"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

func UserCreation(ctx *gin.Context) {
	var Body struct {
		Name     string
		Bio      string
		Email    string
		DoB      time.Time
		Password string
		Profile  *multipart.FileHeader
	}
	err := ctx.Bind(&Body)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": "can't bind the body",
		})
		return
	}
	alreadyData := models.User{}
	initializers.Db.Where("email = ?", Body.Email).First(&alreadyData)
	if alreadyData.ID != "" {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": "User already exists",
		})
		return
	}
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(Body.Password), bcrypt.DefaultCost)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": "can't hash the password",
		})
		return
	}
	imageurl, err := initializers.UploadFile(Body.Profile)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"message": "can't upload the file",
		})
		return
	}
	user := models.User{Name: Body.Name, Bio: Body.Bio, Email: Body.Email, DoB: Body.DoB, Password: string(hashedPassword), ProfilePicUrl: imageurl}
	initializers.Db.Create(&user)
	ctx.JSON(http.StatusAccepted, gin.H{
		"message": fmt.Sprintf("User is created with id:%v, name:%v, email:%v", user.ID, user.Name, user.Email),
	})
}

func UserLogin(ctx *gin.Context) {
	var Body struct {
		Email    string
		Password string
	}
	if err := ctx.BindJSON(&Body); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": "can't bind the body",
		})
		return
	}
	user := models.User{}
	initializers.Db.Where("email = ?", Body.Email).First(&user)
	if user.ID == "" {
		ctx.JSON(http.StatusNotFound, gin.H{
			"message": "User not found",
		})
		return
	}
	err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(Body.Password))
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{
			"message": "Password is incorrect",
		})
		return
	}
	token, err := initializers.CreateToken(user.Email, user.ID)
	if err != nil {
		fmt.Println(err)
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"message": "can't create token",
		})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{
		"message": "User is logged in",
		"token":   token,
	})
}

func UserGet(ctx *gin.Context) {
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
	ctx.JSON(http.StatusOK, gin.H{
		"message": "User is found",
		"user":    user,
	})
}
