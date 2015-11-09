package v1

import (
	"github.com/code-mobi/tvthailand.me/Godeps/_workspace/src/github.com/gin-gonic/gin"
	"github.com/code-mobi/tvthailand.me/Godeps/_workspace/src/github.com/jinzhu/gorm"
	"github.com/code-mobi/tvthailand.me/data"
	"net/http"
	"strconv"
)

func CategoriesHandler(c *gin.Context) {
	db := c.MustGet("DB").(gorm.DB)
	categories, err := data.GetCategories(&db)
	if err != nil {
		c.JSON(http.StatusNotFound, err)
	}
	c.JSON(200, map[string]interface{}{
		"Categories": categories,
	})
}

func RecentlyHandler(c *gin.Context) {
	db := c.MustGet("DB").(gorm.DB)
	start, _ := strconv.Atoi(c.Param("start"))
	shows, err := data.GetShowByRecently(&db, start)
	if err != nil {
		c.JSON(http.StatusNotFound, err)
	}
	c.JSON(200, map[string]interface{}{
		"Shows": shows,
	})
}

func PopularHandler(c *gin.Context) {
	db := c.MustGet("DB").(gorm.DB)
	start, _ := strconv.Atoi(c.Param("start"))
	shows, err := data.GetShowByPopular(&db, start)
	if err != nil {
		c.JSON(http.StatusNotFound, err)
	}
	c.JSON(200, map[string]interface{}{
		"Shows": shows,
	})
}

func CategoryHandler(c *gin.Context) {
	db := c.MustGet("DB").(gorm.DB)
	start, _ := strconv.Atoi(c.Param("start"))
	shows, err := data.GetShowByCategory(&db, c.Param("id"), start)
	if err != nil {
		c.JSON(http.StatusNotFound, err)
	}
	c.JSON(200, map[string]interface{}{
		"Shows": shows,
	})
}

func ChannelsHandler(c *gin.Context) {
	db := c.MustGet("DB").(gorm.DB)
	categories, err := data.GetChannels(&db)
	if err != nil {
		c.JSON(http.StatusNotFound, err)
	}
	c.JSON(200, map[string]interface{}{
		"Channels": categories,
	})
}

func ChannelHandler(c *gin.Context) {
	db := c.MustGet("DB").(gorm.DB)
	start, _ := strconv.Atoi(c.Param("start"))
	shows, err := data.GetShowByChannel(&db, c.Param("id"), start)
	if err != nil {
		c.JSON(http.StatusNotFound, err)
	}
	c.JSON(200, map[string]interface{}{
		"Shows": shows,
	})
}

func ShowHandler(c *gin.Context) {
	db := c.MustGet("DB").(gorm.DB)
	start, _ := strconv.Atoi(c.Param("start"))
	showID, _ := strconv.Atoi(c.Param("show_id"))
	show, err := data.GetShow(&db, showID)
	if err != nil {
		c.JSON(http.StatusNotFound, err)
	}
	episodes, err := data.GetEpisodes(&db, show.ID, start)
	if err != nil {
		c.JSON(http.StatusNotFound, err)
	}
	c.JSON(200, map[string]interface{}{
		"Show":     show,
		"Episodes": episodes,
	})
}

func WatchHandler(c *gin.Context) {
	db := c.MustGet("DB").(gorm.DB)
	episode, err := data.GetVideoList(&db, c.Param("hashID"))
	if err != nil {
		c.JSON(http.StatusNotFound, episode)
	}
	show, _ := data.GetShow(&db, episode.ShowID)
	c.JSON(200, map[string]interface{}{
		"Show":    show,
		"Episode": episode,
	})
}
