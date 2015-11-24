package data

import (
	"bytes"
	"encoding/base64"
	"encoding/gob"
	"fmt"
	"log"
	"time"

	"github.com/code-mobi/tvthailand.me/Godeps/_workspace/src/github.com/jinzhu/gorm"
	"github.com/code-mobi/tvthailand.me/Godeps/_workspace/src/gopkg.in/redis.v3"
	"github.com/code-mobi/tvthailand.me/utils"
)

type Category struct {
	ID           int `gorm:"primary_key"`
	Title        string
	Titleize     string
	Description  string
	Thumbnail    string
	OrderDisplay int  `json:"-"`
	IsOnline     bool `json:"-"`

	CreatedAt time.Time  `json:"-"`
	UpdatedAt time.Time  `json:"-"`
	DeletedAt *time.Time `json:"-"`

	Selected bool `sql:"-" json:"-"`
}

func (s Category) ToGOB64() string {
	b := bytes.Buffer{}
	e := gob.NewEncoder(&b)
	err := e.Encode(s)
	if err != nil {
		log.Println(`failed gob Encode`, err)
	}
	return base64.StdEncoding.EncodeToString(b.Bytes())
}

func (s *Category) FromGOB64(str string) {
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

func Categories(db *gorm.DB) (categories []Category, err error) {
	cachedKey := fmt.Sprintf("Categories")
	redisClient := utils.OpenRedis()
	result, err := redisClient.Get(cachedKey).Result()
	if err != nil || err == redis.Nil {
		err = db.Find(&categories).Error
		if err == nil {
			for i := range categories {
				categories[i].Thumbnail = ThumbnailURLCategory + categories[i].Thumbnail
			}
			redisClient.Set(cachedKey, CategoriesToGOB64(categories), 0)
		}
	} else {
		categories = CategoriesFromGOB64(result)
	}
	return
}

func CategoriesActive(db *gorm.DB) (categories []Category, err error) {
	cachedKey := fmt.Sprintf("CategoriesActive")
	redisClient := utils.OpenRedis()
	result, err := redisClient.Get(cachedKey).Result()
	if err != nil || err == redis.Nil {
		err = db.Scopes(CategoryScope).Find(&categories).Error
		if err == nil {
			for i := range categories {
				categories[i].Thumbnail = ThumbnailURLCategory + categories[i].Thumbnail
			}
			redisClient.Set(cachedKey, CategoriesToGOB64(categories), 0)
		}
	} else {
		categories = CategoriesFromGOB64(result)
	}
	return
}

func CategoryTitleze(db *gorm.DB, titlize string) (category Category, err error) {
	cachedKey := fmt.Sprintf("GetCategory/titlize=%s", titlize)
	redisClient := utils.OpenRedis()
	result, err := redisClient.Get(cachedKey).Result()
	if err != nil || err == redis.Nil {
		err = db.Where("titleize = ?", titlize).First(&category).Error
		if err == nil {
			category.Thumbnail = ThumbnailURLCategory + category.Thumbnail
			redisClient.Set(cachedKey, category.ToGOB64(), 0)
		}
	} else {
		category.FromGOB64(result)
	}
	return
}

func CategoryOptions(db *gorm.DB, selectedID int) (categories []Category) {
	categories, _ = Categories(db)
	if selectedID > 0 {
		for index := range categories {
			if categories[index].ID == selectedID {
				categories[index].Selected = true
				return
			}
		}
	}
	return
}

func CategoryScope(db *gorm.DB) *gorm.DB {
	return db.Where("is_online = ?", true)
}
