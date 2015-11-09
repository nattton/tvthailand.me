package utils

import (
	"encoding/json"
	"fmt"
	"github.com/code-mobi/tvthailand.me/Godeps/_workspace/src/github.com/gin-gonic/gin"
	"github.com/code-mobi/tvthailand.me/Godeps/_workspace/src/github.com/jinzhu/gorm"
	"github.com/code-mobi/tvthailand.me/Godeps/_workspace/src/github.com/mssola/user_agent"
	_ "github.com/go-sql-driver/mysql"
	"html"
	"html/template"
	"log"
	"net/http"
	"net/url"
	"os"
	"reflect"
	"strings"
)

type Configuration struct {
	Port         string
	ReadTimeout  int64
	WriteTimeout int64
	Static       string
}

func LoadConfig() (config Configuration) {
	file, err := os.Open("config.json")
	if err != nil {
		log.Fatalln("Cannot open config file", err)
	}
	decoder := json.NewDecoder(file)
	err = decoder.Decode(&config)
	if err != nil {
		log.Fatalln("Cannot get configuration from file", err)
	}
	return
}

func OpenDB() (gorm.DB, error) {
	db, err := gorm.Open("mysql", os.Getenv("DATABASE_DSN"))
	if err != nil {
		log.Fatal(err)
	}
	db.LogMode(gin.Mode() == "debug")
	return db, err
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

func GenerateHTML(writer http.ResponseWriter, data interface{}, filenames ...string) {
	var files []string
	for _, file := range filenames {
		files = append(files, fmt.Sprintf("templates/%s.tmpl", file))
	}

	var funcMaps = []template.FuncMap{
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
			"marshal": func(v interface{}) template.JS {
				a, _ := json.Marshal(v)
				return template.JS(a)
			},
		},
		{
			"html": func(a string) template.HTML {
				return template.HTML(a)
			},
		},
	}
	templates := template.New("")
	for i := range funcMaps {
		templates = templates.Funcs(funcMaps[i])
	}
	templates = template.Must(templates.ParseFiles(files...))

	templates.ExecuteTemplate(writer, "layout", data)
}
