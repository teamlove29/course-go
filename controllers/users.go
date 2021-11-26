package controllers

import (
	"course-go/config"
	"course-go/models"
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/jinzhu/copier"
	"github.com/jinzhu/gorm"
)

type Users struct {
	DB *gorm.DB
}

type createUserForm struct {
	Email    string `form:"email" binding:"required,email"`
	Password string `form:"password" binding:"required,min=8"`
	Name     string `form:"name" binding:"required"`
}

type updateUserForm struct {
	Email string `form:"email" binding:"omitempty,email"`
	// omitempty = ไม่มีค่า = no validate
	Password string `form:"password" binding:"omitempty,min=8"`
	Name     string `form:"name"`
}

type UserResponse struct {
	ID     uint   `json:"id"`
	Email  string `json:"email"`
	Avatar string `json:"avatar"`
	Name   string `json:"name"`
	Role   string `json:"role"`
}

type usersPaging struct {
	Items  []userResponse `json:"items"`
	Paging *pagingResult  `json:"paging"`
}

// api/v1/users?term=babel
func (u *Users) FindAll(ctx *gin.Context) {
	var users []models.User
	query := u.DB.Order("id desc").Find(&users)

	term := ctx.Query("term")
	if term != "" {
		query = query.Where("name ILIKE ?", "%"+term+"%")
	}
	pagination := pagination{c: ctx, query: query, recodes: &users}
	paging := pagination.paginate()

	serializedUsers := []userResponse{}
	copier.Copy(&serializedUsers, &users)
	ctx.JSON(http.StatusOK, gin.H{
		"users": usersPaging{Items: serializedUsers, Paging: paging},
	})
}

func (u *Users) FindOne(ctx *gin.Context) {
	user, err := u.findUserByID(ctx)
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	var serializedUser userResponse
	copier.Copy(&serializedUser, &user)
	ctx.JSON(http.StatusOK, gin.H{"user": serializedUser})
}

func (u *Users) Create(ctx *gin.Context) {
	var form createUserForm
	if err := ctx.ShouldBindJSON(&form); err != nil {
		ctx.JSON(http.StatusUnprocessableEntity, gin.H{"error": err.Error()})
		return
	}

	var user models.User
	copier.Copy(&user, &form)
	user.Password = user.GenerateEncryptedPassword()

	if err := u.DB.Create(&user).Error; err != nil {
		ctx.JSON(http.StatusUnprocessableEntity, gin.H{"error": err.Error()})
		return
	}

	var serializedUser userResponse
	copier.Copy(&serializedUser, &user)
	ctx.JSON(http.StatusCreated, gin.H{"user": serializedUser})
}

func (u *Users) Update(ctx *gin.Context) {
	var form updateUserForm
	if err := ctx.ShouldBindJSON(&form); err != nil {
		ctx.JSON(http.StatusUnprocessableEntity, gin.H{"error": err.Error()})
		return
	}

	user, err := u.findUserByID(ctx)
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	if form.Password != "" {
		user.Password = user.GenerateEncryptedPassword()
	}

	if err := u.DB.Model(&user).Update(&form).Error; err != nil {
		ctx.JSON(http.StatusUnprocessableEntity, gin.H{"error": err.Error()})
		return
	}

	var serializedUser userResponse
	copier.Copy(&serializedUser, &user)
	ctx.JSON(http.StatusOK, gin.H{"user": serializedUser})
}

func (u *Users) Delete(ctx *gin.Context) {
	user, err := u.findUserByID(ctx)
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	u.DB.Unscoped().Delete(&user)

	ctx.Status(http.StatusNoContent)
}

func (u *Users) Promote(ctx *gin.Context) {
	user, err := u.findUserByID(ctx)
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	user.Promote()
	u.DB.Save(user)

	var serializedUser userResponse
	copier.Copy(&serializedUser, &user)
	ctx.JSON(http.StatusOK, gin.H{"user": serializedUser})
}

func (u *Users) Demote(ctx *gin.Context) {
	user, err := u.findUserByID(ctx)
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	user.Demote()
	u.DB.Save(user)

	var serializedUser userResponse
	copier.Copy(&serializedUser, &user)
	ctx.JSON(http.StatusOK, gin.H{"user": serializedUser})
}

func (u *Users) findUserByID(ctx *gin.Context) (*models.User, error) {
	id := ctx.Param("id")
	var user models.User

	if err := u.DB.First(&user, id).Error; err != nil {
		return nil, err
	}

	return &user, nil
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
