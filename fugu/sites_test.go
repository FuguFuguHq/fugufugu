package fugu

import (
	"github.com/stretchr/testify/assert"
	"net/url"
	"testing"
)

func TestEmptyFromExternals(t *testing.T) {
	externals := make(map[string]Privacy)
	sites := FromExternals(externals)
	assert.Equal(t, 0, len(sites))
}

func TestOneFromExternals(t *testing.T) {
	externals := make(map[string]Privacy)

	u := "https://www.googletagmanager.com/gtag/js"
	p := Privacy{
		Url:    &u,
		Typ:    "Script",
		Cookie: false,
	}
	externals[u] = p
	sites := FromExternals(externals)

	host, _ := url.Parse(u)
	assert.Equal(t, false, sites[host.Host].Cookie)
	assert.Equal(t, "www.googletagmanager.com", *sites[host.Host].Url)
	assert.Equal(t, 1, sites[host.Host].ScriptCount)
	assert.Equal(t, 0, sites[host.Host].CssCount)
	assert.Equal(t, 0, sites[host.Host].ImgCount)
}
