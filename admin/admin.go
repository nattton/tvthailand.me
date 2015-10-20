package admin

import (
	"github.com/code-mobi/tvthailand.me/Godeps/_workspace/src/github.com/go-martini/martini"
	"github.com/code-mobi/tvthailand.me/Godeps/_workspace/src/github.com/jinzhu/gorm"
	"github.com/code-mobi/tvthailand.me/Godeps/_workspace/src/github.com/martini-contrib/render"
	"github.com/code-mobi/tvthailand.me/data"
	"net/http"
	"strconv"
)

func EncryptHandler(db gorm.DB, r render.Render, params martini.Params) {
	episodeID, _ := strconv.Atoi(params["episodeID"])
	data.EncryptEpisode(&db, episodeID)
	r.HTML(http.StatusOK, "index", map[string]interface{}{
		"header": "Encrypt Successfully",
	})
}

func GetEmbedMThaiHandler(db gorm.DB, r render.Render, params martini.Params) {
	showID, _ := strconv.Atoi(params["showID"])
	data.InsertMThaiEmbedVideos(&db, showID)
	r.HTML(http.StatusOK, "show/list", map[string]interface{}{
		"Title": "Embed",
	})
}
