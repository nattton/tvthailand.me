package main

import (
	"fmt"
	"html"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/code-mobi/tvthailand.me/data"
	"github.com/code-mobi/tvthailand.me/utils"
	"gopkg.in/gin-gonic/gin.v1"
)

// flashType : danger, warning, info
func printFlash(writer http.ResponseWriter, flashType, message string) {
	flash := map[string]string{
		flashType: message,
	}
	renderData := map[string]interface{}{
		"flash": flash,
	}
	utils.GenerateHTML(writer, renderData, "layout", "mobile_ads", "index")
}

func indexHandler(c *gin.Context) {
	fmt.Println(c.Request.Host)
	db, _ := utils.OpenDB()
	defer db.Close()
	recents := make(chan []data.Show)
	populars := make(chan []data.Show)
	go func() {
		shows, _ := data.ShowsRecently(db, 0)
		recents <- shows
	}()
	go func() {
		shows, _ := data.ShowsPopular(db, 0)
		populars <- shows
	}()
	categories, _ := data.CategoriesActive(db)
	channels, _ := data.ChannelsActive(db)

	name, _ := c.GetQuery("name")
	fmt.Println("name :", name)
	renderData := map[string]interface{}{
		"host":         c.Request.Host,
		"Description":  "ดูรายการทีวี ละครย้อนหลัง",
		"showRecents":  <-recents,
		"showPopulars": <-populars,
		"categories":   categories,
		"channels":     channels,
		"isMobile":     utils.IsMobileNotPad(c.Request.UserAgent()),
		"isOMU":        name == "omu",
	}
	utils.GenerateHTML(c.Writer, renderData, "layout", "mobile_ads", "index")
}

func mobileAppsHandler(c *gin.Context) {
	renderData := map[string]interface{}{
		"Title":       "Mobile Apps",
		"Description": "Download Apps TV Thailand ได้ที่นี่",
	}
	utils.GenerateHTML(c.Writer, renderData, "layout", "mobile_ads", "static/mobile_apps")
}

func notFoundHandler(c *gin.Context) {
	printFlash(c.Writer, "danger", "Page not found")
}

func goNotFound(c *gin.Context) {
	c.Redirect(http.StatusMovedPermanently, "/not_found")
}

func recentlyHandler(c *gin.Context) {
	db, _ := utils.OpenDB()
	defer db.Close()
	shows, _ := data.ShowsRecently(db, 0)
	renderData := map[string]interface{}{
		"Title":       "รายการล่าสุด",
		"Description": "รายการล่าสุด",
		"header":      "รายการล่าสุด",
		"typeMode":    "recently",
		"shows":       shows,
		"isMobile":    utils.IsMobileNotPad(c.Request.UserAgent()),
	}
	utils.GenerateHTML(c.Writer, renderData, "layout", "mobile_ads", "show/list", "episode/item")
}

func popularHandler(c *gin.Context) {
	db, _ := utils.OpenDB()
	defer db.Close()
	shows, _ := data.ShowsPopular(db, 0)
	renderData := map[string]interface{}{
		"Title":       "Popular",
		"Description": "รายการยอดนิยม",
		"header":      "Popular",
		"typeMode":    "popular",
		"shows":       shows,
		"isMobile":    utils.IsMobileNotPad(c.Request.UserAgent()),
	}
	utils.GenerateHTML(c.Writer, renderData, "layout", "mobile_ads", "show/list", "episode/item")
}

func categoriesHandler(c *gin.Context) {
	db, _ := utils.OpenDB()
	defer db.Close()
	categories, _ := data.CategoriesActive(db)
	renderData := map[string]interface{}{
		"Title":       "หมวด",
		"Description": "หมวดทั้งหมด",
		"header":      "หมวด",
		"categories":  categories,
		"isMobile":    utils.IsMobileNotPad(c.Request.UserAgent()),
	}
	utils.GenerateHTML(c.Writer, renderData, "layout", "mobile_ads", "category/list")
}

func categoryShowHandler(c *gin.Context) {
	titlize := c.Param("titlize")
	db, _ := utils.OpenDB()
	defer db.Close()
	category, err := data.CategoryTitleze(db, titlize)
	if err != nil {
		goNotFound(c)
		return
	}
	shows, _ := data.ShowsCategory(db, category.ID, 0)

	renderData := map[string]interface{}{
		"Title":       category.Title,
		"Description": category.Description,
		"Image":       category.Thumbnail,
		"header":      category.Title,
		"typeMode":    "category",
		"typeId":      category.ID,
		"shows":       shows,
		"isMobile":    utils.IsMobileNotPad(c.Request.UserAgent()),
	}
	utils.GenerateHTML(c.Writer, renderData, "layout", "mobile_ads", "show/list", "episode/item")
}

