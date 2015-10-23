package admin

import (
	"github.com/code-mobi/tvthailand.me/Godeps/_workspace/src/github.com/go-martini/martini"
	"github.com/code-mobi/tvthailand.me/Godeps/_workspace/src/github.com/jinzhu/gorm"
	"github.com/code-mobi/tvthailand.me/Godeps/_workspace/src/github.com/martini-contrib/render"
	"github.com/code-mobi/tvthailand.me/data"
	"net/http"
	"strconv"
)

func EncryptEpisodeHandler(db gorm.DB, r render.Render, params martini.Params) {
	episodeID, _ := strconv.Atoi(params["episodeID"])
	data.EncryptEpisode(&db, episodeID)
	flash := map[string]string{"info": "Encrypt Successfully"}
	r.HTML(http.StatusOK, "admin/index", map[string]interface{}{
		"flash": flash,
	})
}

func AddEmbedMThaiHandler(db gorm.DB, r render.Render, req *http.Request) {
	showID, _ := strconv.Atoi(req.FormValue("show_id"))
	data.InsertMThaiEmbedVideos(&db, showID)
	flash := map[string]string{"info": "Insert MThai Embed Videos Successfully"}
	r.HTML(http.StatusOK, "admin/index", map[string]interface{}{
		"flash": flash,
	})
}

func IndexHandler(r render.Render) {
	r.HTML(http.StatusOK, "admin/index", nil)
}
