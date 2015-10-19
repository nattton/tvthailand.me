package v1

import (
	"github.com/code-mobi/tvthailand.me/Godeps/_workspace/src/github.com/go-martini/martini"
	"github.com/code-mobi/tvthailand.me/Godeps/_workspace/src/github.com/jinzhu/gorm"
	"github.com/code-mobi/tvthailand.me/Godeps/_workspace/src/github.com/martini-contrib/render"
	"github.com/code-mobi/tvthailand.me/Godeps/_workspace/src/github.com/mssola/user_agent"
	"github.com/code-mobi/tvthailand.me/data"
	"net/http"
)

func CategoriesHandler(db gorm.DB, r render.Render, params martini.Params) {
	categories, err := data.GetCategories(&db)
	if err != nil {
		r.JSON(http.StatusNotFound, err)
	}
	r.JSON(200, map[string]interface{}{
		"categories": categories,
	})
}

func CategoryHandler(db gorm.DB, r render.Render, params martini.Params) {
	shows, err := data.GetShowByCategory(&db, params["id"])
	if err != nil {
		r.JSON(http.StatusNotFound, err)
	}
	r.JSON(200, map[string]interface{}{
		"shows": shows,
	})
}

func ChannelsHandler(db gorm.DB, r render.Render, params martini.Params) {
	categories, err := data.GetChannels(&db)
	if err != nil {
		r.JSON(http.StatusNotFound, err)
	}
	r.JSON(200, map[string]interface{}{
		"channels": categories,
	})
}

func ChannelHandler(db gorm.DB, r render.Render, params martini.Params) {
	shows, err := data.GetShowByChannel(&db, params["id"])
	if err != nil {
		r.JSON(http.StatusNotFound, err)
	}
	r.JSON(200, map[string]interface{}{
		"shows": shows,
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
		"show":    show,
		"episode": episode,
	})
}

func WatchOtvHandler(r render.Render, params martini.Params, req *http.Request) {
	ua := user_agent.New(req.UserAgent())
	otvEpisodePlay := data.GetOTVEpisodePlay(params["watchID"], ua.Mobile())
	r.JSON(200, otvEpisodePlay)
}
