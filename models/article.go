package models

import "gorm.io/gorm"

// type Article struct {
// 	ID    uint   `json:"id"`
// 	Title string `json:"title"`
// 	Body  string `json:"body"`
// 	Image string `json:"image"`
// }

type Article struct {
	gorm.Model
	Title      string `gorm:"unique;not null"` // unique == ไม่ซ้ำ
	Excerpt    string `gorm:"not null"`
	Body       string `gorm:"not null"`
	Image      string `gorm:"not null"`
	CategoryID uint
}
