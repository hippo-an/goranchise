package controller

import (
	"github.com/labstack/echo/v4"
	"math"
	"strconv"
)

const (
	DefaultItemsPerPage = 20
	PageQueryKey        = "page"
)

type Pager struct {
	// Items stores the total amount of items in the result
	Items int

	// ItemsPerPage stores the amount of items to display per page
	ItemsPerPage int

	// Page stores the current page number
	Page int

	// Pages store the total amount of pages in the result
	Pages int
}

func NewPager(e echo.Context, itemPerPage int) Pager {

	if itemPerPage <= 0 {
		itemPerPage = DefaultItemsPerPage
	}

	p := Pager{
		ItemsPerPage: itemPerPage,
		Page:         1,
	}

	if page := e.QueryParam(PageQueryKey); page != "" {
		if pageInt, err := strconv.Atoi(page); err == nil {
			if pageInt > 0 {
				p.Page = pageInt
			}
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
