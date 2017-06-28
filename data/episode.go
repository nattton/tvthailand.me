package data

import (
	"encoding/base64"
	"math"
	"strconv"
	"strings"
	"time"

	"fmt"

	"github.com/jinzhu/gorm"
)

type EpisodePage struct {
	PageInfo PageInfo  `json:"pageInfo"`
	Episodes []Episode `json:"episodes"`
}

type Episode struct {
	ID               int `gorm:"primary_key"`
	HashID           string
	ShowID           int
	Ep               int
	Title            string
	Video            string
	VideoEncrypt     string
	VideoEncryptPath string
	SrcType          int
	SrcTypePath      int
	Date             time.Time
	ViewCount        int
	Parts            string
	Password         string
	Thumbnail        string
	User             string

	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt *time.Time

	Playlists  []Playlist `sql:"-"`
	IsURL      bool       `sql:"-"`
	Videos     []string   `sql:"-"`
	VideoCount int        `sql:"-"`
}

func (episode *Episode) SetVideo(videoID string) {
	episode.Video = strings.Trim(videoID, ",")
	var re = strings.NewReplacer(
		"+", "-",
		"=", ",",
		"a", "!",
		"b", "@",
		"c", "#",
		"d", "$",
		"e", "%",
		"f", "^",
		"g", "&",
		"h", "*",
		"i", "(",
		"j", ")",
		"k", "{",
		"l", "}",
		"m", "[",
		"n", "]",
		"o", ":",
		"p", ";",
		"q", "<",
		"r", ">",
		"s", "?",
	)
	encrypt := base64.StdEncoding.EncodeToString([]byte(episode.Video))
	episode.VideoEncrypt = re.Replace(encrypt)
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
	VideoID  string   `json:"-"`
}

type Source struct {
	File string `json:"file"`
}

func EncryptAllEpisodes(db *gorm.DB) {
	var episodes []Episode
	db.Where("hash_id = ? OR video_encrypt = ? OR video_encrypt_path = ?", "", "", "").Order("id desc").Find(&episodes)
	for _, episode := range episodes {
		EncryptEpisode(db, &episode)
	}
}

func ReEncryptAllEpisodes(db *gorm.DB) {
	var episodes []Episode
	db.Order("id desc").Find(&episodes)
	for _, episode := range episodes {
		EncryptEpisode(db, &episode)
	}
}

func EncryptEpisode(db *gorm.DB, episode *Episode) {
	episode.VideoEncrypt = EncryptVideo(episode.Video)

	if episode.SrcType == 0 {
		videoArray := strings.Split(episode.Video, ",")
		videoPath := YoutubeViewURL + strings.Join(videoArray, ","+YoutubeViewURL)
		fmt.Println(videoPath)
		episode.VideoEncryptPath = EncryptVideo(videoPath)
		episode.SrcTypePath = 11
	} else {
		episode.VideoEncryptPath = episode.VideoEncrypt
		episode.SrcTypePath = episode.SrcType
	}

	episode.HashID = Encrypt(strconv.Itoa(episode.ID))
	CreateThumbnail(episode)
	db.Save(&episode)
}

func CreateEpisodeMThaiThumbnail(db *gorm.DB, gtID int) {
	var episodes []Episode
	db.Where("src_type in (?) AND id >= ?", []int{13, 14, 15}, gtID).Order("id asc").Find(&episodes)
	for _, episode := range episodes {
		episode.HashID = Encrypt(strconv.Itoa(episode.ID))
		CreateThumbnail(&episode)
		db.Save(&episode)
	}
}

func CreateThumbnail(episode *Episode) {
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
		_, thumbnailURL := GetMThaiEmbedURL(videoID)
		episode.Thumbnail = thumbnailURL
		if thumbnailURL != "" {
			episode.Thumbnail = thumbnailURL
		} else {
			episode.Thumbnail = "http://video.mthai.com/thumbnail/" + videoID + ".jpg"
		}
	default:
		episode.Thumbnail = "http://thumbnail.instardara.com/chrome.jpg"
	}
}

func GetEpisodes(db *gorm.DB, showID int, offset int) (episodes []Episode, err error) {
	err = db.Where("banned = 0 AND show_id = ?", showID).Order("ep desc, id desc").Offset(offset).Limit(20).Find(&episodes).Error
	for index := range episodes {
		GetEpisodeTitle(&episodes[index])
	}
	return
}

func EpisodesAndPageInfo(db *gorm.DB, showID int, page int32) (episodes []Episode, pageInfo PageInfo, err error) {
	if page < 1 {
		page = 1
	}
	currentPage := page
	pageInfo.ResultsPerPage = 20
	page--
	offset := page * pageInfo.ResultsPerPage

	dbQ := db.Table("episodes").Where("banned = 0 AND show_id = ?", showID).Order("ep desc, id desc")
	dbQ.Count(&pageInfo.TotalResults)

	maxPage := int32(math.Ceil(float64(pageInfo.TotalResults) / float64(pageInfo.ResultsPerPage)))
	if currentPage <= maxPage {
		if currentPage > 1 {
			pageInfo.PreviousPage = currentPage - 1
		}
		if currentPage < maxPage {
			pageInfo.NextPage = currentPage + 1
		}
		err = dbQ.Offset(offset).Limit(pageInfo.ResultsPerPage).Find(&episodes).Error
		for index := range episodes {
			GetEpisodeTitle(&episodes[index])
		}
	} else {

	}
	return
}

func GetEpisodesBySearch(db *gorm.DB, keyword string) (episodes []Episode, err error) {
	err = db.Where("banned = 0 AND title LIKE ?", "%"+keyword+"%").Order("ep desc, id desc").Limit(20).Find(&episodes).Error
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
	episode.Videos = videos
	lengthVideo := len(videos)
	episode.VideoCount = lengthVideo
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
			playlist.Image = "https://www.dailymotion.com/thumbnail/video/" + videoID
			source.File = "https://www.dailymotion.com/embed/video/" + videoID
		case 14:
			SetMThaiThumbnail(episode, &playlist, videoID)
			if embedVideo := GetEmbedVideo(db, videoID); embedVideo.EmbedURL != "" {
				source.File = embedVideo.EmbedURL
				episode.IsURL = false
			} else {
				source.File = "http://video.mthai.com/cool/player/" + videoID + ".html"
				episode.IsURL = true
			}
		case 13, 15:
			SetMThaiThumbnail(episode, &playlist, videoID)
			playlist.Password = episode.Password
			source.File = "http://video.mthai.com/cool/player/" + videoID + ".html"
			episode.IsURL = true
		default:
			episode.IsURL = true
			playlist.Image = "http://thumbnail.instardara.com/chrome.jpg"
			episode.Thumbnail = "http://thumbnail.instardara.com/chrome.jpg"
			source.File = videoID
		}
		playlist.VideoID = videoID
		playlist.File = source.File
		playlist.Sources = append(playlist.Sources, source)
		episode.Playlists = append(episode.Playlists, playlist)
	}
	return
}

func SetMThaiThumbnail(episode *Episode, playlist *Playlist, videoID string) {
	if episode.ID > 674351 {
		playlist.Image = episode.Thumbnail
	} else {
		playlist.Image = "http://video.mthai.com/thumbnail/" + videoID + ".jpg"
	}
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

func GetCountEpisodeByVideoID(db *gorm.DB, videoID string) (count int, err error) {
	err = db.Model(Episode{}).Unscoped().Where("video LIKE ?", "%"+videoID+"%").Count(&count).Error
	return
}
