package main

import (
	"net/http"
	"strconv"

	"github.com/code-mobi/tvthailand.me/data"
	"github.com/code-mobi/tvthailand.me/utils"
	"github.com/gin-gonic/gin"
)

func AjaxRecentlyHandler(c *gin.Context) {
	db, _ := utils.OpenDB()
	defer db.Close()
	offset, _ := strconv.Atoi(c.Query("offset"))
	shows, err := data.ShowsRecently(db, offset)
	if err != nil {
		c.JSON(http.StatusNotFound, err)
	}
	c.JSON(http.StatusOK, map[string]interface{}{
		"shows": shows,
	})
}

func AjaxPopularHandler(c *gin.Context) {
	db, _ := utils.OpenDB()
	defer db.Close()
	offset, _ := strconv.Atoi(c.Query("offset"))
	shows, err := data.ShowsPopular(db, offset)
	if err != nil {
		c.JSON(http.StatusNotFound, err)
	}
	c.JSON(http.StatusOK, map[string]interface{}{
		"shows": shows,
	})
}

func AjaxCategoryHandler(c *gin.Context) {
	db, _ := utils.OpenDB()
	defer db.Close()
	offset, _ := strconv.Atoi(c.Query("offset"))
	id, _ := strconv.Atoi(c.Param("id"))
	shows, err := data.ShowsCategory(db, id, offset)
	if err != nil {
		c.JSON(http.StatusNotFound, err)
	}
	c.JSON(http.StatusOK, map[string]interface{}{
		"shows": shows,
	})
}

func AjaxChannelsHandler(c *gin.Context) {
	db, _ := utils.OpenDB()
	defer db.Close()
	categories, err := data.ChannelsActive(db)
	if err != nil {
		c.JSON(http.StatusNotFound, err)
	}
	c.JSON(http.StatusOK, map[string]interface{}{
		"channels": categories,
	})
}

func AjaxChannelHandler(c *gin.Context) {
	db, _ := utils.OpenDB()
	defer db.Close()
	offset, _ := strconv.Atoi(c.Query("offset"))
	id, _ := strconv.Atoi(c.Param("id"))
	shows, err := data.ShowsChannel(db, id, offset)
	if err != nil {
		c.JSON(http.StatusNotFound, err)
	}
	c.JSON(http.StatusOK, map[string]interface{}{
		"shows": shows,
	})
}

func AjaxShowHandler(c *gin.Context) {
	db, _ := utils.OpenDB()
	defer db.Close()
	page, _ := strconv.Atoi(c.Query("page"))
	showID, _ := strconv.Atoi(c.Param("show_id"))
	episodes, pageInfo, err := data.EpisodesAndPageInfo(db, showID, int32(page))
	if err != nil {
		c.JSON(http.StatusNotFound, err)
	}
	c.JSON(200, map[string]interface{}{
		"pageInfo": pageInfo,
		"episodes": episodes,
	})
}
