package utils

import (
	"bytes"
	"fmt"
	"github.com/code-mobi/tvthailand.me/Godeps/_workspace/src/github.com/gin-gonic/gin"
	_ "github.com/code-mobi/tvthailand.me/Godeps/_workspace/src/github.com/go-sql-driver/mysql"
	"github.com/code-mobi/tvthailand.me/Godeps/_workspace/src/github.com/jinzhu/gorm"
	"github.com/code-mobi/tvthailand.me/Godeps/_workspace/src/github.com/mssola/user_agent"
	"github.com/code-mobi/tvthailand.me/Godeps/_workspace/src/gopkg.in/redis.v3"
	"html/template"
	"log"
	"net/http"
	"os"
	"strings"
	"time"
)

func OpenDB() (gorm.DB, error) {
	db, err := gorm.Open("mysql", os.Getenv("DATABASE_DSN"))
	if err != nil {
		log.Fatal(err)
	}
	db.LogMode(gin.Mode() == "debug")
	return db, err
}

func OpenRedis() *redis.Client {
	client := redis.NewClient(&redis.Options{
		Addr:     os.Getenv("REDIS_HOST") + ":6379",
		Password: "",
		DB:       0,
	})
	return client
}

func IsMobile(userAgent string) bool {
	ua := user_agent.New(userAgent)
	isiPad := strings.Contains(userAgent, "iPad")
	return ua.Mobile() && !isiPad
}

// parse HTML templates
// pass in a list of file names, and get a template
func ParseTemplateFiles(filenames ...string) (t *template.Template) {
	var files []string
	t = template.New("layout")
	for _, file := range filenames {
		files = append(files, fmt.Sprintf("templates/%s.tmpl", file))
	}
	t = template.Must(t.ParseFiles(files...))
	return
}

func GenerateHTML(writer http.ResponseWriter, renderData interface{}, filenames ...string) {
	var files []string
	for _, file := range filenames {
		files = append(files, fmt.Sprintf("templates/%s.tmpl", file))
	}

	var funcMaps = []template.FuncMap{
		{"add": add},
		{"last": lastItem},
		{"escStr": escStr},
		{"urlEsc": urlEsc},
		{"marshal": marshal},
		{"html": templateHTML},
	}
	templates := template.New("")
	for i := range funcMaps {
		templates = templates.Funcs(funcMaps[i])
	}
	templates = template.Must(templates.ParseFiles(files...))

	var doc bytes.Buffer
	templates.ExecuteTemplate(&doc, "layout", renderData)
	if renderData != nil {
		CachedKey := renderData.(map[string]interface{})["CachedKey"]
		if CachedKey != nil {
			redisClient := OpenRedis()
			redisClient.Set(CachedKey.(string), doc.String(), 5*time.Minute)
		}
	}

	fmt.Fprint(writer, doc.String())
}
