package data

import (
	"fmt"
	"html"
	"log"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/jinzhu/gorm"
)

const (
	MThaiURL = "http://video.mthai.com/cool/player/%s.html"
)

type EmbedVideo struct {
	ID       int
	VideoID  string
	EmbedURL string
}

func GetEmbedVideo(db *gorm.DB, videoID string) (embedVideo EmbedVideo) {
	err := db.Where("video_id = ?", videoID).First(&embedVideo).Error
	if err != nil {
		embedURL, _ := GetMThaiEmbedURL(videoID)
		embedVideo = EmbedVideo{VideoID: videoID, EmbedURL: embedURL}
		db.Create(embedVideo)
	}
	return
}

func InsertMThaiEmbedVideos(db *gorm.DB, showID int) {
	var episodes []Episode
	db.Where("src_type = ? AND show_id = ?", 14, showID).Find(&episodes)
	for _, episode := range episodes {
		videos := strings.Split(episode.Video, ",")
		for _, v := range videos {
			embedVideo := EmbedVideo{}
			err := db.Where("video_id = ?", v).First(&embedVideo).Error
			if err != nil {
				embedURL, _ := GetMThaiEmbedURL(v)
				embedVideo = EmbedVideo{VideoID: v, EmbedURL: embedURL}
				db.Create(embedVideo)
			}
		}
	}
}

func GetMThaiEmbedURL(id string) (iframeURL string, thumbnailURL string) {
	url := fmt.Sprintf(MThaiURL, id)
	doc, err := goquery.NewDocument(url)
	if err != nil {
		log.Fatal(err)
	}
	var iframe string
	doc.Find(".input-copy-text").Each(func(i int, s *goquery.Selection) {
		val, _ := s.Attr("value")
		if strings.Contains(val, "iframe") && iframe == "" {
			iframe = val
		}
	})
	iframe = html.UnescapeString(iframe)
	reader := strings.NewReader(iframe)
	docIframe, err := goquery.NewDocumentFromReader(reader)
	docIframe.Find("iframe").Each(func(i int, s *goquery.Selection) {
		val, _ := s.Attr("src")
		if val != "" {
			iframeURL = val
		}
	})
	doc.Find("link[itemprop=thumbnailUrl]").Each(func(i int, s *goquery.Selection) {
		val, _ := s.Attr("href")
		thumbnailURL = val
		log.Println("thumbnailURL:", thumbnailURL)
	})
	return
}