func channelsHandler(c *gin.Context) {
	db, _ := utils.OpenDB()
	defer db.Close()
	channels, _ := data.ChannelsActive(db)
	renderData := map[string]interface{}{
		"Title":       "ช่องทีวี / Live",
		"Description": "ดูทีวี / Live / รายการสด",
		"header":      "ช่องทีวี / Live",
		"channels":    channels,
		"isMobile":    utils.IsMobileNotPad(c.Request.UserAgent()),
	}
	utils.GenerateHTML(c.Writer, renderData, "layout", "mobile_ads", "channel/list")
}

func channelShowHandler(c *gin.Context) {
	id := c.Param("id")
	db, _ := utils.OpenDB()
	defer db.Close()
	channel, _ := data.GetChannel(db, id)
	shows, _ := data.ShowsChannel(db, channel.ID, 0)
	renderData := map[string]interface{}{
		"Title":       channel.Title,
		"Description": channel.Description,
		"Image":       channel.Thumbnail,
		"header":      channel.Title,
		"channel":     channel,
		"typeMode":    "category",
		"typeId":      channel.ID,
		"shows":       shows,
	}
	utils.GenerateHTML(c.Writer, renderData, "layout", "mobile_ads", "show/list", "episode/item")
}

func searchShowHandler(c *gin.Context) {
	db, _ := utils.OpenDB()
	defer db.Close()
	qs := c.Request.URL.Query()
	keyword := qs.Get("keyword")
	chShows := make(chan []data.Show)
	chEpisodes := make(chan []data.Episode)

	renderData := map[string]interface{}{}
	if keyword != "" {
		header := "ผลการค้นหา : " + keyword
		renderData["Title"] = header
		renderData["header"] = header
		go func() {
			shows, _ := data.ShowsSearch(db, keyword)
			chShows <- shows
		}()
		go func() {
			episodes, _ := data.GetEpisodesBySearch(db, keyword)
			chEpisodes <- episodes
		}()

		renderData["shows"] = <-chShows
		renderData["episodes"] = <-chEpisodes
	} else {
		renderData["Title"] = "Search"
		renderData["header"] = "กรุณากรอกชื่อเรื่องที่ต้องการค้นหา"
	}
	utils.GenerateHTML(c.Writer, renderData, "layout", "mobile_ads", "show/list", "episode/item")
}

func showHandler(c *gin.Context) {
	db, _ := utils.OpenDB()
	defer db.Close()
	showID, _ := strconv.Atoi(c.Param("id"))
	show, err := data.GetShow(db, showID)
	if err != nil {
		log.Println(err)
		goNotFound(c)
		return
	}
	if !show.Web {
		goNotFound(c)
		return
	}

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
	show, _ := data.GetShow(db, showID)
	if !show.Web {
		goNotFound(c)
		return
	}
	renderShow(c, show)
}

func showOtvHandler(c *gin.Context) {
	db, _ := utils.OpenDB()
	defer db.Close()
	otvID, _ := strconv.Atoi(c.Param("id"))
	show, _ := data.ShowWithOtv(db, otvID)
	if show.ID == 0 {
		goNotFound(c)
		return
	}
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
	episodes, pageInfo, _ := data.EpisodesAndPageInfo(db, show.ID, int32(page))
	title := fmt.Sprintf("ดู %s ย้อนหลัง", show.Title)
	description := fmt.Sprintf("%s | %s", show.Title, show.Description)
	renderData := map[string]interface{}{
		"Title":       title,
		"Description": description,
		"Image":       show.Thumbnail,
		"show":        show,
		"episodes":    episodes,
		"pageInfo":    pageInfo,
	}
	utils.GenerateHTML(c.Writer, renderData, "layout", "mobile_ads", "show/index", "episode/item")
}

func renderShowOtv(c *gin.Context, show data.Show) {
	limit := 20
	page, _ := strconv.Atoi(c.Query("page"))
	var offset int
	if page <= 1 {
		page = 1
		offset = 0
	} else {
		offset = (page - 1) * limit
	}
	_, episodes, err := data.GetOTVEpisodelist(show.OtvID, offset, limit)
	pageInfo := data.PageInfo{}
	if len(episodes.EpisodeList) == limit {
		pageInfo.NextPage = int32(page + 1)
	}
	if page > 1 {
		pageInfo.PreviousPage = int32(page - 1)
	}
	if err != nil {
		printFlash(c.Writer, "danger", "OTV Server Error")
		return
	}
	title := fmt.Sprintf("ดู %s ย้อนหลัง", show.Title)
	description := fmt.Sprintf("%s | %s", show.Title, show.Description)
	renderData := map[string]interface{}{
		"Title":       title,
		"Description": description,
		"Image":       show.Thumbnail,
		"show":        show,
		"episodes":    episodes,
		"pageInfo":    pageInfo,
	}
	utils.GenerateHTML(c.Writer, renderData, "layout", "mobile_ads", "show/otv_index", "episode/otv_item")
}

