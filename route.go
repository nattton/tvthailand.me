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

func indexHandler(db gorm.DB, r render.Render, req *http.Request) {
	recents, _ := data.GetShowByRecently(&db, 0)
	populars, _ := data.GetShowByPopular(&db, 0)
	r.HTML(http.StatusOK, "index", map[string]interface{}{
		"showRecents":  recents,
		"showPopulars": populars,
		"isMobile":     user_agent.New(req.UserAgent()).Mobile(),
	})
}

func notFoundHandler(r render.Render) {
	r.HTML(http.StatusOK, "not_found", nil)
}

func goOutHandler(r render.Render) {
	r.Redirect("/not_found", http.StatusMovedPermanently)
}

func encryptHandler(db gorm.DB, r render.Render, params martini.Params) {
	episodeID, _ := strconv.Atoi(params["episodeID"])
	data.EncryptEpisode(&db, episodeID)
	r.HTML(http.StatusOK, "index", map[string]interface{}{
		"header": "Encrypt Successfully",
	})
}

func recentlyHandler(db gorm.DB, r render.Render, req *http.Request) {
	shows, _ := data.GetShowByRecently(&db, 0)
	r.HTML(http.StatusOK, "show/list", map[string]interface{}{
		"Title":    "รายการล่าสุด",
		"header":   "รายการล่าสุด",
		"apiPath":  "/recently/",
		"shows":    shows,
		"isMobile": user_agent.New(req.UserAgent()).Mobile(),
	})
}

func popularHandler(db gorm.DB, r render.Render, req *http.Request) {
	shows, _ := data.GetShowByPopular(&db, 0)
	r.HTML(http.StatusOK, "show/list", map[string]interface{}{
		"Title":    "Popular",
		"header":   "Popular",
		"apiPath":  "/popular/",
		"shows":    shows,
		"isMobile": user_agent.New(req.UserAgent()).Mobile(),
	})
}

func categoriesHandler(db gorm.DB, r render.Render, req *http.Request) {
	categories, _ := data.GetCategories(&db)
	r.HTML(http.StatusOK, "category/list", map[string]interface{}{
		"header":     "หมวด",
		"categories": categories,
		"isMobile":   user_agent.New(req.UserAgent()).Mobile(),
	})
}

func categoryShowHandler(db gorm.DB, r render.Render, params martini.Params) {
	titlize := params["titlize"]
	start, _ := strconv.Atoi(params["start"])
	category, _ := data.GetCategory(&db, titlize)
	shows, _ := data.GetShowByCategory(&db, category.ID, start)
	r.HTML(http.StatusOK, "show/list", map[string]interface{}{
		"Title":   category.Title,
		"header":  category.Title,
		"apiPath": "/category/" + category.ID + "/",
		"shows":   shows,
	})
}

func channelsHandler(db gorm.DB, r render.Render) {
	channels, _ := data.GetChannels(&db)
	r.HTML(http.StatusOK, "channel/list", map[string]interface{}{
		"header":   "ช่องทีวี",
		"channels": channels,
	})
}

func channelShowHandler(db gorm.DB, r render.Render, params martini.Params) {
	id := params["id"]
	start, _ := strconv.Atoi(params["start"])
	channel, _ := data.GetChannel(&db, id)
	shows, _ := data.GetShowByChannel(&db, channel.ID, start)
	r.HTML(http.StatusOK, "show/list", map[string]interface{}{
		"Title":   channel.Title,
		"header":  channel.Title,
		"apiPath": "/channel/" + channel.ID + "/",
		"shows":   shows,
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
	r.HTML(http.StatusOK, "show/list", map[string]interface{}{
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

func watchHandler(db gorm.DB, r render.Render, params martini.Params, req *http.Request) {
	watchID, _ := strconv.Atoi(params["watchID"])
	playIndex, _ := strconv.Atoi(params["playIndex"])
	episode, err := data.GetEpisode(&db, watchID)
	if err != nil {
		goOutHandler(r)
	}
	show, err := data.GetShow(&db, episode.ShowID)
	if maxIndex := len(episode.Playlists) - 1; maxIndex < playIndex {
		playIndex = maxIndex
	}
	playlistItem := episode.Playlists[playIndex]
	r.HTML(http.StatusOK, "watch/index", map[string]interface{}{
		"Title":        show.Title + " | " + episode.Title,
		"playIndex":    playIndex,
		"episode":      episode,
		"show":         show,
		"playlistItem": playlistItem,
		"isMobile":     user_agent.New(req.UserAgent()).Mobile(),
	})
}

func watchOtvHandler(db gorm.DB, r render.Render, params martini.Params, req *http.Request) {
	ua := user_agent.New(req.UserAgent())
	isMobile := ua.Mobile()
	otvEpisodePlay := data.GetOTVEpisodePlay(params["watchID"], isMobile)
	watchID, _ := strconv.Atoi(params["watchID"])
	playIndex, _ := strconv.Atoi(params["playIndex"])
	if maxIndex := len(otvEpisodePlay.EpisodeDetail.PartItems) - 1; maxIndex < playIndex {
		playIndex = maxIndex
	}
	partItem := otvEpisodePlay.EpisodeDetail.PartItems[playIndex]
	partItem.IframeHTML = strings.Replace(strings.Replace(partItem.IframeHTML, "&lt;", "<", -1), "&gt;", ">", -1)
	r.HTML(http.StatusOK, "watch/otv_index", map[string]interface{}{
		"Title":          partItem.Title,
		"partItem":       partItem,
		"otvEpisodePlay": otvEpisodePlay,
		"playIndex":      playIndex,
		"watchID":        watchID,
		"isMobile":       isMobile,
	})
}
