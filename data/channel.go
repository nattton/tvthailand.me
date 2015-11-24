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

type Channel struct {
	ID          int `gorm:"primary_key"`
	Title       string
	Description string
	Thumbnail   string
	URL         string
	HasShow     bool
	IsOnline    bool `json:"-"`

	CreatedAt time.Time  `json:"-"`
	UpdatedAt time.Time  `json:"-"`
	DeletedAt *time.Time `json:"-"`

	Selected bool `sql:"-" json:"-"`
}

func ChannelsToGOB64(s []Channel) string {
	b := bytes.Buffer{}
	e := gob.NewEncoder(&b)
	err := e.Encode(s)
	if err != nil {
		log.Println(`failed gob Encode`, err)
	}
	return base64.StdEncoding.EncodeToString(b.Bytes())
}

func ChannelsFromGOB64(str string) (s []Channel) {
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

func Channels(db *gorm.DB) (channels []Channel, err error) {
	cachedKey := fmt.Sprintf("Channels")
	redisClient := utils.OpenRedis()
	result, err := redisClient.Get(cachedKey).Result()
	if err != nil || err == redis.Nil {
		err = db.Order("order_display").Find(&channels).Error
		if err == nil {
			for i := range channels {
				channels[i].Thumbnail = ThumbnailURLChannel + channels[i].Thumbnail
			}
			redisClient.Set(cachedKey, ChannelsToGOB64(channels), 0)
		}
	} else {
		channels = ChannelsFromGOB64(result)
	}
	return
}

func ChannelsActive(db *gorm.DB) (channels []Channel, err error) {
	cachedKey := fmt.Sprintf("ChannelsActive")
	redisClient := utils.OpenRedis()
	result, err := redisClient.Get(cachedKey).Result()
	if err != nil || err == redis.Nil {
		err = db.Scopes(ChannelScope).Order("order_display").Find(&channels).Error
		if err == nil {
			for i := range channels {
				channels[i].Thumbnail = ThumbnailURLChannel + channels[i].Thumbnail
			}
			redisClient.Set(cachedKey, ChannelsToGOB64(channels), 0)
		}
	} else {
		channels = ChannelsFromGOB64(result)
	}
	return
}

func GetChannel(db *gorm.DB, id string) (channel Channel, err error) {
	err = db.First(&channel, id).Error
	channel.Thumbnail = ThumbnailURLChannel + channel.Thumbnail
	return
}

func ChannelOptions(db *gorm.DB, selectedID int) (channels []Channel) {
	channels, _ = Channels(db)
	if selectedID > 0 {
		for index := range channels {
			if channels[index].ID == selectedID {
				channels[index].Selected = true
				return
			}
		}
	}
	return
}

func ChannelScope(db *gorm.DB) *gorm.DB {
	return db.Where("is_online = ?", true)
}
