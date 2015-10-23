package main

import (
	"encoding/json"
	"github.com/code-mobi/tvthailand.me/Godeps/_workspace/src/github.com/go-martini/martini"
	"github.com/code-mobi/tvthailand.me/Godeps/_workspace/src/github.com/jinzhu/gorm"
	"github.com/code-mobi/tvthailand.me/Godeps/_workspace/src/github.com/martini-contrib/auth"
	"github.com/code-mobi/tvthailand.me/Godeps/_workspace/src/github.com/martini-contrib/render"
	"github.com/code-mobi/tvthailand.me/admin"
	"github.com/code-mobi/tvthailand.me/api/v1"
	"html"
	"html/template"
	"log"
	"net/url"
	"os"
	"reflect"
	"strings"
)

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "3000"
	}
	db, err := gorm.Open("mysql", os.Getenv("DATABASE_DSN"))
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()
	db.LogMode(martini.Env != "production")

	m := martini.Classic()
	m.Map(db)
	m.Use(render.Renderer(render.Options{
		Directory:  "templates",
		Layout:     "layout",
		Extensions: []string{".tmpl", ".html"},
		Delims:     render.Delims{"{[{", "}]}"},
		Charset:    "UTF-8",
		IndentJSON: false,
		Funcs: []template.FuncMap{
			{
				"last": func(x int, a interface{}) bool {
					return x == reflect.ValueOf(a).Len()-1
				},
			},
			{
				"escStr": func(a ...string) string {
					return html.EscapeString(strings.Join(a, "-"))
				},
			},
			{
				"urlEsc": func(a ...string) string {
					return url.QueryEscape(strings.Join(a, "-"))
				},
			},
			{
				"toJson": func(a interface{}) string {
					b, _ := json.Marshal(a)
					r := strings.NewReplacer("\\", "")
					return r.Replace(string(b))
				},
			},
		},
	}))

	authAdmin := auth.BasicFunc(func(username, password string) bool {
		return auth.SecureCompare(username, "saly") && auth.SecureCompare(password, "admin888")
	})

	m.Get("/", indexHandler)
	m.Get("/recently", recentlyHandler)
	m.Get("/popular", popularHandler)
	m.Get("/categories", categoriesHandler)
	m.Get("/category/:titlize", categoryShowHandler)
	m.Get("/channels", channelsHandler)
	m.Get("/channel/:id", channelShowHandler)
	m.Get("/channel/:id/**", channelShowHandler)
	m.Get("/search", searchShowHandler)
	m.Get("/show/:id", showHandler)
	m.Get("/show/:id/**", showHandler)
	m.Get("/show_tv/:id", showTvHandler)
	m.Get("/show_tv/:id/**", showTvHandler)
	m.Get("/show_otv/:id", showOtvHandler)
	m.Get("/show_otv/:id/**", showOtvHandler)
	m.Get("/watch/(?P<watchID>[0-9]+)/(?P<playIndex>[0-9]+)", watchHandler)
	m.Get("/watch/(?P<watchID>[0-9]+)/(?P<playIndex>[0-9]+)/**", watchHandler)
	m.Get("/watch_otv/(?P<watchID>[0-9]+)/(?P<playIndex>[0-9]+)", watchOtvHandler)
	m.Get("/watch_otv/(?P<watchID>[0-9]+)/(?P<playIndex>[0-9]+)/**", watchOtvHandler)
	m.Group("/admin", func(r martini.Router) {
		m.Get("/", authAdmin, admin.IndexHandler)
		m.Get("/encrypt_episode", admin.EncryptEpisodeHandler)
		m.Get("/encrypt_episode/:episodeID", authAdmin, admin.EncryptEpisodeHandler)
		m.Get("/mthai_embed/(?P<showID>[0-9]+)", authAdmin, admin.GetEmbedMThaiHandler)
	})

	m.Group("/api/v1", func(r martini.Router) {
		r.Get("/recently/:start", v1.RecentlyHandler)
		r.Get("/popular/:start", v1.PopularHandler)
		r.Get("/categories", v1.CategoriesHandler)
		r.Get("/category/:id", v1.CategoryHandler)
		r.Get("/category/:id/(?P<start>[0-9]+)", v1.CategoryHandler)
		r.Get("/channels", v1.ChannelsHandler)
		r.Get("/channel/:id", v1.ChannelHandler)
		r.Get("/channel/:id/(?P<start>[0-9]+)", v1.ChannelHandler)
		r.Get("/episode/:hashID", v1.EpisodeHandler)
		r.Get("/watch/:hashID", v1.WatchHandler)
		m.Get("/watch_otv/(?P<watchID>[0-9]+)", v1.WatchOtvHandler)
	})
	m.Get("/mobile_apps", func(r render.Render) {
		r.HTML(200, "static/mobile_apps", nil)
	})
	m.Get("/not_found", notFoundHandler)
	m.NotFound(notFoundHandler)
	m.RunOnAddr(":" + port)
}
