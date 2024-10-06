package controllers

import (
	"github.com/hippo-an/goranchise/msg"
	"github.com/labstack/echo/v4"
	"html/template"
	"net/http"
)

type Page struct {
	AppName    string
	Title      string
	Context    echo.Context
	Reverse    func(name string, params ...interface{}) string
	Path       string
	Data       interface{}
	Layout     string
	Name       string
	IsHome     bool
	IsAuth     bool
	StatusCode int
	MetaTags   struct {
		Description string
		Keywords    []string
	}
}

func NewPage(c echo.Context) Page {
	p := Page{
		Context:    c,
		Reverse:    c.Echo().Reverse,
		Path:       c.Request().URL.Path,
		StatusCode: http.StatusOK,
	}

	p.IsAuth = p.Path == "/"
	return p
}

func (p Page) SetMessage(typ msg.Type, value string) {
	msg.Set(p.Context, typ, value)
}
func (p Page) GetMessages(typ msg.Type) []template.HTML {
	strs := msg.Get(p.Context, typ)
	ret := make([]template.HTML, len(strs))
	for k, v := range strs {
		ret[k] = template.HTML(v)
	}
	return ret
}
