package controller

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestNewPager(t *testing.T) {
	ctx, _ := newContext("/")

	pager := NewPager(ctx, 10)
	assert.Equal(t, 10, pager.ItemsPerPage)
	assert.Equal(t, 1, pager.Page)
	assert.Equal(t, 0, pager.Items)
	assert.Equal(t, 0, pager.Pages)

	ctx, _ = newContext(fmt.Sprintf("/abc?%s=%d", PageQueryKey, 2))
	pager = NewPager(ctx, 10)
	assert.Equal(t, 2, pager.Page)

	ctx, _ = newContext(fmt.Sprintf("/abc?%s=%d", PageQueryKey, -2))
	pager = NewPager(ctx, 10)
	assert.Equal(t, 1, pager.Page)
}

func TestPager_SetItems(t *testing.T) {
	ctx, _ := newContext("/")
	pager := NewPager(ctx, 20)
	pager.SetItems(100)
	assert.Equal(t, 100, pager.Items)
	assert.Equal(t, 5, pager.Pages)
}
func TestPager_IsBeginning(t *testing.T) {
	ctx, _ := newContext("/")
	pager := NewPager(ctx, 20)
	pager.Pages = 10
	assert.True(t, pager.IsBeginning())
	pager.Page = 2
	assert.False(t, pager.IsBeginning())
	pager.Page = 1
	assert.True(t, pager.IsBeginning())
}
func TestPager_IsEnd(t *testing.T) {
	ctx, _ := newContext("/")
	pager := NewPager(ctx, 20)
	pager.Pages = 10
	assert.False(t, pager.IsEnd())
	pager.Page = 10
	assert.True(t, pager.IsEnd())
	pager.Page = 1
	assert.False(t, pager.IsEnd())
}
func TestPager_GetOffset(t *testing.T) {
	ctx, _ := newContext("/")
	pager := NewPager(ctx, 20)
	assert.Equal(t, 0, pager.GetOffset())
	pager.Page = 2
	assert.Equal(t, 20, pager.GetOffset())
	pager.Page = 3
	assert.Equal(t, 40, pager.GetOffset())
}
