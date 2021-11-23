package controllers

import (
	"math"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
)

type pagingResult struct {
	Page      int `json:"page"`
	Limit     int `json:"limit"`
	PrevPage  int `json:"prevPage"`
	NextPage  int `json:"nextPage"`
	Count     int `json:"count"`
	TotalPage int `json:"totalPage"`
}

type pagination struct {
	c       *gin.Context
	query   *gorm.DB
	recodes interface{}
}

func (p *pagination) paginate() *pagingResult {
	// 1. Get limit , page ?limit=10&page=2
	page, _ := strconv.Atoi(p.c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(p.c.DefaultQuery("limit", "12"))

	// 2. count records

	ch := make(chan int)
	go p.conuntRecords(ch)

	// 3. Find records
	// limit, offset
	// limit => 10

	// page => 1, 1 - 10, offset = 0
	// page => 2, 11 - 20, offset = 10
	// page => 3, 21 - 30, offset = 20

	// offset = skip
	// 1 - 1 = 0 * 0 = 0
	offset := (page - 1) * limit
	p.query.Limit(limit).Offset(offset).Find(p.recodes)

	// 4. total page
	// allData / limit
	// รอผลลัพจาก ch
	count := <-ch
	totalPage := int(math.Ceil(float64(count) / float64(limit)))

	// 5. Find nextPage
	var nextPage int
	if page == totalPage {
		nextPage = totalPage
	} else {
		nextPage = page + 1
	}

	// 6. create pagingResult
	return &pagingResult{
		Page:      page,
		Limit:     limit,
		Count:     count,
		PrevPage:  page - 1,
		NextPage:  nextPage,
		TotalPage: totalPage,
	}
}

func (p *pagination) conuntRecords(ch chan int) {
	var count int
	p.query.Model(p.recodes).Count(&count)

	ch <- count
}
