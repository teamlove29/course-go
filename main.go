package main

import (
	"course-go/config"
	"course-go/migrations"
	"course-go/routes"
	"log"
	"os"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {

	// export APP_ENV=production

	// export HOST=http://127.0.0.1:8080
	// export PORT=8080
	// uuidgen
	// export SECRET_KEY=$(uuidgen)
	// go build
	// ./course-go
	// ถ้ารันบน server เอาไฟล์ binary ไปรันก็พอ

	if os.Getenv("APP_ENV") != "production" {
		err := godotenv.Load()
		if err != nil {
			log.Fatal("Error loading .env file")
		}
	}

	if err := godotenv.Load(); err != nil {
		// if err print and stop.
		log.Fatal("Error loading .env file")
	}

	config.InitDB()
	defer config.CloseDB()
	migrations.Migrate()

	// AllowAll http
	corsConfig := cors.DefaultConfig()
	corsConfig.AllowAllOrigins = true
	corsConfig.AddAllowHeaders("Authorization")

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
