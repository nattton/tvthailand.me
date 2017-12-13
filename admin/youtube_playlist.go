package admin

import (
	"log"

	"github.com/code-mobi/tvthailand.me/data"
	"github.com/code-mobi/tvthailand.me/utils"
	"github.com/code-mobi/tvthailand.me/youtube"
	"github.com/gin-gonic/gin"
)

type playlistSearchForm struct {
	PlaylistID string `form:"playlistId"`
	MaxResults int    `form:"maxResults"`
	PageToken  string `form:"pageToken"`
}

func YoutubePlaylistHandler(c *gin.Context) {
	db, _ := utils.OpenDB()
	defer db.Close()

	var form playlistSearchForm
	c.Bind(&form)

	playlistOptions := data.YoutubePlaylistOptions(db, form.PlaylistID)
	if form.MaxResults == 0 {
		form.MaxResults = 40
	}
	renderData := map[string]interface{}{
		"form":            form,
		"playlistOptions": playlistOptions,
	}
	log.Println(form)
	if form.PlaylistID != "" {
		y := youtube.NewYoutube()
		api, _ := y.GetExVideoByPlaylistID(form.PlaylistID, form.MaxResults, form.PageToken)
		form.PageToken = api.NextPageToken
		renderData["form"] = form
		renderData["items"] = api.Items
	}

	utils.GenerateHTML(c.Writer, renderData, "admin/layout", "admin/youtube_playlist")
}
