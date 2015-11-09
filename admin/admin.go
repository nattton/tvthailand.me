package admin

import (
	"github.com/code-mobi/tvthailand.me/Godeps/_workspace/src/github.com/gin-gonic/gin"
	"github.com/code-mobi/tvthailand.me/Godeps/_workspace/src/github.com/jinzhu/gorm"
	"github.com/code-mobi/tvthailand.me/analytic"
	"github.com/code-mobi/tvthailand.me/data"
	"github.com/code-mobi/tvthailand.me/utils"
	"io/ioutil"
	"strconv"
)

func IndexHandler(c *gin.Context) {
	utils.GenerateHTML(c.Writer, nil, "layout", "mobile_ads", "admin/index")
}

func EncryptEpisodeHandler(c *gin.Context) {
	db := c.MustGet("DB").(gorm.DB)
	episodeID, _ := strconv.Atoi(c.Param("episodeID"))
	data.EncryptEpisode(&db, episodeID)
	flash := map[string]string{
		"info": "Encrypt Successfully",
	}
	data := map[string]interface{}{
		"flash": flash,
	}
	utils.GenerateHTML(c.Writer, data, "layout", "mobile_ads", "admin/index")
}

func AddEmbedMThaiHandler(c *gin.Context) {
	db := c.MustGet("DB").(gorm.DB)
	showID, _ := strconv.Atoi(c.Request.FormValue("show_id"))
	data.InsertMThaiEmbedVideos(&db, showID)
	flash := map[string]string{"info": "Insert MThai Embed Videos Successfully"}
	utils.GenerateHTML(c.Writer, map[string]interface{}{
		"flash": flash,
	},
		"layout", "mobile_ads", "admin/index")
}

func AnalyticProcessHandler(c *gin.Context) {
	db := c.MustGet("DB").(gorm.DB)
	c.Request.ParseMultipartForm(1024)
	fnPleaseChooseFile := func() {
		flash := map[string]string{"warning": "Please Choose Anlytice.json File"}
		utils.GenerateHTML(c.Writer, map[string]interface{}{
			"flash": flash,
		},
			"layout", "mobile_ads", "admin/index")
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
	data, err := ioutil.ReadAll(file)
	if err != nil {
		fnPleaseChooseFile()
		return
	}
	analytic.UpdateView(&db, data)
	flash := map[string]string{
		"info": "Update View Count Successfully",
	}
	utils.GenerateHTML(c.Writer, map[string]interface{}{
		"flash": flash,
	},
		"layout", "mobile_ads", "admin/index")
}
