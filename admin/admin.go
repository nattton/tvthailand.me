package admin

import (
	"github.com/code-mobi/tvthailand.me/Godeps/_workspace/src/github.com/go-martini/martini"
	"github.com/code-mobi/tvthailand.me/Godeps/_workspace/src/github.com/jinzhu/gorm"
	"github.com/code-mobi/tvthailand.me/Godeps/_workspace/src/github.com/martini-contrib/render"
	"github.com/code-mobi/tvthailand.me/analytic"
	"github.com/code-mobi/tvthailand.me/data"
	"io/ioutil"
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

func AnalyticProcessHandler(db gorm.DB, r render.Render, req *http.Request) {
	req.ParseMultipartForm(1024)
	fnPleaseChooseFile := func() {
		flash := map[string]string{"warning": "Please Choose Anlytice.json File"}
		r.HTML(http.StatusOK, "admin/index", map[string]interface{}{
			"flash": flash,
		})
	}

	if len(req.MultipartForm.File["uploaded"]) == 0 {
		fnPleaseChooseFile()
		return
	}

	fileHeader := req.MultipartForm.File["uploaded"][0]
	file, err := fileHeader.Open()
	if err != nil {
		fnPleaseChooseFile()
		return
	}
	data, err := ioutil.ReadAll(file)
	if err != nil {
		fnPleaseChooseFile()
		return
	}
	analytic.UpdateView(&db, data)
	flash := map[string]string{"info": "Update View Count Successfully"}
	r.HTML(http.StatusOK, "admin/index", map[string]interface{}{
		"flash": flash,
	})
}
