package admin

import (
	"github.com/code-mobi/tvthailand.me/Godeps/_workspace/src/github.com/gin-gonic/gin"
	"github.com/code-mobi/tvthailand.me/analytic"
	"github.com/code-mobi/tvthailand.me/data"
	"github.com/code-mobi/tvthailand.me/utils"
	"io/ioutil"
	"strconv"
)

func IndexHandler(c *gin.Context) {
	utils.GenerateHTML(c.Writer, nil, "admin/layout", "admin/index")
}

func EncryptEpisodeHandler(c *gin.Context) {
	db, _ := utils.OpenDB()
	defer db.Close()
	episodeID, _ := strconv.Atoi(c.Param("episodeID"))
	data.EncryptEpisode(&db, episodeID)
	flash := map[string]string{
		"info": "Encrypt Successfully",
	}
	renderData := map[string]interface{}{
		"flash": flash,
	}
	utils.GenerateHTML(c.Writer, renderData, "admin/layout", "index")
}

func AddEmbedMThaiHandler(c *gin.Context) {
	db, _ := utils.OpenDB()
	defer db.Close()
	showID, _ := strconv.Atoi(c.Request.FormValue("show_id"))
	data.InsertMThaiEmbedVideos(&db, showID)
	flash := map[string]string{"info": "Insert MThai Embed Videos Successfully"}
	utils.GenerateHTML(c.Writer, map[string]interface{}{
		"flash": flash,
	},
		"admin/layout", "admin/index")
}

func AnalyticHandler(c *gin.Context) {
	utils.GenerateHTML(c.Writer, nil, "admin/layout", "admin/analytic")
}

func AnalyticProcessHandler(c *gin.Context) {
	db, _ := utils.OpenDB()
	defer db.Close()
	c.Request.ParseMultipartForm(1024)
	fnPleaseChooseFile := func() {
		flash := map[string]string{"warning": "Please Choose Anlytice.json File"}
		utils.GenerateHTML(c.Writer, map[string]interface{}{
			"flash": flash,
		},
			"admin/layout", "admin/analytic")
	}

	if len(c.Request.MultipartForm.File["uploaded"]) == 0 {
		fnPleaseChooseFile()
		return
	}

	fileHeader := c.Request.MultipartForm.File["uploaded"][0]
	file, err := fileHeader.Open()
	if err != nil {
		fnPleaseChooseFile()
		return
	}
	dataFile, err := ioutil.ReadAll(file)
	if err != nil {
		fnPleaseChooseFile()
		return
	}
	shows := analytic.UpdateView(&db, dataFile)
	flash := map[string]string{
		"info": "Update View Count Successfully",
	}
	renderData := map[string]interface{}{
		"flash": flash,
		"shows": shows,
	}

	utils.GenerateHTML(c.Writer, renderData, "admin/layout", "admin/analytic")
}
