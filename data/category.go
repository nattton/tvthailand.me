package data

import (
	_ "github.com/go-sql-driver/mysql"
	"github.com/jinzhu/gorm"
	"time"
)

type Category struct {
	ID           string `gorm:"primary_key"`
	Title        string
	Titleize     string
	Description  string
	Thumbnail    string
	OrderDisplay int  `json:"-"`
	IsOnline     bool `json:"-"`

	CreatedAt time.Time  `json:"-"`
	UpdatedAt time.Time  `json:"-"`
	DeletedAt *time.Time `json:"-"`
}

var CategoriesMenu []Category

func GetCategories(db *gorm.DB) (categories []Category, err error) {
	err = db.Find(&categories).Error
	for i := range categories {
		categories[i].Thumbnail = ThumbnailURLCategory + categories[i].Thumbnail
	}
	return
}

func GetCategory(db *gorm.DB, titlize string) (category Category, err error) {
	err = db.Where("titleize = ?", titlize).First(&category).Error
	if err != nil {
		return
	}
	category.Thumbnail = ThumbnailURLCategory + category.Thumbnail
	return
}