func watchHandler(c *gin.Context) {
	db, _ := utils.OpenDB()
	defer db.Close()
	watchID, _ := strconv.Atoi(c.Param("watchID"))
	playIndex, _ := strconv.Atoi(c.Param("playIndex"))
	episode, err := data.GetEpisode(db, watchID)
	if err != nil {
		goNotFound(c)
		return
	}
	show, err := data.GetShow(db, episode.ShowID)
	if !show.Web {
		goNotFound(c)
		return
	}

	if maxIndex := len(episode.Playlists) - 1; maxIndex < playIndex {
		playIndex = maxIndex
	}

	if playIndex == -1 {
		goNotFound(c)
		return
	}

	playlistItem := episode.Playlists[playIndex]
	var embedURL string
	if !episode.IsURL {
		switch episode.SrcType {
		case 1, 14:
			embedURL = playlistItem.Sources[0].File
		}
	}

	episodes, _ := data.GetEpisodes(db, show.ID, 0)
	title := fmt.Sprintf("%s | %s", show.Title, episode.Title)
	renderData := map[string]interface{}{
		"Title":        title,
		"Description":  title,
		"Image":        playlistItem.Image,
		"playIndex":    playIndex,
		"episode":      episode,
		"show":         show,
		"episodes":     episodes,
		"playlistItem": playlistItem,
		"embedURL":     embedURL,
		"isMobile":     utils.IsMobile(c.Request.UserAgent()),
	}
	utils.GenerateHTML(c.Writer, renderData, "layout", "mobile_ads", "watch/index", "episode/item")
}

func watchOtvHandler(c *gin.Context) {
	limit := 20
	db, _ := utils.OpenDB()
	defer db.Close()
	isMobile := utils.IsMobileNotPad(c.Request.UserAgent())
	_, otvEpisodePlay, err := data.GetOTVEpisodePlay(c.Param("watchID"), isMobile)
	if err != nil {
		printFlash(c.Writer, "danger", "OTV Server Error")
		return
	}

	watchID, _ := strconv.Atoi(c.Param("watchID"))
	playIndex, _ := strconv.Atoi(c.Param("playIndex"))
	if maxIndex := len(otvEpisodePlay.EpisodeDetail.PartItems) - 1; maxIndex < playIndex {
		playIndex = maxIndex
	}

	if playIndex == -1 {
		goNotFound(c)
		return
	}

	partItem := otvEpisodePlay.EpisodeDetail.PartItems[playIndex]
	partItem.IframeHTML = html.UnescapeString(partItem.IframeHTML)
	// if !isMobile {
	// 	r := strings.NewReplacer("/v/", "/playlist/",
	// 		"iframe", `iframe class="embed-responsive-item"`)
	// 	partItem.IframeHTML = r.Replace(partItem.IframeHTML)
	// }
	r := strings.NewReplacer("iframe", `iframe class="embed-responsive-item"`)
	partItem.IframeHTML = r.Replace(partItem.IframeHTML)

	otvID, _ := strconv.Atoi(otvEpisodePlay.SeasonDetail.ContentSeasonID)
	show, _ := data.ShowWithOtv(db, otvID)
	page, _ := strconv.Atoi(c.Query("page"))
	var offset int
	if page <= 1 {
		page = 1
		offset = 0
	} else {
		offset = (page - 1) * limit
	}
	_, episodes, err := data.GetOTVEpisodelist(show.OtvID, offset, limit)
	if err != nil {
		printFlash(c.Writer, "danger", "OTV Server Error")
		return
	}
	pageInfo := data.PageInfo{}
	if len(episodes.EpisodeList) == limit {
		pageInfo.NextPage = int32(page + 1)
	}
	if page > 1 {
		pageInfo.PreviousPage = int32(page - 1)
	}

	title := fmt.Sprintf("ดู %s ย้อนหลัง", partItem.Title)
	description := fmt.Sprintf("%s | %s", partItem.Title, otvEpisodePlay.EpisodeDetail.Detail)
	renderData := map[string]interface{}{
		"Title":          title,
		"Description":    description,
		"Image":          partItem.Thumbnail,
		"partItem":       partItem,
		"otvEpisodePlay": otvEpisodePlay,
		"playIndex":      playIndex,
		"watchID":        watchID,
		"show":           show,
		"episodes":       episodes,
		"isMobile":       isMobile,
		"pageInfo":       pageInfo,
	}
	utils.GenerateHTML(c.Writer, renderData, "layout", "mobile_ads", "watch/otv_index", "episode/otv_item")
}

func OPlayHandler(c *gin.Context) {
	responseBody, _, _ := data.GetOTVEpisodePlay(c.Param("watchID"), false)
	fmt.Fprintf(c.Writer, string(responseBody))
}
