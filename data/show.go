package data

import (
	_ "github.com/code-mobi/tvthailand.me/Godeps/_workspace/src/github.com/go-sql-driver/mysql"
	"github.com/code-mobi/tvthailand.me/Godeps/_workspace/src/github.com/jinzhu/gorm"
	"time"
)

type Show struct {
	ID          int `gorm:"primary_key"`
	Title       string
	Description string
	Thumbnail   string
	Poster      string
	Detail      string
	LastEpname  string
	ViewCount   int
	Rating      float32
	VoteCount   int
	IsOtv       bool
	OtvID       string

	CreatedAt time.Time  `json:"-"`
	UpdatedAt time.Time  `json:"-"`
	DeletedAt *time.Time `json:"-"`
}

func GetShow(db *gorm.DB, id int) (show Show, err error) {
	err = db.First(&show, id).Error
	return
}

func GetShowByOtv(db *gorm.DB, id int) (show Show, err error) {
	err = db.Where("otv_id = ?", id).First(&show).Error
	return
}

func GetShowByRecently(db *gorm.DB, start int) (shows []Show, err error) {
	err = db.Scopes(ShowScope).Order("update_date desc").Offset(start).Limit(20).Find(&shows).Error
	for i := range shows {
		shows[i].Thumbnail = ThumbnailURLTv + shows[i].Thumbnail
	}
	return
}

func GetShowByPopular(db *gorm.DB, start int) (shows []Show, err error) {
	err = db.Scopes(ShowScope).Order("view_count desc").Offset(start).Limit(20).Find(&shows).Error
	for i := range shows {
		shows[i].Thumbnail = ThumbnailURLTv + shows[i].Thumbnail
	}
	return
}

func GetShowByCategory(db *gorm.DB, id string, start int) (shows []Show, err error) {
	err = db.Scopes(ShowScope).Where("category_id = ?", id).Order("update_date desc").Offset(start).Limit(20).Find(&shows).Error
	for i := range shows {
		shows[i].Thumbnail = ThumbnailURLTv + shows[i].Thumbnail
	}
	return
}

func GetShowByCategoryPopular(db *gorm.DB, id string) (shows []Show, err error) {
	err = db.Scopes(ShowScope).Where("category_id = ?", id).Order("view_count desc").Limit(20).Find(&shows).Error
	for i := range shows {
		shows[i].Thumbnail = ThumbnailURLTv + shows[i].Thumbnail
	}
	return
}

func GetShowByChannel(db *gorm.DB, id string, start int) (shows []Show, err error) {
	err = db.Scopes(ShowScope).Where("channel_id = ?", id).Order("update_date desc").Offset(start).Limit(20).Find(&shows).Error
	for i := range shows {
		shows[i].Thumbnail = ThumbnailURLTv + shows[i].Thumbnail
	}
	return
}

func GetShowByChannelPopular(db *gorm.DB, id string) (shows []Show, err error) {
	err = db.Scopes(ShowScope).Where("channel_id = ?", id).Order("view_count desc").Limit(20).Find(&shows).Error
	for i := range shows {
		shows[i].Thumbnail = ThumbnailURLTv + shows[i].Thumbnail
	}
	return
}

func GetShowBySearch(db *gorm.DB, keyword string) (shows []Show) {
	db.Scopes(ShowScope).Where("title LIKE ?", "%"+keyword+"%").Order("title asc").Limit(20).Find(&shows)
	for i := range shows {
		shows[i].Thumbnail = ThumbnailURLTv + shows[i].Thumbnail
	}
	return
}

func ShowScope(db *gorm.DB) *gorm.DB {
	return db.Where("is_online = 1 AND build_max > 1000")
}
