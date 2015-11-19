package admin

import (
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"

	"github.com/code-mobi/tvthailand.me/Godeps/_workspace/src/github.com/gin-gonic/gin"
	"github.com/code-mobi/tvthailand.me/data"
	"github.com/code-mobi/tvthailand.me/utils"
)

const (
	dateFMT      = "2006-01-02"
	dateSmallFMT = "20060102"
)

type episodeForm struct {
	ID        int    `form:"id"`
	ShowID    int    `form:"show_id" binding:"required"`
	Ep        int    `form:"ep"`
	Title     string `form:"title"`
	Video     string `form:"video" binding:"required"`
	PartTotal int    `form:"part_total"`
	SrcType   int    `form:"src_type"`
	Date      string `form:"date"`
	Password  string `form:"password"`
	User      string `form:"user"`
}

// GetEpisodeHandler GET /admin/episode
func GetEpisodeHandler(c *gin.Context) {
	utils.GenerateHTML(c.Writer, nil, "admin/layout", "admin/index")
}

// EncryptEpisodeHandler GET /encrypt_episode/:episodeID
func EncryptEpisodeHandler(c *gin.Context) {
	db, _ := utils.OpenDB()
	defer db.Close()
	episodeID, _ := strconv.Atoi(c.Param("episodeID"))
	if episodeID > 0 {
		data.EncryptAllEpisodes(&db)
	} else {
		episode, err := data.GetEpisode(&db, episodeID)
		if err != nil {
			printFlash(c.Writer, "danger", "Not Found Episode")
		}
		data.EncryptEpisode(&db, &episode)
	}
	flash := map[string]string{
		"info": "Encrypt Successfully",
	}
	renderData := map[string]interface{}{
		"flash": flash,
	}
	utils.GenerateHTML(c.Writer, renderData, "admin/layout", "admin/index")
}

// SaveEpisodeHandler POST /admin/episode
func SaveEpisodeHandler(c *gin.Context) {
	var form episodeForm
	err := c.Bind(&form)
	if err != nil {
		printFlash(c.Writer, "danger", err.Error())
		return
	}

	db, _ := utils.OpenDB()
	defer db.Close()
	log.Println(form)
	var episode data.Episode
	if form.ID > 0 {
		episode, _ = data.GetEpisode(&db, form.ID)
	}

	episode.Title = form.Title
	// Set Video
	episode.SetVideo(form.Video)
	// Set Parts
	if form.PartTotal > 0 {
		lengthVideo := len(strings.Split(episode.Video, ","))
		if lengthVideo < form.PartTotal {
			episode.Parts = fmt.Sprintf("%d/%d", lengthVideo, form.PartTotal)
			if episode.Title == "" {
				episode.Title = fmt.Sprintf("(%s)", episode.Parts)
			} else {
				episode.Title = fmt.Sprintf("%s - (%s)", episode.Title, episode.Parts)
			}
		}
	}
	// Set Date
	episode.Date, err = time.Parse(dateFMT, form.Date)
	if err != nil {
		printFlash(c.Writer, "danger", err.Error())
		return
	}
	// Set Ep
	if form.Ep > 0 {
		episode.Ep = form.Ep
	} else {
		episode.Ep, _ = strconv.Atoi(episode.Date.Format(dateSmallFMT))
	}

	episode.ShowID = form.ShowID
	episode.SrcType = form.SrcType
	episode.Password = form.Password
	episode.User = form.User
	db.Save(&episode)
	data.EncryptEpisode(&db, &episode)
	// Set Show UpdateDate
	data.ShowUpdateDate(&db, episode.ShowID)
	// Set Status BotVideo to Updated
	data.SetBotVideoUpdated(&db, episode.Video)
	
	DeleteCached(episode.ShowID)

	message := fmt.Sprintf(`Save Episode <a href="/watch/%d/0">%d : %s</a> Successful`, episode.ID, episode.ID, episode.Title)
	printFlash(c.Writer, "info", message)
}

// DeleteCached by ShowID
func DeleteCached(showID int) {
	utils.DeleteListCached("API_RECENTLY")
	utils.DeleteListCached(fmt.Sprintf("API_SHOW:%d", showID))
}
