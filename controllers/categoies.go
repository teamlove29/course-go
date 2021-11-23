package controllers

import (
	"course-go/models"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/jinzhu/copier"
	"github.com/jinzhu/gorm"
)

type Categories struct {
	DB *gorm.DB
}

type categoryResponse struct {
	Name string `json:"name"`
	Desc string `json:"desc"`
}

type createCategoryForm struct {
	Name string `form:"name" binding:"required"`
	Desc string `form:"desc"binding:"required"`
}

type updateCategoryForm struct {
	Name string `form:"name"`
	Desc string `form:"desc"`
}

func (c *Categories) FindAll(ctx *gin.Context) {
	var categories []models.Category

	c.DB.Order("id desc").Find(&categories)
	var serializedCategories []categoryResponse
	copier.Copy(&serializedCategories, &categories)
	ctx.JSON(http.StatusOK, gin.H{"categories": serializedCategories})
}

func (c *Categories) Create(ctx *gin.Context) {

	var form createCategoryForm

	if err := ctx.ShouldBind(&form); err != nil {
		ctx.JSON(http.StatusUnprocessableEntity, gin.H{"error": err.Error()})
		return
	}

	var categories models.Category
	copier.Copy(&categories, &form)
	if err := c.DB.Create(&categories).Error; err != nil {
		ctx.JSON(http.StatusUnprocessableEntity, gin.H{"error": err.Error()})
		return
	}

	var serializedCategories categoryResponse
	copier.Copy(&serializedCategories, &categories)

	ctx.JSON(http.StatusOK, gin.H{"categories": serializedCategories})

}

func (c *Categories) FindOne(ctx *gin.Context) {
	categories, err := c.FineByID(ctx)
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	var serializedCategories categoryResponse
	copier.Copy(&serializedCategories, &categories)
	ctx.JSON(http.StatusOK, gin.H{"categoies": serializedCategories})
}

func (c *Categories) FineByID(ctx *gin.Context) (*models.Category, error) {
	var categoies models.Category

	id := ctx.Param("id")

	if err := c.DB.First(&categoies, id).Error; err != nil {
		return nil, err
	}

	return &categoies, nil

}

func (c *Categories) Update(ctx *gin.Context) {
	var form updateCategoryForm

	if err := ctx.ShouldBind(&form); err != nil {
		ctx.JSON(http.StatusUnprocessableEntity, gin.H{"error": err.Error()})
	}

	categories, err := c.FineByID(ctx)
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
	}

	if err := c.DB.Model(&categories).Update(&form).Error; err != nil {
		ctx.JSON(http.StatusUnprocessableEntity, gin.H{"error": err.Error()})
		return
	}

	var serializedCategories categoryResponse
	copier.Copy(&serializedCategories, &categories)
	ctx.JSON(http.StatusOK, gin.H{"categoies": serializedCategories})

}

func (c *Categories) Delete(ctx *gin.Context) {

	categories, err := c.FineByID(ctx)
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	c.DB.Unscoped().Delete(&categories)

	ctx.Status(http.StatusNoContent)

}
