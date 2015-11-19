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

type Show struct {
	ID          int       `gorm:"primary_key" json:"id"`
	CategoryID  int       `json:"-"`
	ChannelID   int       `json:"-"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	Thumbnail   string    `json:"thumbnail"`
	Poster      string    `json:"-"`
	Detail      string    `json:"-"`
	LastEpname  string    `json:"-"`
	ViewCount   int       `json:"-"`
	Rating      float32   `json:"-"`
	VoteCount   int       `json:"-"`
	IsOtv       bool      `json:"-"`
	OtvID       string    `json:"-"`
	UpdateDate  time.Time `json:"-"`

	CreatedAt time.Time  `json:"-"`
	UpdatedAt time.Time  `json:"-"`
	DeletedAt *time.Time `json:"-"`
}

func (s Show) ToGOB64() string {
	b := bytes.Buffer{}
	e := gob.NewEncoder(&b)
	err := e.Encode(s)
	if err != nil {
		log.Println(`failed gob Encode`, err)
	}
	return base64.StdEncoding.EncodeToString(b.Bytes())
}

func (s *Show) FromGOB64(str string) {
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

func ShowsToGOB64(s []Show) string {
	b := bytes.Buffer{}
	e := gob.NewEncoder(&b)
	err := e.Encode(s)
	if err != nil {
		log.Println(`failed gob Encode`, err)
	}
	return base64.StdEncoding.EncodeToString(b.Bytes())
}

func ShowsFromGOB64(str string) (s []Show) {
	by, err := base64.StdEncoding.DecodeString(str)
	if err != nil {
		fmt.Println(`failed base64 Decode`, err)
	}
	b := bytes.Buffer{}
	b.Write(by)
	d := gob.NewDecoder(&b)
	err = d.Decode(&s)
	if err != nil {
		fmt.Println(`failed gob Decode`, err)
	}
	return
}

func GetShow(db *gorm.DB, id int) (show Show, err error) {
	cachedKey := fmt.Sprintf("Show/id=%d", id)
	redisClient := utils.OpenRedis()
	result, err := redisClient.Get(cachedKey).Result()
	if err != nil || err == redis.Nil {
		err = db.First(&show, id).Error
		if err == nil {
			show.Thumbnail = ThumbnailURLTv + show.Thumbnail
			redisClient.Set(cachedKey, show.ToGOB64(), 0)
		}
	} else {
		show.FromGOB64(result)
	}
	return
}

func ShowWithOtv(db *gorm.DB, id int) (show Show, err error) {
	cachedKey := fmt.Sprintf("ShowWithOtv/otv_id=%d", id)
	redisClient := utils.OpenRedis()
	result, err := redisClient.Get(cachedKey).Result()
	if err != nil || err == redis.Nil {
		err = db.Where("otv_id = ?", id).First(&show).Error
		if err == nil {
			show.Thumbnail = ThumbnailURLTv + show.Thumbnail
			redisClient.Set(cachedKey, show.ToGOB64(), 0)
		}
	} else {
		show.FromGOB64(result)
	}
	return
}

func ShowsRecently(db *gorm.DB, offset int) (shows []Show, err error) {
	cachedKey := fmt.Sprintf("ShowsRecently/offset=%d", offset)
	redisClient := utils.OpenRedis()
	result, err := redisClient.Get(cachedKey).Result()
	if err != nil || err == redis.Nil {
		err = db.Scopes(ShowScope).Order("update_date desc").Offset(offset).Limit(20).Find(&shows).Error
		if err == nil {
			for i := range shows {
				shows[i].Thumbnail = ThumbnailURLTv + shows[i].Thumbnail
			}
			redisClient.Set(cachedKey, ShowsToGOB64(shows), 5*time.Minute)
		}
	} else {
		shows = ShowsFromGOB64(result)
	}
	return
}

func ShowsPopular(db *gorm.DB, offset int) (shows []Show, err error) {
	cachedKey := fmt.Sprintf("ShowsPopular/offset=%d", offset)
	redisClient := utils.OpenRedis()
	result, err := redisClient.Get(cachedKey).Result()
	if err != nil || err == redis.Nil {
		err = db.Scopes(ShowScope).Order("view_count desc").Offset(offset).Limit(20).Find(&shows).Error
		if err == nil {
			for i := range shows {
				shows[i].Thumbnail = ThumbnailURLTv + shows[i].Thumbnail
			}
			redisClient.Set(cachedKey, ShowsToGOB64(shows), 0)
		}
	} else {
		shows = ShowsFromGOB64(result)
	}
	return
}

func ShowsCategory(db *gorm.DB, id string, offset int) (shows []Show, err error) {
	cachedKey := fmt.Sprintf("ShowsCategory/id=%s/offset=%d", id, offset)
	redisClient := utils.OpenRedis()
	result, err := redisClient.Get(cachedKey).Result()
	if err != nil || err == redis.Nil {
		err = db.Scopes(ShowScope).Where("category_id = ?", id).Order("update_date desc").Offset(offset).Limit(20).Find(&shows).Error
		if err == nil {
			for i := range shows {
				shows[i].Thumbnail = ThumbnailURLTv + shows[i].Thumbnail
			}
			redisClient.Set(cachedKey, ShowsToGOB64(shows), 10*time.Minute)
		}
	} else {
		shows = ShowsFromGOB64(result)
	}
	return
}

func ShowsCategoryPopular(db *gorm.DB, id string) (shows []Show, err error) {
	err = db.Scopes(ShowScope).Where("category_id = ?", id).Order("view_count desc").Limit(20).Find(&shows).Error
	for i := range shows {
		shows[i].Thumbnail = ThumbnailURLTv + shows[i].Thumbnail
	}
	return
}

func ShowsChannel(db *gorm.DB, id string, offset int) (shows []Show, err error) {
	cachedKey := fmt.Sprintf("ShowsChannel/id=%s/offset=%d", id, offset)
	redisClient := utils.OpenRedis()
	result, err := redisClient.Get(cachedKey).Result()
	if err != nil || err == redis.Nil {
		err = db.Scopes(ShowScope).Where("channel_id = ?", id).Order("update_date desc").Offset(offset).Limit(20).Find(&shows).Error
		if err == nil {
			for i := range shows {
				shows[i].Thumbnail = ThumbnailURLTv + shows[i].Thumbnail
			}
			redisClient.Set(cachedKey, ShowsToGOB64(shows), 10*time.Minute)
		}
	} else {
		shows = ShowsFromGOB64(result)
	}
	return
}

func ShowsChannelPopular(db *gorm.DB, id string) (shows []Show, err error) {
	err = db.Scopes(ShowScope).Where("channel_id = ?", id).Order("view_count desc").Limit(20).Find(&shows).Error
	for i := range shows {
		shows[i].Thumbnail = ThumbnailURLTv + shows[i].Thumbnail
	}
	return
}

func ShowsSearch(db *gorm.DB, keyword string) (shows []Show, err error) {
	db.Scopes(ShowScope).Where("title LIKE ?", "%"+keyword+"%").Order("update_date desc, title asc").Limit(20).Find(&shows)
	for i := range shows {
		shows[i].Thumbnail = ThumbnailURLTv + shows[i].Thumbnail
	}
	return
}

func ShowUpdateDate(db *gorm.DB, id int) {
	db.Model(Show{}).Where("id = ?", id).
		Updates(Show{UpdateDate: time.Now(), UpdatedAt: time.Now()})
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
