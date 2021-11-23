package routes

import (
	"course-go/config"
	"course-go/controllers"

	"github.com/gin-gonic/gin"
)

func Serve(r *gin.Engine) {

	db := config.GetDB()
	v1 := r.Group("/api/v1")

	articlesGroup := v1.Group("/articles")
	articleController := controllers.Articles{DB: db}
	{
		articlesGroup.GET("", articleController.FindAll)
		articlesGroup.GET("/:id", articleController.FindOne)
		articlesGroup.POST("", articleController.Create)
		articlesGroup.PATCH("/:id", articleController.Update)
		articlesGroup.DELETE("/:id", articleController.Delete)
	}

	categoiesGroup := v1.Group("/categories")
	categoiesController := controllers.Categories{DB: db}
	{
		categoiesGroup.GET("", categoiesController.FindAll)
		categoiesGroup.GET("/:id", categoiesController.FindOne)
		categoiesGroup.POST("", categoiesController.Create)
		categoiesGroup.PATCH("/:id", categoiesController.Update)
		categoiesGroup.DELETE("/:id", categoiesController.Delete)
	}

	authGroup := v1.Group("/auth")
	authController := controllers.Auth{DB: db}
	{
		authGroup.POST("/sign-up", authController.Signup)
	}

}
