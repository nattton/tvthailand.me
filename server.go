package main

import (
	"encoding/json"
	"github.com/code-mobi/tvthailand.me/api/v1"
	"github.com/go-martini/martini"
	"github.com/jinzhu/gorm"
	"github.com/martini-contrib/render"
	"html/template"
	"log"
	"os"
	"reflect"
	"strings"
)

func main() {
	db, err := gorm.Open("mysql", os.Getenv("DATABASE_URL_PATH"))
	if err != nil {
		log.Fatal(err)
	}
	db.LogMode(true)
	defer db.Close()

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
				"toJson": func(a interface{}) string {
					b, _ := json.Marshal(a)
					r := strings.NewReplacer("\\", "")
					return r.Replace(string(b))
				},
			},
		},
	}))

	m.Get("/", indexHandler)
	m.Get("/categories", categoryHandler)
	m.Get("/category/:titlize", categoryShowHandler)
	// m.Get("/channel", channel)
	m.Get("/channel/:id", channelShowHandler)
	m.Get("/channel/:id/:title", channelShowHandler)
	m.Get("/search", searchShowHandler)
	m.Get("/show/:id", showHandler)
	m.Get("/show/:id/:title", showHandler)
	m.Get("/show_otv/:id", showOtvHandler)
	m.Get("/show_otv/:id/:title", showOtvHandler)
	m.Get("/watch/(?P<watchID>[0-9]+)", watchHandler)
	m.Get("/watch/(?P<watchID>[0-9]+)/(?P<playIndex>[0-9]+)", watchHandler)
	m.Get("/watch_otv/(?P<watchID>[0-9]+)", watchOtvHandler)
	m.Get("/watch_otv/(?P<watchID>[0-9]+)/(?P<playIndex>[0-9]+)", watchOtvHandler)
	m.Group("/admin", func(r martini.Router) {
		m.Get("/encrypt", encryptHandler)
	})

	m.Group("/api/v1", func(r martini.Router) {
		r.Get("/categories", v1.CategoriesHandler)
		r.Get("/category/:id", v1.CategoryHandler)
		r.Get("/channels", v1.ChannelsHandler)
		r.Get("/channel/:id", v1.ChannelHandler)
		r.Get("/episode/:hashID", v1.EpisodeHandler)
		r.Get("/watch/:hashID", v1.WatchHandler)
		m.Get("/watch_otv/(?P<watchID>[0-9]+)", v1.WatchOtvHandler)
	})
	m.Get("/not_found", notFoundHandler)
	m.NotFound(notFoundHandler)
	m.Run()
}
