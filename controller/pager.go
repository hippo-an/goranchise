package controller

import (
	"github.com/labstack/echo/v4"
	"math"
	"strconv"
)

const (
	DefaultItemsPerPage = 20
)

type Pager struct {
	Items        int
	Page         int
	ItemsPerPage int
	Pages        int
}

func NewPager(e echo.Context) Pager {
	p := Pager{
		ItemsPerPage: DefaultItemsPerPage,
		Page:         1,
	}

	if page := e.QueryParam("page"); page != "" {
		if pageInt, err := strconv.Atoi(page); err == nil {
			p.Page = pageInt
		}
	}

	return p
}

func (p *Pager) SetItems(items int) {
	p.Items = items
	p.Pages = int(math.Ceil(float64(items) / float64(p.ItemsPerPage)))
	if p.Page > p.Pages {
		p.Page = p.Pages
	}
}
func (p *Pager) IsBeginning() bool {
	return p.Page == 1
}

func (p *Pager) IsEnd() bool {
	return p.Page >= p.Pages
}

func (p *Pager) GetOffset() int {
	if p.Page == 0 {
		p.Page = 1
	}
	return (p.Page - 1) * p.ItemsPerPage
}
