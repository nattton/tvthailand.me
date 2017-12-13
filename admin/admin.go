package admin

import (
	"io/ioutil"
	"net/http"
	"strconv"

	"github.com/code-mobi/tvthailand.me/analytic"
	"github.com/code-mobi/tvthailand.me/data"
	"github.com/code-mobi/tvthailand.me/utils"
	"github.com/gin-gonic/gin"
)

func printFlash(writer http.ResponseWriter, flashType, message string) {
	flash := map[string]string{
		flashType: message,
	}
	renderData := map[string]interface{}{
		"flash": flash,
	}
	utils.GenerateHTML(writer, renderData, "admin/layout", "admin/index")
}

func IndexHandler(c *gin.Context) {
	utils.GenerateHTML(c.Writer, nil, "admin/layout", "admin/index")
}

func AddEmbedMThaiHandler(c *gin.Context) {
	db, _ := utils.OpenDB()
	defer db.Close()
	showID, _ := strconv.Atoi(c.Request.FormValue("show_id"))
	data.InsertMThaiEmbedVideos(db, showID)
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
		flash := map[string]string{"warning": "Please Choose Anlytice CSV File"}
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
	shows := analytic.UpdateView(db, dataFile)
	flash := map[string]string{
		"info": "Update View Count Successfully",
	}
	renderData := map[string]interface{}{
		"flash": flash,
		"shows": shows,
	}

	utils.GenerateHTML(c.Writer, renderData, "admin/layout", "admin/analytic")
}

func FlushHandler(c *gin.Context) {
	redisClient := utils.OpenRedis()
	redisClient.FlushAll()
	flash := map[string]string{
		"info": "Flush All Successfully",
	}
	renderData := map[string]interface{}{
		"flash": flash,
	}
	utils.GenerateHTML(c.Writer, renderData, "admin/layout", "admin/index")
}
