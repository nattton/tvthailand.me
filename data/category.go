package data

import (
	"bytes"
	"encoding/base64"
	"encoding/gob"
	"fmt"
	"github.com/code-mobi/tvthailand.me/Godeps/_workspace/src/github.com/jinzhu/gorm"
	"github.com/code-mobi/tvthailand.me/utils"
	"log"
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

func CategoriesToGOB64(s []Category) string {
	b := bytes.Buffer{}
	e := gob.NewEncoder(&b)
	err := e.Encode(s)
	if err != nil {
		log.Println(`failed gob Encode`, err)
	}
	return base64.StdEncoding.EncodeToString(b.Bytes())
}

func CategoriesFromGOB64(str string) (s []Category) {
	by, err := base64.StdEncoding.DecodeString(str)
	if err != nil {
		log.Println(`failed base64 Decode`, err)
	}
	b := bytes.Buffer{}
	b.Write(by)
	d := gob.NewDecoder(&b)
	err = d.Decode(&s)
	if err != nil {
		log.Println(`failed gob Decode`, err)
	}
	return
}

func CategoriesActive(db *gorm.DB) (categories []Category, err error) {
	cachedKey := fmt.Sprintf("CategoriesActive")
	redisClient := utils.OpenRedis()
	result, err := redisClient.Get(cachedKey).Result()
	if err != nil {
		err = db.Scopes(CategoryScope).Find(&categories).Error
		for i := range categories {
			categories[i].Thumbnail = ThumbnailURLCategory + categories[i].Thumbnail
		}
		redisClient.Set(cachedKey, CategoriesToGOB64(categories), 24*time.Hour)
	} else {
		categories = CategoriesFromGOB64(result)
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

func CategoryScope(db *gorm.DB) *gorm.DB {
	return db.Where("is_online = ?", true)
}
