package controllers

import (
	"errors"
	"fmt"
	"github.com/labstack/echo/v4"
	"net/http"
)

type Error struct {
	Controller
}

func (e *Error) Handler(err error, ctx echo.Context) {
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
	p := NewPage(ctx)
	p.Layout = "main"
	p.Title = http.StatusText(code)

	p.Name = fmt.Sprintf("errors/%d", code)
	p.StatusCode = code
	if err = e.RenderPage(ctx, p); err != nil {
		ctx.Logger().Error(err)
	}
}
