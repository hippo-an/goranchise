package controller

import (
	"github.com/hippo-an/goranchise/context"
	"github.com/hippo-an/goranchise/msg"
	"github.com/labstack/echo/v4"
	echomw "github.com/labstack/echo/v4/middleware"
	"html/template"
	"net/http"
	"time"
)

type Page struct {
	AppName    string
	Title      string
	Context    echo.Context
	ToURL      func(name string, params ...interface{}) string
	Path       string
	URL        string
	Data       interface{}
	Layout     string
	PageName   string
	IsHome     bool
	IsAuth     bool
	StatusCode int
	MetaTags   struct {
		Description string
		Keywords    []string
	}
	Pager   Pager
	CSRF    string
	Headers map[string]string
	Cache   struct {
		Enabled    bool
		Expiration time.Duration
		Tags       []string
	}
	RequestId string
}

func NewPage(c echo.Context) Page {
	p := Page{
		Context:    c,
		ToURL:      c.Echo().Reverse,
		Path:       c.Request().URL.Path,
		URL:        c.Request().URL.String(),
		StatusCode: http.StatusOK,
		Pager:      NewPager(c, DefaultItemsPerPage),
		Headers:    make(map[string]string),
		RequestId:  c.Response().Header().Get(echo.HeaderXRequestID),
	}

	p.IsHome = p.Path == "/"

	if csrf := c.Get(echomw.DefaultCSRFConfig.ContextKey); csrf != nil {
		p.CSRF = csrf.(string)
	}

	// auth settings for page
	if u := c.Get(context.AuthenticatedUserKey); u != nil {
		p.IsAuth = true
	}

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
