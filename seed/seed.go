package seed

import (
	"course-go/config"
	"course-go/migrations"
	"course-go/models"
	"math/rand"
	"strconv"

	"github.com/bxcodec/faker/v3"
	log "github.com/sirupsen/logrus"
)

func Load() {
	db := config.GetDB()

	// Clean Database
	db.DropTableIfExists("articles", "categories", "migrations")
	migrations.Migrate()

	// Add categoies
	log.Info("Creating categories...")

	numOfcategories := 20
	categories := make([]models.Category, 0, numOfcategories)

	for i := 1; i <= numOfcategories; i++ {
		category := models.Category{
			Name: faker.Word(),
			Desc: faker.Paragraph(),
		}
		db.Create(&category)
		categories = append(categories, category)
	}

	// Add articles
	log.Info("Creating articles...")

	numOfArticles := 50
	articles := make([]models.Article, 0, numOfArticles)

	for i := 1; i <= numOfArticles; i++ {
		article := models.Article{
			Title:      faker.Sentence(),
			Excerpt:    faker.Sentence(),
			Body:       faker.Paragraph(),
			Image:      "https://source.unsplash.com/random/300x200?" + strconv.Itoa(i),
			CategoryID: uint(rand.Intn(numOfcategories) + 1),
		}
		db.Create(&article)
		articles = append(articles, article)
	}

}
