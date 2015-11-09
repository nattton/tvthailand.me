package main

import (
	"github.com/code-mobi/tvthailand.me/Godeps/_workspace/src/github.com/gin-gonic/gin"
	"github.com/code-mobi/tvthailand.me/data"
	"github.com/code-mobi/tvthailand.me/utils"
	"html"
	"net/http"
	"os"
	"strconv"
	"strings"
)

func indexHandler(c *gin.Context) {
	db, _ := utils.OpenDB()
	defer db.Close()
	recents, _ := data.GetShowByRecently(&db, 0)
	populars, _ := data.GetShowByPopular(&db, 0)
	data := map[string]interface{}{
		"showRecents":  recents,
		"showPopulars": populars,
		"isMobile":     utils.IsMobile(c.Request.UserAgent()),
	}

	utils.GenerateHTML(c.Writer, data, "layout", "mobile_ads", "index")
}

func notFoundHandler(c *gin.Context) {
	utils.GenerateHTML(c.Writer, nil, "layout", "mobile_ads", "not_found")
}

func goOutHandler(c *gin.Context) {
	c.Redirect(http.StatusMovedPermanently, "/not_found")
}

func recentlyHandler(c *gin.Context) {
	db, _ := utils.OpenDB()
	defer db.Close()
	shows, _ := data.GetShowByRecently(&db, 0)
	data := map[string]interface{}{
		"Title":    "รายการล่าสุด",
		"header":   "รายการล่าสุด",
		"typeMode": "recently",
		"shows":    shows,
		"isMobile": utils.IsMobile(c.Request.UserAgent()),
	}
	utils.GenerateHTML(c.Writer, data, "layout", "mobile_ads", "show/list", "episode/item")
}

func popularHandler(c *gin.Context) {
	db, _ := utils.OpenDB()
	defer db.Close()
	shows, _ := data.GetShowByPopular(&db, 0)
	data := map[string]interface{}{
		"Title":    "Popular",
		"header":   "Popular",
		"typeMode": "popular",
		"shows":    shows,
		"isMobile": utils.IsMobile(c.Request.UserAgent()),
	}
	utils.GenerateHTML(c.Writer, data, "layout", "mobile_ads", "show/list", "episode/item")
}

func categoriesHandler(c *gin.Context) {
	db, _ := utils.OpenDB()
	defer db.Close()
	categories, _ := data.GetCategories(&db)
	data := map[string]interface{}{
		"header":     "หมวด",
		"categories": categories,
		"isMobile":   utils.IsMobile(c.Request.UserAgent()),
	}
	utils.GenerateHTML(c.Writer, data, "layout", "mobile_ads", "category/list")
}

func categoryShowHandler(c *gin.Context) {
	db, _ := utils.OpenDB()
	defer db.Close()
	titlize := c.Param("titlize")
	category, _ := data.GetCategory(&db, titlize)
	shows, _ := data.GetShowByCategory(&db, category.ID, 0)
	data := map[string]interface{}{
		"Title":    category.Title,
		"header":   category.Title,
		"typeMode": "category",
		"typeId":   category.ID,
		"shows":    shows,
	}
	utils.GenerateHTML(c.Writer, data, "layout", "mobile_ads", "show/list", "episode/item")
}

func channelsHandler(c *gin.Context) {
	db, _ := utils.OpenDB()
	defer db.Close()
	channels, _ := data.GetChannels(&db)
	data := map[string]interface{}{
		"header":   "ช่องทีวี",
		"channels": channels,
	}
	utils.GenerateHTML(c.Writer, data, "layout", "mobile_ads", "channel/list")
}

func channelShowHandler(c *gin.Context) {
	db, _ := utils.OpenDB()
	defer db.Close()
	id := c.Param("id")
	channel, _ := data.GetChannel(&db, id)
	shows, _ := data.GetShowByChannel(&db, channel.ID, 0)
	data := map[string]interface{}{
		"Title":    channel.Title,
		"header":   channel.Title,
		"channel":  channel,
		"typeMode": "category",
		"typeId":   channel.ID,
		"shows":    shows,
	}
	utils.GenerateHTML(c.Writer, data, "layout", "mobile_ads", "show/list", "episode/item")
}

func searchShowHandler(c *gin.Context) {
	db, _ := utils.OpenDB()
	defer db.Close()
	qs := c.Request.URL.Query()
	keyword := qs.Get("keyword")
	var shows []data.Show
	var episodes []data.Episode
	var header string
	var title string
	if keyword != "" {
		shows, _ = data.GetShowBySearch(&db, keyword)
		episodes, _ = data.GetEpisodesBySearch(&db, keyword)
		header = "ผลการค้นหา : " + keyword
		title = header
	} else {
		title = "Search"
		header = "กรุณาพิมพชื่อเรื่องที่ต้องการค้นหา"
	}
	data := map[string]interface{}{
		"Title":    title,
		"keyword":  keyword,
		"header":   header,
		"shows":    shows,
		"episodes": episodes,
	}
	utils.GenerateHTML(c.Writer, data, "layout", "mobile_ads", "show/list", "episode/item")
}

