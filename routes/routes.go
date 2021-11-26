package routes

import (
	"course-go/config"
	"course-go/controllers"
	"course-go/middleware"

	"github.com/gin-gonic/gin"
)

func Serve(r *gin.Engine) {

	db := config.GetDB()
	v1 := r.Group("/api/v1")
	authenticate := middleware.Authenticate().MiddlewareFunc()
	authorize := middleware.Authorize()

	articlesGroup := v1.Group("/articles")
	articleController := controllers.Articles{DB: db}
	articlesGroup.GET("", articleController.FindAll)
	articlesGroup.GET("/:id", articleController.FindOne)
	articlesGroup.Use(authenticate, authorize)
	{
		articlesGroup.POST("", authenticate, articleController.Create)
		articlesGroup.PATCH("/:id", articleController.Update)
		articlesGroup.DELETE("/:id", articleController.Delete)
	}

	categoiesGroup := v1.Group("/categories")
	categoiesController := controllers.Categories{DB: db}
	categoiesGroup.GET("", categoiesController.FindAll)
	categoiesGroup.GET("/:id", categoiesController.FindOne)
	categoiesGroup.Use(authenticate, authorize)
	{
		categoiesGroup.POST("", categoiesController.Create)
		categoiesGroup.PATCH("/:id", categoiesController.Update)
		categoiesGroup.DELETE("/:id", categoiesController.Delete)
	}

	authGroup := v1.Group("/auth")
	authController := controllers.Auth{DB: db}
	{
		authGroup.POST("/sign-up", authController.Signup)
		authGroup.POST("/sign-in", middleware.Authenticate().LoginHandler)
		authGroup.GET("/profile", authenticate, authController.GetProfile)
		authGroup.PATCH("/profile", authenticate, authController.UpdateProfile)
	}

	dashboradController := controllers.Dashboard{DB: db}
	dashboradGroup := v1.Group("/dashboard")
	dashboradGroup.Use(authenticate, authorize)
	{
		dashboradGroup.GET("", dashboradController.GetInfo)
	}

	usersController := controllers.Users{DB: db}
	usersGroup := v1.Group("users")
	usersGroup.Use(authenticate, authorize)
	{
		usersGroup.GET("", usersController.FindAll)
		usersGroup.POST("", usersController.Create)
		usersGroup.GET("/:id", usersController.FindOne)
		usersGroup.PATCH("/:id", usersController.Update)
		usersGroup.DELETE("/:id", usersController.Delete)
		usersGroup.PATCH("/:id/promote", usersController.Promote) // เปลี่ยนสถานะ role
		usersGroup.PATCH("/:id/demote", usersController.Demote)   // ลดสถานะ role = Member
	}

}
