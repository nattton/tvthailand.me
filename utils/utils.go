package utils

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/code-mobi/tvthailand.me/Godeps/_workspace/src/github.com/gin-gonic/gin"
	_ "github.com/code-mobi/tvthailand.me/Godeps/_workspace/src/github.com/go-sql-driver/mysql"
	"github.com/code-mobi/tvthailand.me/Godeps/_workspace/src/github.com/jinzhu/gorm"
	"github.com/code-mobi/tvthailand.me/Godeps/_workspace/src/github.com/mssola/user_agent"
	"github.com/code-mobi/tvthailand.me/Godeps/_workspace/src/gopkg.in/redis.v3"
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
	return ua.Mobile()
}

func IsMobileNotPad(userAgent string) bool {
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

func GenerateHTML(writer http.ResponseWriter, data interface{}, filenames ...string) {
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
		{"metaDescription": metaDescription},
		{"fmtTime": fmtTime},
	}
	templates := template.New("")
	for i := range funcMaps {
		templates = templates.Funcs(funcMaps[i])
	}
	templates = template.Must(templates.ParseFiles(files...))

	templates.ExecuteTemplate(writer, "layout", data)
}

// DeleteHashCached Delete Hash Cached by CachedKey
func DeleteHashCached(cachedKey string) {
	redisClient := OpenRedis()
	resultList, err := redisClient.HKeys(cachedKey).Result()
	if err == redis.Nil {
		fmt.Println(err)
	}
	fmt.Println(resultList)
	redisClient.Del(resultList...)
	redisClient.Del(cachedKey)
}
