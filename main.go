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

	// คำสั่ง run heroku
	// install heroku
	// heroku --version
	// heroku login
	// heroku create = จะได้้ url
	// git push heroku main (อัพ git ขึ้นไปที่ heroku)
	// create file: Procfile
	// git push heroku main

	// add tool to heruko
	// heroku addons:create heroku-postgresql:hobby-dev
	// heroku config
	// change .env DATABASE_CONNECTION to DATABASE_URL
	// git push heroku main

	// heroku ไม่อ่าน .env จึงต้อง setting
	// heroku config:set SECRET_KEY=$(uuidgen)
	// heroku config:set HOST=https://lit-tundra-22431.herokuapp.com
	// heroku config:set GIN_MODE=release APP_ENV=production (ใช้แค่ ginmode ก็ได้)
	// heroku open ทำการ test ได้เลย

	// fix error git pull --rebase

	if os.Getenv("APP_ENV") != "production" {
		err := godotenv.Load()
		if err != nil {
			log.Fatal("Error loading .env file")
		}
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
