package main

import (
	"flag"
	"github.com/code-mobi/tvthailand.me/Godeps/_workspace/src/github.com/gin-gonic/gin"
	"github.com/code-mobi/tvthailand.me/admin"
	"os"
)

var commandParam CommandParam

func init() {
	flag.StringVar(&commandParam.Command, "command", "", "COMMAND = runbotch [-channel] [-q] | runbotpl [-playlist] | updateuser | migrate_botvideo | mthaithumbnail")
	flag.StringVar(&commandParam.Channel, "channel", "", "CHANNEL")
	flag.StringVar(&commandParam.Playlist, "playlist", "", "Playlist")
	flag.StringVar(&commandParam.Query, "q", "", "QUERY")
	flag.IntVar(&commandParam.Start, "start", 0, "START")
	flag.IntVar(&commandParam.Stop, "stop", 0, "STOP")
	flag.Parse()
}

func main() {
	if commandParam.Command != "" {
		processCommand(commandParam)
	} else {
		runServer()
	}
}

func runServer() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "3000"
	}

	router := gin.New()
	router.Use(gin.Logger())
	router.Static("/bower_components", "./public/bower_components")
	router.Static("/static", "./public/static")
	router.Static("/favicon", "./public/favicon")
	router.StaticFile("/favicon.ico", "./public/favicon.ico")
	router.StaticFile("/robot.txt", "./public/robot.txt")

	router.GET("/", indexHandler)
	router.GET("/recently", recentlyHandler)
	router.GET("/popular", popularHandler)
	router.GET("/categories", categoriesHandler)
	router.GET("/category/:titlize", categoryShowHandler)
	router.GET("/channels", channelsHandler)
	router.GET("/channel/:id", channelShowHandler)
	router.GET("/channel/:id/*title", channelShowHandler)
	router.GET("/search", searchShowHandler)
	router.GET("/show/:id", showHandler)
	router.GET("/show/:id/*title", showHandler)
	router.GET("/show_tv/:id/*title", showTvHandler)
	router.GET("/show_otv/:id/*title", showOtvHandler)
	router.GET("/watch/:watchID/", watchHandler)
	router.GET("/watch/:watchID/:playIndex", watchHandler)
	router.GET("/watch/:watchID/:playIndex/*title", watchHandler)
	router.GET("/watch_otv/:watchID/", watchOtvHandler)
	router.GET("/watch_otv/:watchID/:playIndex", watchOtvHandler)
	router.GET("/watch_otv/:watchID/:playIndex/*title", watchOtvHandler)
	router.GET("/oplay/:watchID/*title", OPlayHandler)
	router.GET("/mobile_apps", mobileAppsHandler)

	routerAjax := router.Group("/ajax")
	{
		routerAjax.GET("/recently", AjaxRecentlyHandler)
		routerAjax.GET("/popular", AjaxPopularHandler)
		routerAjax.GET("/category/:id", AjaxCategoryHandler)
		routerAjax.GET("/channels", AjaxChannelsHandler)
		routerAjax.GET("/channel/:id", AjaxChannelHandler)
		routerAjax.GET("/show/:show_id/episodes", AjaxShowHandler)
	}

	routerAuthorized := router.Group("/admin", gin.BasicAuth(gin.Accounts{
		"saly":    "admin888",
		"lucifer": "gundamman",
	}))
	routerAuthorized.GET("/", admin.IndexHandler)
	routerAuthorized.GET("/encrypt_episode", admin.EncryptEpisodeHandler)
	routerAuthorized.GET("/encrypt_episode/:episodeID", admin.EncryptEpisodeHandler)
	routerAuthorized.POST("/mthai_embed", admin.AddEmbedMThaiHandler)
	routerAuthorized.GET("/analytic", admin.AnalyticHandler)
	routerAuthorized.POST("/analytic", admin.AnalyticProcessHandler)
	routerAuthorized.GET("/flush", admin.FlushHandler)
	routerAuthorized.GET("/episode", admin.GetEpisodeHandler)
	routerAuthorized.POST("/episode", admin.SaveEpisodeHandler)
	router.NoRoute(notFoundHandler)
	router.Run(":" + port)
}
