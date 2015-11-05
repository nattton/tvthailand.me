package main

import (
	"fmt"
	"github.com/code-mobi/tvthailand.me/data"
	"github.com/code-mobi/tvthailand.me/utils"
)

type CommandParam struct {
	Command  string
	Channel  string
	Playlist string
	Query    string
	Start    int
	Stop     int
}

func processCommand(cmd CommandParam) {
	db, err := utils.OpenDB()
	if err != nil {
		panic(err.Error())
	}
	defer db.Close()

	fmt.Println(cmd.Command)
	switch cmd.Command {
	case "runbotch":
		if commandParam.Channel != "" {
			data.RunBotChannel(&db, commandParam.Channel, commandParam.Query)
		} else {
			data.RunBotChannels(&db)
		}
	case "runbotpl":
		if commandParam.Playlist != "" {
			data.RunBotPlaylist(&db, commandParam.Playlist)
		} else {
			data.RunBotPlaylists(&db)
		}
	case "updateuser":
		data.UpdateUserChannel(&db)
	case "checkuser":
		data.CheckActiveUser(&db)
	case "migrate_botvideo":
		data.MigrateUsernameToChannelID(&db)
	}
}
