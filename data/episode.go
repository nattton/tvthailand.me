package data

import (
	"github.com/code-mobi/tvthailand.me/Godeps/_workspace/src/github.com/jinzhu/gorm"
	"strconv"
	"strings"
	"time"
)

type Episode struct {
	ID        int `gorm:"primary_key"`
	HashID    string
	ShowID    int
	Ep        int
	Title     string
	Video     string
	SrcType   int
	Date      time.Time
	ViewCount int
	Parts     string
	Password  string
	Thumbnail string

	CreatedAt time.Time  `json:"-"`
	UpdatedAt time.Time  `json:"-"`
	DeletedAt *time.Time `json:"-"`

	Playlists []Playlist `sql:"-"`
	IsURL     bool       `sql:"-"`
}

type Video struct {
	ID        string
	Thumbnail string
	URL       string
}

type Playlist struct {
	Title    string   `json:"title"`
	Image    string   `json:"image"`
	Sources  []Source `json:"sources"`
	Password string   `json:"password"`
	File     string   `json:"-"`
}

type Source struct {
	File string `json:"file"`
}

func EncryptEpisode(db *gorm.DB, episodeID int) {
	var episodes []Episode
	if episodeID > 0 {
		episode := Episode{}
		db.First(&episode, episodeID)
		episodes = append(episodes, episode)
	} else {
		db.Where("hash_id = ?", "").Order("id desc").Find(&episodes)
	}
	for _, episode := range episodes {
		episode.HashID = Encrypt(strconv.Itoa(episode.ID))
		CreatThumbnail(&episode)
		db.Save(&episode)
	}
}

func CreatThumbnail(episode *Episode) {
	videos := strings.Split(strings.Trim(episode.Video, ","), ",")
	var videoID string
	if len(videos) > 0 {
		videoID = videos[0]
	}
	switch episode.SrcType {
	case 0:
		episode.Thumbnail = "https://i.ytimg.com/vi/" + videoID + "/0.jpg"
	case 1:
		episode.Thumbnail = "http://www.dailymotion.com/thumbnail/video/" + videoID
	case 13, 14, 15:
		episode.Thumbnail = "http://video.mthai.com/thumbnail/" + videoID + ".jpg"
	default:
		episode.Thumbnail = "http://thumbnail.instardara.com/chrome.jpg"
	}
}

func GetEpisodes(db *gorm.DB, id int) (episodes []Episode) {
	db.Where("banned = 0 AND show_id = ?", id).Order("ep desc, id desc").Limit(40).Find(&episodes)
	for index := range episodes {
		GetEpisodeTitle(&episodes[index])
	}
	return
}

func GetEpisode(db *gorm.DB, id int) (episode Episode, err error) {
	err = db.First(&episode, id).Error
	SetVideoList(db, &episode)
	return
}

func GetVideoList(db *gorm.DB, hashID string) (episode Episode, err error) {
	db.Where("hash_id = ?", hashID).First(&episode)
	SetVideoList(db, &episode)
	return
}

func SetVideoList(db *gorm.DB, episode *Episode) {
	GetEpisodeTitle(episode)
	videos := strings.Split(strings.Trim(episode.Video, ","), ",")
	lengthVideo := len(videos)
	for i := range videos {
		playlist := Playlist{}
		playlist.Title = episode.Title
		if lengthVideo > 1 {
			playlist.Title += " Part " + strconv.Itoa(i+1) + "/" + strconv.Itoa(lengthVideo)
		}
		videoID := videos[i]
		source := Source{}
		switch episode.SrcType {
		case 0:
			playlist.Image = "https://i.ytimg.com/vi/" + videoID + "/0.jpg"
			source.File = "https://www.youtube.com/watch?v=" + videoID
		case 1:
			playlist.Image = "http://www.dailymotion.com/thumbnail/video/" + videoID
			source.File = "http://www.dailymotion.com/embed/video/" + videoID
		case 14:
			playlist.Image = "http://video.mthai.com/thumbnail/" + videoID + ".jpg"
			if embedVideo := GetEmbedVideo(db, videoID); embedVideo.EmbedURL != "" {
				source.File = embedVideo.EmbedURL
			} else {
				episode.IsURL = true
				source.File = "http://video.mthai.com/cool/player/" + videoID + ".html"
			}
		case 13, 15:
			playlist.Image = "http://video.mthai.com/thumbnail/" + videoID + ".jpg"
			playlist.Password = episode.Password
			source.File = "http://video.mthai.com/cool/player/" + videoID + ".html"
			episode.IsURL = true
		default:
			episode.IsURL = true
			playlist.Image = "http://thumbnail.instardara.com/chrome.jpg"
			episode.Thumbnail = "http://thumbnail.instardara.com/chrome.jpg"
			source.File = videoID
		}
		playlist.File = source.File
		playlist.Sources = append(playlist.Sources, source)
		episode.Playlists = append(episode.Playlists, playlist)
	}
	return
}

func GetEpisodeTitle(episode *Episode) {
	var title string
	if episode.Ep < 20000000 {
		title = "EP." + strconv.Itoa(episode.Ep)
		if episode.Title != "" {
			episode.Title = title + " - " + episode.Title
		} else {
			episode.Title = title
		}
	} else {
		title = "วันที่ " + episode.Date.Format(DateLongFMT)
		if episode.Title != "" {
			episode.Title = episode.Title + " - " + title
		} else {
			episode.Title = title
		}
	}
	return
}
