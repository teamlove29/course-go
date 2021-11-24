package migrations

import (
	"course-go/config"
	"log"

	"gopkg.in/gormigrate.v1"
)

func Migrate() {
	db := config.GetDB()
	m := gormigrate.New(
		db,
		gormigrate.DefaultOptions,
		[]*gormigrate.Migration{
			m1636648144CreateArticlesTable(),
			m1637308954CreateCategoiesTable(),
			m1637321814AddCategoryIdToArticles(),
			m1637685002CreateUsersTable(),
			m1637774714AddUserIdToArticles(),
		},
	)

	if err := m.Migrate(); err != nil {
		log.Fatalf("Could not migrate: %v", err)
	}
}
