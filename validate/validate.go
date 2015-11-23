package validate

import (
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/code-mobi/tvthailand.me/Godeps/_workspace/src/github.com/facebookgo/httpcontrol"
	"github.com/code-mobi/tvthailand.me/Godeps/_workspace/src/github.com/jinzhu/gorm"
)
// URL check URL isValid
func URL(url string) (bool, error) {
	client := &http.Client{
		Transport: &httpcontrol.Transport{
			RequestTimeout: time.Minute,
			MaxTries:       3,
		},
	}
	resp, err := client.Get(url)
	if err != nil {
		log.Println(err)
		return false, err
	}
	if resp.StatusCode == 404 {
		return false, nil
	}
	return true, nil
}

// RunWebURL Validate Web URL 
func RunWebURL(db *gorm.DB, start int, limit int, delay time.Duration) {
	episodes, _ := EpisodeWebURL(db, start, limit)
	for _, episode := range episodes {
		videos := strings.Split(episode.Video, ",")
		for _, video := range videos {
			time.Sleep(delay)
			isValid, err := URL(video)
			if err != nil {
				log.Println(episode.ID, "URL Err : ", err)
			} else if !isValid {
				log.Println("### ", episode.ID, video, "is not valid : ")
				db.Model(&episode).UpdateColumn("banned", 1)
				break
			} else if isValid {
				log.Println(episode.ID, video, "is valid")
			}
		}
	}
}
