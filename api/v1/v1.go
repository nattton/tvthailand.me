package v1

import (
	"github.com/code-mobi/tvthailand.me/Godeps/_workspace/src/github.com/go-martini/martini"
	"github.com/code-mobi/tvthailand.me/Godeps/_workspace/src/github.com/jinzhu/gorm"
	"github.com/code-mobi/tvthailand.me/Godeps/_workspace/src/github.com/martini-contrib/render"
	"github.com/code-mobi/tvthailand.me/Godeps/_workspace/src/github.com/mssola/user_agent"
	"github.com/code-mobi/tvthailand.me/data"
	"net/http"
	"strconv"
)

func CategoriesHandler(db gorm.DB, r render.Render, params martini.Params) {
	categories, err := data.GetCategories(&db)
	if err != nil {
		r.JSON(http.StatusNotFound, err)
	}
	r.JSON(200, map[string]interface{}{
		"Categories": categories,
	})
}

func RecentlyHandler(db gorm.DB, r render.Render, params martini.Params) {
	start, _ := strconv.Atoi(params["start"])
	shows, err := data.GetShowByRecently(&db, start)
	if err != nil {
		r.JSON(http.StatusNotFound, err)
	}
	r.JSON(200, map[string]interface{}{
		"Shows": shows,
	})
}

func PopularHandler(db gorm.DB, r render.Render, params martini.Params) {
	start, _ := strconv.Atoi(params["start"])
	shows, err := data.GetShowByPopular(&db, start)
	if err != nil {
		r.JSON(http.StatusNotFound, err)
	}
	r.JSON(200, map[string]interface{}{
		"Shows": shows,
	})
}

func CategoryHandler(db gorm.DB, r render.Render, params martini.Params) {
	start, _ := strconv.Atoi(params["start"])
	shows, err := data.GetShowByCategory(&db, params["id"], start)
	if err != nil {
		r.JSON(http.StatusNotFound, err)
	}
	r.JSON(200, map[string]interface{}{
		"Shows": shows,
	})
}

func ChannelsHandler(db gorm.DB, r render.Render, params martini.Params) {
	categories, err := data.GetChannels(&db)
	if err != nil {
		r.JSON(http.StatusNotFound, err)
	}
	r.JSON(200, map[string]interface{}{
		"Channels": categories,
	})
}

func ChannelHandler(db gorm.DB, r render.Render, params martini.Params) {
	start, _ := strconv.Atoi(params["start"])
	shows, err := data.GetShowByChannel(&db, params["id"], start)
	if err != nil {
		r.JSON(http.StatusNotFound, err)
	}
	r.JSON(200, map[string]interface{}{
		"Shows": shows,
	})
}

func EpisodeHandler(db gorm.DB, r render.Render, params martini.Params) {
	episode, err := data.GetVideoList(&db, params["hashID"])
	if err != nil {
		r.JSON(http.StatusNotFound, err)
	}
	r.JSON(200, episode)
}

func WatchHandler(db gorm.DB, r render.Render, params martini.Params) {
	episode, err := data.GetVideoList(&db, params["hashID"])
	if err != nil {
		r.JSON(http.StatusNotFound, episode)
	}
	show, _ := data.GetShow(&db, episode.ShowID)
	r.JSON(200, map[string]interface{}{
		"Show":    show,
		"Episode": episode,
	})
}

func WatchOtvHandler(r render.Render, params martini.Params, req *http.Request) {
	ua := user_agent.New(req.UserAgent())
	otvEpisodePlay := data.GetOTVEpisodePlay(params["watchID"], ua.Mobile())
	r.JSON(200, otvEpisodePlay)
}
