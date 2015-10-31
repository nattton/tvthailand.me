package utils

import (
	"github.com/code-mobi/tvthailand.me/Godeps/_workspace/src/github.com/go-martini/martini"
	"github.com/code-mobi/tvthailand.me/Godeps/_workspace/src/github.com/jinzhu/gorm"
	"github.com/code-mobi/tvthailand.me/Godeps/_workspace/src/github.com/mssola/user_agent"
	"log"
	"os"
	"strings"
)

func OpenDB() (gorm.DB, error) {
	db, err := gorm.Open("mysql", os.Getenv("DATABASE_DSN"))
	if err != nil {
		log.Fatal(err)
	}
	db.LogMode(martini.Env != "production")
	return db, err
}

func IsMobile(userAgent string) bool {
	ua := user_agent.New(userAgent)
	isiPad := strings.Contains(userAgent, "iPad")
	return ua.Mobile() && !isiPad
}