func showHandler(c *gin.Context) {
	db, _ := utils.OpenDB()
	defer db.Close()
	showID, _ := strconv.Atoi(c.Param("id"))
	show, _ := data.GetShow(&db, showID)
	if show.IsOtv && (os.Getenv("WATCH_OTV") == "1" || show.ChannelID == 3) {
		renderShowOtv(c, show)
	} else {
		renderShow(c, show)
	}
}

func showTvHandler(c *gin.Context) {
	db, _ := utils.OpenDB()
	defer db.Close()
	showID, _ := strconv.Atoi(c.Param("id"))
	show, _ := data.GetShow(&db, showID)
	renderShow(c, show)
}

func showOtvHandler(c *gin.Context) {
	db, _ := utils.OpenDB()
	defer db.Close()
	otvID, _ := strconv.Atoi(c.Param("id"))
	show, _ := data.GetShowByOtv(&db, otvID)
	if show.IsOtv {
		renderShowOtv(c, show)
	} else {
		renderShow(c, show)
	}
}

func renderShow(c *gin.Context, show data.Show) {
	db, _ := utils.OpenDB()
	defer db.Close()
	page, _ := strconv.Atoi(c.Query("page"))
	episodes, pageInfo, _ := data.GetEpisodesAndPageInfo(&db, show.ID, int32(page))
	// episodes, _ := data.GetEpisodes(&db, show.ID, 0)
	data := map[string]interface{}{
		"Title":    show.Title,
		"show":     show,
		"episodes": episodes,
		"pageInfo": pageInfo,
	}
	utils.GenerateHTML(c.Writer, data, "layout", "mobile_ads", "show/index", "episode/item")
}

func renderShowOtv(c *gin.Context, show data.Show) {
	_, episodes := data.GetOTVEpisodelist(show.OtvID)
	data := map[string]interface{}{
		"Title":    show.Title,
		"show":     show,
		"episodes": episodes,
	}
	utils.GenerateHTML(c.Writer, data, "layout", "mobile_ads", "show/otv_index")
}

func watchHandler(c *gin.Context) {
	db, _ := utils.OpenDB()
	defer db.Close()
	watchID, _ := strconv.Atoi(c.Param("watchID"))
	playIndex, _ := strconv.Atoi(c.Param("playIndex"))
	episode, err := data.GetEpisode(&db, watchID)
	if err != nil {
		// goOutHandler(r)
	}
	show, err := data.GetShow(&db, episode.ShowID)
	if maxIndex := len(episode.Playlists) - 1; maxIndex < playIndex {
		playIndex = maxIndex
	}

	if playIndex == -1 {
		c.Redirect(http.StatusMovedPermanently, "/not_found")
	}

	playlistItem := episode.Playlists[playIndex]
	var embedURL string
	if !episode.IsURL {
		switch episode.SrcType {
		case 1, 14:
			embedURL = playlistItem.Sources[0].File
		}
	}

	episodes, _ := data.GetEpisodes(&db, show.ID, 0)
	data := map[string]interface{}{
		"Title":        show.Title + " | " + episode.Title,
		"playIndex":    playIndex,
		"episode":      episode,
		"show":         show,
		"episodes":     episodes,
		"playlistItem": playlistItem,
		"embedURL":     embedURL,
		"isMobile":     utils.IsMobile(c.Request.UserAgent()),
	}
	utils.GenerateHTML(c.Writer, data, "layout", "mobile_ads", "watch/index")
}

func watchOtvHandler(c *gin.Context) {
	db, _ := utils.OpenDB()
	defer db.Close()
	isMobile := utils.IsMobile(c.Request.UserAgent())
	_, otvEpisodePlay := data.GetOTVEpisodePlay(c.Param("watchID"), isMobile)

	watchID, _ := strconv.Atoi(c.Param("watchID"))
	playIndex, _ := strconv.Atoi(c.Param("playIndex"))
	if maxIndex := len(otvEpisodePlay.EpisodeDetail.PartItems) - 1; maxIndex < playIndex {
		playIndex = maxIndex
	}

	if playIndex == -1 {
		c.Redirect(http.StatusMovedPermanently, "/not_found")
	}

	partItem := otvEpisodePlay.EpisodeDetail.PartItems[playIndex]
	partItem.IframeHTML = html.UnescapeString(partItem.IframeHTML)
	if !isMobile {
		partItem.IframeHTML = strings.Replace(partItem.IframeHTML, "/v/", "/playlist/", 1)
	}

	otvID, _ := strconv.Atoi(otvEpisodePlay.SeasonDetail.ContentSeasonID)
	show, _ := data.GetShowByOtv(&db, otvID)

	_, episodes := data.GetOTVEpisodelist(show.OtvID)
	data := map[string]interface{}{
		"Title":          partItem.Title,
		"partItem":       partItem,
		"otvEpisodePlay": otvEpisodePlay,
		"playIndex":      playIndex,
		"watchID":        watchID,
		"show":           show,
		"episodes":       episodes,
		"isMobile":       isMobile,
	}
	utils.GenerateHTML(c.Writer, data, "layout", "mobile_ads", "watch/otv_index")
}
