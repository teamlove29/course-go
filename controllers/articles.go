package controllers

import (
	"course-go/models"
	"mime/multipart"
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/jinzhu/copier"
	"github.com/jinzhu/gorm"
)

type Articles struct {
	DB *gorm.DB
}

type createArticleForm struct {
	Title   string                `form:"title" binding : "required"`
	Body    string                `form:"body" binding : "required"`
	Excerpt string                `form:"excerpt" binding : "required"`
	Image   *multipart.FileHeader `form:"image"binding : "required"`
}

type updateArticleForm struct {
	Title      string                `form:"title"`
	Body       string                `form:"body"`
	Excerpt    string                `form:"excerpt"`
	Image      *multipart.FileHeader `form:"image"`
	CategoryID uint                  `form:"categoryId"`
}

type articleResponse struct {
	ID         uint   `json:"id"`
	Title      string `json:"title"`
	Excerpt    string `json:"excerpt"`
	Body       string `json:"body"`
	Image      string `json:"image"`
	CategoryID uint   `json:"categoryId"`
	Category   struct {
		ID   uint   `json: "id"`
		Name string `json: "name"`
	} `json:"category"`
	User struct {
		Name   string `json: "name"`
		Avatar string `json: "avatar"`
	} `json:"user"`
}

type createOrUpdateResponse struct {
	ID         uint   `json:"id"`
	Title      string `json:"title"`
	Excerpt    string `json:"excerpt"`
	Body       string `json:"body"`
	Image      string `json:"image"`
	CategoryID uint   `json:"categoryId"`
	UserID     uint   `json:"userId"`
}

type articlesPaging struct {
	Items  []articleResponse `json:"items"`
	Paging *pagingResult     `json:"paging"`
}

// /api/v1/articles?categoryID=1&term=aut
func (a *Articles) FindAll(c *gin.Context) {
	var articles []models.Article // empty slice

	query := a.DB.Preload("User").Preload("Category").Order("id desc")

	categoryID := c.Query("categoryId")
	if categoryID != "" {
		query = query.Where("category_id = ?", categoryID)
	}

	term := c.Query("term")
	if term != "" {
		query = query.Where("title ILIKE ?", "%"+term+"%")
	}

	// a.DB.Find(&article)
	// article => limit => 12, page => 1
	// articles?limit => limit => 10, page => 1
	// articles?page=10 => limit => 12, page => 10
	// articles?page=2&linit=4 => limit => 4, page => 2
	pagination := pagination{c: c, query: query, recodes: &articles}
	paging := pagination.paginate()

	serializedArticles := []articleResponse{} // nil slice
	copier.Copy(&serializedArticles, &articles)

	c.JSON(http.StatusOK, gin.H{"articles": articlesPaging{Items: serializedArticles, Paging: paging}})
}

func (a *Articles) FindOne(c *gin.Context) {
	article, err := a.findArticleByID(c)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	serializedArticle := articleResponse{}
	copier.Copy(&serializedArticle, &article)
	c.JSON(http.StatusOK, gin.H{"article": serializedArticle})
}

func (a *Articles) Create(c *gin.Context) {

	var form createArticleForm
	var article models.Article

	user, _ := c.Get("sub")
	copier.Copy(&article, &form)
	article.User = *user.(*models.User)

	if err := c.ShouldBind(&form); err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{"error": err.Error()})
		return
	}

	copier.Copy(&article, &form)

	if err := a.DB.Create(&article).Error; err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{"error": err.Error()})
		return
	}

	a.setArticleImage(c, &article)
	serializedArticle := createOrUpdateResponse{}
	copier.Copy(&serializedArticle, &article)

	c.JSON(http.StatusCreated, gin.H{"article": serializedArticle})

}

func (a *Articles) setArticleImage(c *gin.Context, article *models.Article) error {
	file, err := c.FormFile("image")
	if err != nil || file == nil {
		return err
	}

	// เช็คว่ามีรูปอยู่ในไฟล์อยู่รึเปล่า
	if article.Image != "" {
		// http://127.0.0.1:5000/uploads/articles/<ID>/image.png
		// 1. /uploads/articles/<ID>/image.png
		article.Image = strings.Replace(article.Image, os.Getenv("HOST"), "", 1)
		// 2. <WD>/uploads/articles/<ID>/image.png <WD> = working dir
		pwd, _ := os.Getwd()
		// 3. Remove <WD>/uploads/articles/<ID>/image.png
		os.Remove(pwd + article.Image)
	}

	// create Path
	// uploads/articles/123
	path := "uploads/articles/" + strconv.Itoa(int(article.ID))
	os.MkdirAll(path, 0755)

	// Upload file
	// uploads/articles/123/filename
	filename := path + "/" + file.Filename
	if err := c.SaveUploadedFile(file, filename); err != nil {
		return err
	}
	// Attach file to article
	article.Image = os.Getenv("HOST") + "/" + filename
	// update to sql
	a.DB.Save(article)

	return nil
}

func (a *Articles) findArticleByID(c *gin.Context) (*models.Article, error) {
	var article models.Article

	id := c.Param("id")

	// SELECT * FROM articles;
	// SELECT * FROM Users WHERE user_id = user_id;
	// SELECT * FROM Categories WHERE category_id = category_id;
	if err := a.DB.Preload("User").Preload("Category").First(&article, id).Error; err != nil {
		return nil, err
	}

	return &article, nil

}

func (a *Articles) Update(c *gin.Context) {
	var form updateArticleForm

	if err := c.ShouldBind(&form); err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{"error": err.Error()})
		return
	}

	article, err := a.findArticleByID(c)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	if err := a.DB.Model(&article).Update(&form).Error; err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{"error": err.Error()})
		return
	}

	a.setArticleImage(c, article)

	var serializedArticle createOrUpdateResponse
	copier.Copy(&serializedArticle, &article)
	c.JSON(http.StatusOK, gin.H{"article": serializedArticle})

}

func (a *Articles) Delete(c *gin.Context) {

	article, err := a.findArticleByID(c)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	a.DB.Unscoped().Delete(&article)
	// a.DB.Delete(&article)

	c.Status(http.StatusNoContent)

}
