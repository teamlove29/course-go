package controllers

import (
	"course-go/config"
	"course-go/models"
	"mime/multipart"
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/jinzhu/copier"
	"github.com/jinzhu/gorm"
)

type Auth struct {
	DB *gorm.DB
}

type authForm struct {
	Email    string `form:"email" binding:"required,email"`
	Password string `form:"password" binding:"required,min=8"`
}

type authResponse struct {
	ID    uint   `json:"id"`
	Email string `json:"email"`
}

type userResponse struct {
	ID     uint   `json:"id"`
	Email  string `json:"email"`
	Avatar string `json:"avatar"`
	Name   string `json:"name"`
	Role   string `json:"role"`
}

type updateProfileForm struct {
	Email  string                `form:"email"`
	Name   string                `form:"name"`
	Avatar *multipart.FileHeader `form:"avatar"`
}

func (a *Auth) Signup(ctx *gin.Context) {

	var form authForm
	if err := ctx.ShouldBind(&form); err != nil {
		ctx.JSON(http.StatusUnprocessableEntity, gin.H{"error": err.Error()})
		return
	}

	var user models.User
	copier.Copy(&user, &form)
	user.Password = user.GenerateEncryptedPassword()
	if err := a.DB.Create(&user).Error; err != nil {
		ctx.JSON(http.StatusUnprocessableEntity, gin.H{"error": err.Error()})
		return
	}

	var serializedUser authResponse
	copier.Copy(&serializedUser, &user)
	ctx.JSON(http.StatusCreated, gin.H{"user": serializedUser})
}

// auth/profile => jwt => sub(UserID) => User => responUser
func (a *Auth) GetProfile(ctx *gin.Context) {

	sub, _ := ctx.Get("sub")
	user := sub.(*models.User)

	var seriailzedUser userResponse
	copier.Copy(&seriailzedUser, &user)
	ctx.JSON(http.StatusOK, gin.H{"User": seriailzedUser})

}

func (a *Auth) UpdateProfile(ctx *gin.Context) {
	var form updateProfileForm

	if err := ctx.ShouldBind(&form); err != nil {
		ctx.JSON(http.StatusUnprocessableEntity, gin.H{"Error": err.Error()})
		return
	}

	sub, _ := ctx.Get("sub")
	user := sub.(*models.User)

	setUserImage(ctx, user)
	if err := a.DB.Model(&user).Update(&form).Error; err != nil {
		ctx.JSON(http.StatusUnprocessableEntity, gin.H{"Error": err.Error()})
		return
	}

	var seriailzedUser userResponse
	copier.Copy(&seriailzedUser, &user)
	ctx.JSON(http.StatusOK, gin.H{"User": seriailzedUser})

}

func setUserImage(ctx *gin.Context, user *models.User) error {
	file, err := ctx.FormFile("avatar")
	if err != nil || file == nil {
		return err
	}

	// เช็คว่ามีรูปอยู่ในไฟล์อยู่รึเปล่า
	if user.Avatar != "" {
		// http://127.0.0.1:5000/uploads/users/<ID>/image.png
		// 1. /uploads/articles/<ID>/image.png
		user.Avatar = strings.Replace(user.Avatar, os.Getenv("HOST"), "", 1)
		// 2. <WD>/uploads/users/<ID>/image.png <WD> = working dir
		pwd, _ := os.Getwd()
		// 3. Remove <WD>/uploads/users/<ID>/image.png
		os.Remove(pwd + user.Avatar)
	}

	// create Path
	// uploads/users/123
	path := "uploads/users/" + strconv.Itoa(int(user.ID))
	os.MkdirAll(path, 0755)

	// Upload file
	// uploads/users/123/filename
	filename := path + "/" + file.Filename
	if err := ctx.SaveUploadedFile(file, filename); err != nil {
		return err
	}
	// Attach file to article
	user.Avatar = os.Getenv("HOST") + "/" + filename
	// update to sql
	db := config.GetDB()
	db.Save(user)

	return nil

}
