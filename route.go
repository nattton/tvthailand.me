package main

import (
	"github.com/code-mobi/tvthailand.me/Godeps/_workspace/src/github.com/go-martini/martini"
	"github.com/code-mobi/tvthailand.me/Godeps/_workspace/src/github.com/jinzhu/gorm"
	"github.com/code-mobi/tvthailand.me/Godeps/_workspace/src/github.com/martini-contrib/render"
	"github.com/code-mobi/tvthailand.me/Godeps/_workspace/src/github.com/mssola/user_agent"
	"github.com/code-mobi/tvthailand.me/data"
	"net/http"
	"strconv"
	"strings"
)

func indexHandler(db gorm.DB, r render.Render) {
	recents, _ := data.GetShowByRecently(&db)
	populars, _ := data.GetShowByPopular(&db)
	r.HTML(http.StatusOK, "index", map[string]interface{}{
		"showRecents":  recents,
		"showPopulars": populars,
	})
}

func notFoundHandler(r render.Render) {
	r.HTML(http.StatusOK, "not_found", nil)
}

func goOutHandler(r render.Render) {
	r.Redirect("/not_found", http.StatusMovedPermanently)
}

func encryptHandler(db gorm.DB, r render.Render) {
	data.EncryptEpisode(&db)
	r.HTML(http.StatusOK, "index", map[string]interface{}{
		"header": "Encrypt Successfully",
	})
}

func categoryHandler(db gorm.DB, r render.Render) {
	categories, _ := data.GetCategories(&db)
	r.HTML(http.StatusOK, "category_channel", map[string]interface{}{
		"header":     "หมวด",
		"categories": categories,
	})
}

func categoryShowHandler(db gorm.DB, r render.Render, params martini.Params) {
	titlize := params["titlize"]
	category, _ := data.GetCategory(&db, titlize)
	shows, _ := data.GetShowByCategory(&db, category.ID)
	populars, _ := data.GetShowByCategoryPopular(&db, category.ID)
	r.HTML(http.StatusOK, "category_channel", map[string]interface{}{
		"Title":    category.Title,
		"header":   category.Title,
		"shows":    shows,
		"populars": populars,
	})
}

func channelShowHandler(db gorm.DB, r render.Render, params martini.Params) {
	id := params["id"]
	channel, _ := data.GetChannel(&db, id)
	shows, _ := data.GetShowByChannel(&db, id)
	r.HTML(http.StatusOK, "category_channel", map[string]interface{}{
		"Title":  channel.Title,
		"header": channel.Title,
		"shows":  shows,
	})
}

func searchShowHandler(db gorm.DB, r render.Render, params martini.Params, req *http.Request) {
	qs := req.URL.Query()
	keyword := qs.Get("keyword")
	var shows []data.Show
	var header string
	var title string
	if keyword != "" {
		shows = data.GetShowBySearch(&db, keyword)
		header = "ผลการค้นหา : " + keyword
		title = header
	} else {
		title = "Search"
		header = "กรุณาพิมพชื่อเรื่องที่ต้องการค้นหา"
	}
	r.HTML(http.StatusOK, "category_channel", map[string]interface{}{
		"Title":   title,
		"keyword": keyword,
		"header":  header,
		"shows":   shows,
	})
}

func showHandler(db gorm.DB, r render.Render, params martini.Params) {
	showID, _ := strconv.Atoi(params["id"])
	show, _ := data.GetShow(&db, showID)
	if show.IsOtv {
		renderShowOtv(db, r, show)
	} else {
		renderShow(db, r, show)
	}
}

func showOtvHandler(db gorm.DB, r render.Render, params martini.Params) {
	otvID, _ := strconv.Atoi(params["id"])
	show, _ := data.GetShowByOtv(&db, otvID)
	if show.IsOtv {
		renderShowOtv(db, r, show)
	} else {
		renderShow(db, r, show)
	}
}

func renderShow(db gorm.DB, r render.Render, show data.Show) {
	episodes := data.GetEpisodes(&db, show.ID)
	r.HTML(http.StatusOK, "show/index", map[string]interface{}{
		"Title":    show.Title,
		"show":     show,
		"episodes": episodes,
	})
}

func renderShowOtv(db gorm.DB, r render.Render, show data.Show) {
	episodes := data.GetOTVEpisodelist(show.OtvID)
	r.HTML(http.StatusOK, "show/otv_index", map[string]interface{}{
		"Title":    show.Title,
		"show":     show,
		"episodes": episodes,
	})
}

func watchHandler(db gorm.DB, r render.Render, params martini.Params) {
	watchID, _ := strconv.Atoi(params["watchID"])
	playIndex, _ := strconv.Atoi(params["playIndex"])
	episode, err := data.GetEpisode(&db, watchID)
	if err != nil {
		goOutHandler(r)
	}
	show, err := data.GetShow(&db, episode.ShowID)
	r.HTML(http.StatusOK, "watch/index", map[string]interface{}{
		"Title":     show.Title + " | " + episode.Title,
		"playIndex": playIndex,
		"episode":   episode,
	})
}

func watchOtvHandler(db gorm.DB, r render.Render, params martini.Params, req *http.Request) {
	ua := user_agent.New(req.UserAgent())
	otvEpisodePlay := data.GetOTVEpisodePlay(params["watchID"], ua.Mobile())
	watchID, _ := strconv.Atoi(params["watchID"])
	playIndex, _ := strconv.Atoi(params["playIndex"])
	partItem := otvEpisodePlay.EpisodeDetail.PartItems[playIndex]
	partItem.IframeHTML = strings.Replace(strings.Replace(partItem.IframeHTML, "&lt;", "<", -1), "&gt;", ">", -1)
	r.HTML(http.StatusOK, "watch/otv_index", map[string]interface{}{
		"Title":          partItem.Title,
		"partItem":       partItem,
		"otvEpisodePlay": otvEpisodePlay,
		"playIndex":      playIndex,
		"watchID":        watchID,
	})
}
