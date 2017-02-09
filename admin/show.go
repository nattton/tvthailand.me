package admin

import (
	"strconv"
	"strings"
	"time"

	"github.com/code-mobi/tvthailand.me/Godeps/_workspace/src/github.com/gin-gonic/gin"
	"github.com/code-mobi/tvthailand.me/data"
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
	CoverPhoto  string
	Detail      string
	UpdateDate  time.Time
	OtvID       string
	OtvLogo     string
	OtvAPIName  string
	IsOtv       bool
	IsActive    bool
	IsOnline    bool
	Ios         bool
	Android     bool
	Wp          bool
	S40         bool
	ThRestrict  bool
	V3          bool
	Owner       string

	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt *time.Time
}

type showForm struct {
	ID          int    `gorm:"primary_key" form:"id"`
	CategoryID  int    `form:"category_id" binding:"required"`
	ChannelID   int    `form:"channel_id" binding:"required"`
	Title       string `form:"title" binding:"required"`
	Description string `form:"description" binding:"required"`
	Thumbnail   string `form:"thumbnail"`
	CoverPhoto  string `form:"cover_photo"`
	Poster      string `form:"poster"`
	Detail      string `form:"detail"`
	UpdateDate  string `form:"update_date"`
	OtvID       string `form:"otv_id"`
	OtvLogo     string `form:"otv_logo"`
	IsOtv       bool   `form:"is_otv"`
	IsActive    bool   `form:"is_active"`
	IsOnline    bool   `form:"is_online"`
	Ios         bool   `form:"ios"`
	Android     bool   `form:"android"`
	Wp          bool   `form:"wp"`
	S40         bool   `form:"s40"`
	ThRestrict  bool   `form:"th_restrict"`
	V3          bool   `form:"v3"`
	Owner       string `form:"owner"`
}

// ShowsHandler display list shows
// GET /admin/shows
// GET /admin/shows/search?q=?
func ShowsHandler(c *gin.Context) {
	db, _ := utils.OpenDB()
	defer db.Close()
	var shows []Show
	q := strings.TrimSpace(c.Query("q"))
	if q != "" {
		db.Where("id = ? OR title LIKE ?", q, "%"+q+"%").Order("title asc").Limit(20).Find(&shows)
	} else {
		db.Order("id desc").Limit(20).Find(&shows)
	}

	renderData := map[string]interface{}{
		"q":     q,
		"shows": shows,
	}
	flash, exist := c.Get("flash")
	if exist {
		renderData["flash"] = flash
	}
	utils.GenerateHTML(c.Writer, renderData, "admin/layout", "admin/show_list")
}

func ShowNewHandler(c *gin.Context) {
	db, _ := utils.OpenDB()
	defer db.Close()

	var form showForm
	c.Bind(&form)
	// Form Default Value
	if form.Title == "" {
		form.IsActive = true
		form.IsOnline = true
		form.Ios = true
		form.Android = true
		form.Wp = true
		form.S40 = true
		form.V3 = true
	}
	renderData := map[string]interface{}{
		"show":            form,
		"categoryOptions": data.CategoryOptions(&db, form.CategoryID),
		"channelOptions":  data.ChannelOptions(&db, form.ChannelID),
	}
	flash, exist := c.Get("flash")
	if exist {
		renderData["flash"] = flash
	}
	utils.GenerateHTML(c.Writer, renderData, "admin/layout", "admin/show_edit")
}

func ShowEditHandler(c *gin.Context) {
	db, _ := utils.OpenDB()
	defer db.Close()
	var show Show
	db.First(&show, c.Param("id"))
	renderData := map[string]interface{}{
		"show":            show,
		"categoryOptions": data.CategoryOptions(&db, show.CategoryID),
		"channelOptions":  data.ChannelOptions(&db, show.ChannelID),
	}
	flash, exist := c.Get("flash")
	if exist {
		renderData["flash"] = flash
	}
	utils.GenerateHTML(c.Writer, renderData, "admin/layout", "admin/show_edit")
}

func ShowUpdateHandler(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))

	var form showForm
	err := c.Bind(&form)
	if err != nil {
		showSaveError(c, id, err)
		return
	}

	db, _ := utils.OpenDB()
	defer db.Close()

	var show Show
	if id > 0 {
		err := db.First(&show, id).Error
		if err != nil {
			showSaveError(c, show.ID, err)
			return
		}
	}

	show.CategoryID = form.CategoryID
	show.ChannelID = form.ChannelID
	show.Title = form.Title
	show.Description = form.Description
	show.Thumbnail = form.Thumbnail
	show.CoverPhoto = form.CoverPhoto
	show.Poster = form.Poster
	show.Detail = form.Detail
	if form.UpdateDate != "" {
		updateDate, err := time.Parse("2006-01-02 15:04:05", form.UpdateDate)
		if err == nil {
			show.UpdateDate = updateDate
		}
	}

	show.OtvID = form.OtvID
	show.OtvLogo = form.OtvLogo
	show.IsOtv = form.IsOtv

	show.IsActive = form.IsActive
	show.IsOnline = form.IsOnline
	show.Ios = form.Ios
	show.Android = form.Android
	show.Wp = form.Wp
	show.S40 = form.S40
	show.ThRestrict = form.ThRestrict
	show.V3 = form.V3
	show.Owner = form.Owner

	err = db.Save(&show).Error
	if err != nil {
		showSaveError(c, id, err)
		return
	}

	flash := map[string]string{"info": "Save Successful"}
	renderData := map[string]interface{}{
		"flash":           flash,
		"show":            show,
		"categoryOptions": data.CategoryOptions(&db, show.CategoryID),
		"channelOptions":  data.ChannelOptions(&db, show.ChannelID),
	}
	utils.GenerateHTML(c.Writer, renderData, "admin/layout", "admin/show_edit")
}

func ShowToggleActivateHandler(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))

	db, _ := utils.OpenDB()
	var show Show
	err := db.First(&show, id).Error
	if err != nil {
		showSaveError(c, show.ID, err)
		return
	}

	show.IsActive = !show.IsActive
	db.Save(&show)
	flash := map[string]string{"info": "Update Successful"}
	c.Set("flash", flash)
	ShowsHandler(c)
}

func showSaveError(c *gin.Context, id int, err error) {
	flash := map[string]string{"danger": strings.Replace(err.Error(), "\n", "<br/>", -1)}
	c.Set("flash", flash)
	if id > 0 {
		ShowEditHandler(c)
	} else {
		ShowNewHandler(c)
	}
}
