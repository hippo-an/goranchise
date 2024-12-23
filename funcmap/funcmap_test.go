package funcmap

import (
	"fmt"
	"github.com/hippo-an/goranchise/config"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestHasField(t *testing.T) {
	type example struct {
		name string
	}
	var e example
	assert.True(t, HasField(e, "name"))
	assert.False(t, HasField(e, "asdfasdf"))
}

func TestLink(t *testing.T) {
	link := string(Link("/abc", "Text", "/abc"))
	expected := `<a class="is-active" href="/abc">Text</a>`
	assert.Equal(t, expected, link)
	link = string(Link("/abc", "Text", "/abc", "first", "second"))
	expected = `<a class="first second is-active" href="/abc">Text</a>`
	assert.Equal(t, expected, link)
	link = string(Link("/abc", "Text", "/def"))
	expected = `<a class="" href="/abc">Text</a>`
	assert.Equal(t, expected, link)
}

func TestGetFuncMap(t *testing.T) {
	file := File("test.png")
	expected := fmt.Sprintf("/%s/test.png?v=%s", config.StaticDir, CacheKey)
	assert.Equal(t, expected, file)
}
