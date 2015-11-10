package utils

import (
	"encoding/json"
	"html"
	"html/template"
	"net/url"
	"reflect"
	"strings"
)

func add(a, b int) int {
	return a + b
}

func marshal(v interface{}) template.JS {
	a, _ := json.Marshal(v)
	return template.JS(a)
}

func templateHTML(a string) template.HTML {
	return template.HTML(a)
}

func escStr(a ...string) string {
	return html.EscapeString(strings.Join(a, "-"))
}

func urlEsc(a ...string) string {
	return url.QueryEscape(strings.Join(a, "-"))
}

func lastItem(x int, a interface{}) bool {
	return x == reflect.ValueOf(a).Len()-1
}
