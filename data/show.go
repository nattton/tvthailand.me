package data

import (
	"github.com/code-mobi/tvthailand.me/Godeps/_workspace/src/github.com/jinzhu/gorm"
	"time"
)

type Show struct {
	ID          int     `gorm:"primary_key" json:"id"`
	CategoryID  int     `json:"-"`
	ChannelID   int     `json:"-"`
	Title       string  `json:"title"`
	Description string  `json:"description"`
	Thumbnail   string  `json:"thumbnail"`
	Poster      string  `json:"-"`
	Detail      string  `json:"-"`
	LastEpname  string  `json:"-"`
	ViewCount   int     `json:"-"`
	Rating      float32 `json:"-"`
	VoteCount   int     `json:"-"`
	IsOtv       bool    `json:"-"`
	OtvID       string  `json:"-"`

	CreatedAt time.Time  `json:"-"`
	UpdatedAt time.Time  `json:"-"`
	DeletedAt *time.Time `json:"-"`
}

func GetShow(db *gorm.DB, id int) (show Show, err error) {
	err = db.First(&show, id).Error
	show.Thumbnail = ThumbnailURLTv + show.Thumbnail
	return
}

func GetShowByOtv(db *gorm.DB, id int) (show Show, err error) {
	err = db.Where("otv_id = ?", id).First(&show).Error
	return
}

func GetShowByRecently(db *gorm.DB, offset int) (shows []Show, err error) {
	err = db.Scopes(ShowScope).Order("update_date desc").Offset(offset).Limit(20).Find(&shows).Error
	for i := range shows {
		shows[i].Thumbnail = ThumbnailURLTv + shows[i].Thumbnail
	}
	return
}

func GetShowByPopular(db *gorm.DB, offset int) (shows []Show, err error) {
	err = db.Scopes(ShowScope).Order("view_count desc").Offset(offset).Limit(20).Find(&shows).Error
	for i := range shows {
		shows[i].Thumbnail = ThumbnailURLTv + shows[i].Thumbnail
	}
	return
}

func GetShowByCategory(db *gorm.DB, id string, offset int) (shows []Show, err error) {
	err = db.Scopes(ShowScope).Where("category_id = ?", id).Order("update_date desc").Offset(offset).Limit(20).Find(&shows).Error
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

func GetShowByChannel(db *gorm.DB, id string, offset int) (shows []Show, err error) {
	err = db.Scopes(ShowScope).Where("channel_id = ?", id).Order("update_date desc").Offset(offset).Limit(20).Find(&shows).Error
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

func GetShowBySearch(db *gorm.DB, keyword string) (shows []Show, err error) {
	db.Scopes(ShowScope).Where("title LIKE ?", "%"+keyword+"%").Order("update_date desc, title asc").Limit(20).Find(&shows)
	for i := range shows {
		shows[i].Thumbnail = ThumbnailURLTv + shows[i].Thumbnail
	}
	return
}

func ShowScope(db *gorm.DB) *gorm.DB {
	return db.Where("is_online = ? AND build_max > ?", true, 1000)
}

func ResetShowViewCount(db *gorm.DB) (err error) {
	err = db.Model(Show{}).UpdateColumn("view_count", 0).Error
	return
}

func UpdateShowViewCount(db *gorm.DB, title string, viewCount int) int64 {
	return db.Model(Show{}).Where("title = ?", title).UpdateColumn("view_count", viewCount).RowsAffected
}
