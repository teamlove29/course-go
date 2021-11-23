package migrations

import (
	"course-go/models"

	"github.com/jinzhu/gorm"
	"gopkg.in/gormigrate.v1"
)

// func m1636648144CreateArticlesTable() *gormigrate.Migration {
// 	return &gormigrate.Migration{
// 		ID: "1636648144",
// 		Migrate: func(tx *gorm.DB) error {

// 			return tx.AutoMigrate(&models.Article{}).Error
// 		},
// 		Rollback: func(tx *gorm.DB) error {
// 			return tx.Migrator().DropTable("articles")
// 		},
// 	}
// }

func m1636648144CreateArticlesTable() *gormigrate.Migration {
	return &gormigrate.Migration{
		ID: "1596813596",
		Migrate: func(tx *gorm.DB) error {
			return tx.AutoMigrate(&models.Article{}).Error
		},
		Rollback: func(tx *gorm.DB) error {
			return tx.DropTable("articles").Error
		},
	}
}
