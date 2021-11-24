package migrations

import (
	"course-go/models"

	"github.com/jinzhu/gorm"
	"gopkg.in/gormigrate.v1"
)

func m1637774714AddUserIdToArticles() *gormigrate.Migration {
	return &gormigrate.Migration{
		ID: "1637774714",
		Migrate: func(tx *gorm.DB) error {
			return tx.AutoMigrate(&models.Article{}).Error
		},
		Rollback: func(tx *gorm.DB) error {
			return tx.Model(&models.Article{}).DropTable("user_id").Error
		},
	}
}
