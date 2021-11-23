package main

import (
	"course-go/config"
	"course-go/migrations"
	"course-go/routes"
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	if err := godotenv.Load(); err != nil {
		// if err print and stop.
		log.Fatal("Error loading .env file")
	}

	config.InitDB()
	defer config.CloseDB()
	migrations.Migrate()

	// create faker database
	// seed.Load()

	r := gin.Default()
	r.Static("/uploads", "./uploads")

	uploadDirs := [...]string{"articles", "users"}

	for _, dir := range uploadDirs {
		os.MkdirAll("uploads/"+dir, 0755)
	}

	routes.Serve(r)

	r.Run(":" + os.Getenv("PORT"))
}
