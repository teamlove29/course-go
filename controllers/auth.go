package controllers

import (
	"course-go/models"
	"mime/multipart"
	"net/http"

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
