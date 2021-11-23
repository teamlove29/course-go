package migrations

import (
	"course-go/models"

	"github.com/jinzhu/gorm"
	"gopkg.in/gormigrate.v1"
)

func m1637308954CreateCategoiesTable() *gormigrate.Migration {
	return &gormigrate.Migration{
		ID: "1637308954",
		Migrate: func(tx *gorm.DB) error {
			return tx.AutoMigrate(&models.Category{}).Error
		},
		Rollback: func(tx *gorm.DB) error {
			return tx.DropTable("cetagories").Error
		},
	}
}
