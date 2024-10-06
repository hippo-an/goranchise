package funcmap

import (
	"fmt"
	"github.com/Masterminds/sprig"
	"github.com/labstack/gommon/random"
	"html/template"
	"reflect"
	"strings"
)

var CacheKey = random.String(10)

func GetFuncMap() template.FuncMap {
	funcMap := sprig.FuncMap()

	f := template.FuncMap{
		"hasField": HasField,
		"file":     File,
		"link":     Link,
	}

	for k, v := range f {
		funcMap[k] = v
	}

	return funcMap
}

func HasField(v interface{}, name string) bool {
	rv := reflect.ValueOf(v)

	if rv.Kind() == reflect.Ptr {
		rv = rv.Elem()
	}

	if rv.Kind() != reflect.Struct {
		return false
	}

	return rv.FieldByName(name).IsValid()
}

func File(filepath string) string {
	return fmt.Sprintf("%s?v=%s", filepath, CacheKey)
}

func Link(url, text, currentPath string, classes ...string) template.HTML {
	if currentPath == url {
		classes = append(classes, "active")
	}

	html := fmt.Sprintf(`<a class="%s" href="%s">%s</a>`, strings.Join(classes, " "), url, text)

	return template.HTML(html)
}
