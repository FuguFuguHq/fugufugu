package fugu

import (
	"github.com/stretchr/testify/assert"
	"net/url"
	"testing"
)

func TestEmptyFromExternals(t *testing.T) {
	products := Products()
	externals := make(map[string]Privacy)
	sites := FromExternals(products, externals)
	assert.Equal(t, 0, len(sites))
}

func TestOneFromExternals(t *testing.T) {
	products := Products()
	externals := make(map[string]Privacy)

	u := "https://www.googletagmanager.com/gtag/js"
	p := Privacy{
		Url:    &u,
		Typ:    "Script",
		Cookie: false,
	}
	externals[u] = p
	sites := FromExternals(products, externals)

	host, _ := url.Parse(u)
	assert.Equal(t, false, sites[host.Host].Cookie)
	assert.Equal(t, "www.googletagmanager.com", *sites[host.Host].Url)
	assert.Equal(t, 1, sites[host.Host].ScriptCount)
	assert.Equal(t, 0, sites[host.Host].CssCount)
	assert.Equal(t, 0, sites[host.Host].ImgCount)
}

func TestThreeFromExternals(t *testing.T) {
	externals := make(map[string]Privacy)
	products := Products()

	u1 := "https://www.googletagmanager.com/1"
	p1 := Privacy{
		Url:    &u1,
		Typ:    "Script",
		Cookie: false,
	}
	u2 := "https://www.googletagmanager.com/2"
	p2 := Privacy{
		Url:    &u2,
		Typ:    "Image",
		Cookie: false,
	}
	u3 := "https://www.googletagmanager.com/3"
	p3 := Privacy{
		Url:    &u3,
		Typ:    "Css",
		Cookie: false,
	}
	externals[u1] = p1
	externals[u2] = p2
	externals[u3] = p3
	sites := FromExternals(products, externals)

	host, _ := url.Parse(u1)
	assert.Equal(t, false, sites[host.Host].Cookie)
	assert.Equal(t, "www.googletagmanager.com", *sites[host.Host].Url)
	assert.Equal(t, 1, sites[host.Host].ScriptCount)
	assert.Equal(t, 1, sites[host.Host].CssCount)
	assert.Equal(t, 1, sites[host.Host].ImgCount)
}
