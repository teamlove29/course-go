package config

import (
	"os"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
)

var db *gorm.DB

type Product struct {
	gorm.Model
	Code  string
	Price uint
}

func InitDB() {

	// db, err := gorm.Open(postgres.Open(os.Getenv("DATABASE_CONNECTION")), &gorm.Config{})
	var err error
	db, err = gorm.Open("postgres", os.Getenv("DATABASE_URL"))
	if err != nil {
		panic("failed to connect database")
	}
	db.AutoMigrate(&Product{})
	// db.LogMode(true)
	// db.LogMode(gin.Mode() == gin.DebugMode)
}

func GetDB() *gorm.DB {
	return db
}

func CloseDB() {
	db.Close()
}
