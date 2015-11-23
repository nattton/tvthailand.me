package admin

import (
	"time"

	"github.com/code-mobi/tvthailand.me/Godeps/_workspace/src/github.com/gin-gonic/gin"
	"github.com/code-mobi/tvthailand.me/utils"
)

type Show struct {
	ID          int `gorm:"primary_key"`
	CategoryID  int
	ChannelID   int
	Title       string
	Description string
	Thumbnail   string
	Poster      string
	Detail      string
	LastEpname  string
	ViewCount   int
	Rating      float32
	VoteCount   int
	IsOtv       bool
	OtvID       string
	UpdateDate  time.Time

	CreatedAt time.Time
	UpdatedAt time.Time
	// DeletedAt *time.Time
}

func ShowsHandler(c *gin.Context) {
	db, _ := utils.OpenDB()
	var shows []Show
	db.Order("id desc").Limit(20).Find(&shows)
	renderData := map[string]interface{}{
		"shows": shows,
	}
	utils.GenerateHTML(c.Writer, renderData, "admin/layout", "admin/show")
}


func SearchShowsHandler(c *gin.Context) {
	db, _ := utils.OpenDB()
	q := c.Query("q")
	var shows []Show
	db.Where("title LIKE ?", "%"+q+"%").Order("title asc").Limit(20).Find(&shows)
	renderData := map[string]interface{}{
		"shows": shows,
	}
	utils.GenerateHTML(c.Writer, renderData, "admin/layout", "admin/show")
}

func ShowEditHandler(c *gin.Context) {
	db,_ := utils.OpenDB()
	var show Show
	db.First(&show, c.Param("id"))
	renderData := map[string]interface{}{
		"show": show,
	}
	utils.GenerateHTML(c.Writer, renderData, "admin/layout", "admin/show_edit")
}
