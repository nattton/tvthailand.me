package main

import (
	"github.com/code-mobi/tvthailand.me/Godeps/_workspace/src/github.com/gin-gonic/gin"
	"github.com/code-mobi/tvthailand.me/Godeps/_workspace/src/github.com/jinzhu/gorm"
	"github.com/code-mobi/tvthailand.me/data"
	"net/http"
	"strconv"
)

func AjaxRecentlyHandler(c *gin.Context) {
	db := c.MustGet("DB").(gorm.DB)
	offset, _ := strconv.Atoi(c.Query("offset"))
	shows, err := data.GetShowByRecently(&db, offset)
	if err != nil {
		c.JSON(http.StatusNotFound, err)
	}
	c.JSON(http.StatusOK, map[string]interface{}{
		"shows": shows,
	})
}

func AjaxPopularHandler(c *gin.Context) {
	db := c.MustGet("DB").(gorm.DB)
	offset, _ := strconv.Atoi(c.Query("offset"))
	shows, err := data.GetShowByPopular(&db, offset)
	if err != nil {
		c.JSON(http.StatusNotFound, err)
	}
	c.JSON(http.StatusOK, map[string]interface{}{
		"shows": shows,
	})
}

func AjaxCategoryHandler(c *gin.Context) {
	db := c.MustGet("DB").(gorm.DB)
	offset, _ := strconv.Atoi(c.Query("offset"))
	shows, err := data.GetShowByCategory(&db, c.Param("id"), offset)
	if err != nil {
		c.JSON(http.StatusNotFound, err)
	}
	c.JSON(http.StatusOK, map[string]interface{}{
		"shows": shows,
	})
}

func AjaxChannelsHandler(c *gin.Context) {
	db := c.MustGet("DB").(gorm.DB)
	categories, err := data.GetChannels(&db)
	if err != nil {
		c.JSON(http.StatusNotFound, err)
	}
	c.JSON(http.StatusOK, map[string]interface{}{
		"channels": categories,
	})
}

func AjaxChannelHandler(c *gin.Context) {
	db := c.MustGet("DB").(gorm.DB)
	offset, _ := strconv.Atoi(c.Query("offset"))
	shows, err := data.GetShowByChannel(&db, c.Param("id"), offset)
	if err != nil {
		c.JSON(http.StatusNotFound, err)
	}
	c.JSON(http.StatusOK, map[string]interface{}{
		"shows": shows,
	})
}

func AjaxShowHandler(c *gin.Context) {
	db := c.MustGet("DB").(gorm.DB)
	page, _ := strconv.Atoi(c.Query("page"))
	showID, _ := strconv.Atoi(c.Param("show_id"))
	episodes, pageInfo, err := data.GetEpisodesAndPageInfo(&db, showID, int32(page))
	if err != nil {
		c.JSON(http.StatusNotFound, err)
	}
	c.JSON(200, map[string]interface{}{
		"pageInfo": pageInfo,
		"episodes": episodes,
	})
}
