package data

import (
	_ "github.com/go-sql-driver/mysql"
	"github.com/jinzhu/gorm"
	"time"
)

type Channel struct {
	ID          string `gorm:"primary_key"`
	Title       string
	Description string
	Thumbnail   string
	URL         string
	HasShow     bool
	IsOnline    bool `json:"-"`

	CreatedAt time.Time  `json:"-"`
	UpdatedAt time.Time  `json:"-"`
	DeletedAt *time.Time `json:"-"`
}

func GetChannels(db *gorm.DB) (channels []Channel, err error) {
	err = db.Order("order_display").Find(&channels).Error
	for i := range channels {
		channels[i].Thumbnail = ThumbnailURLChannel + channels[i].Thumbnail
	}
	return
}

func GetChannel(db *gorm.DB, id string) (channel Channel, err error) {
	err = db.First(&channel, id).Error
	channel.Thumbnail = ThumbnailURLChannel + channel.Thumbnail
	return
}
