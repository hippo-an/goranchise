package handlers

import (
	"errors"
	"fmt"
	"github.com/hippo-an/goranchise/controller"
	"github.com/labstack/echo/v4"
	"net/http"
)

type Error struct {
	controller.Controller
}

func (e *Error) Handler(err error, ctx echo.Context) {
	if ctx.Response().Committed {
		return
	}

	code := http.StatusInternalServerError
	var he *echo.HTTPError
	if errors.As(err, &he) {
		code = he.Code
	}
	if code >= 500 {
		ctx.Logger().Error(err)
	} else {
		ctx.Logger().Info(err)
	}
	p := controller.NewPage(ctx)
	p.Layout = "main"
	p.PageName = "error"
	p.Title = http.StatusText(code)

	p.PageName = fmt.Sprintf("errors/%d", code)
	p.StatusCode = code
	p.Data = err.Error()
	if err = e.RenderPage(ctx, p); err != nil {
		ctx.Logger().Error(err)
	}
}
