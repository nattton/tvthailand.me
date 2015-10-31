package main

import (
	"fmt"
	"github.com/code-mobi/tvthailand.me/utils"
)

type CommandParam struct {
	Command string
	User    string
	Channel string
	Q       string
	Start   int
	Stop    int
}

func processCommand(param CommandParam) {
	db, err := utils.OpenDB()
	if err != nil {
		panic(err.Error())
	}
	defer db.Close()

	fmt.Println(param.Command)
	switch param.Command {
	case "botrun":

	}
}
