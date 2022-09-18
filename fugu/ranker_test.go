package fugu

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestRanker(t *testing.T) {
	url := "https://google.com"
	s := SitePrivacy{
		Url:    &url,
		Cookie: false,
	}
	assert.Equal(t, 0, Rank(s))
}

func TestRankerImages(t *testing.T) {
	url := "https://google.com"
	s := SitePrivacy{
		Url:      &url,
		Cookie:   false,
		ImgCount: 1,
	}
	assert.Equal(t, 1, Rank(s))

	s = SitePrivacy{
		Url:      &url,
		Cookie:   false,
		ImgCount: 11,
	}
	assert.Equal(t, 2, Rank(s))
}

func TestRankerScripts(t *testing.T) {
	url := "https://google.com"
	s := SitePrivacy{
		Url:         &url,
		Cookie:      false,
		ScriptCount: 1,
	}
	assert.Equal(t, 100, Rank(s))

	s = SitePrivacy{
		Url:         &url,
		Cookie:      false,
		ScriptCount: 11,
	}
	assert.Equal(t, 200, Rank(s))
}

func TestRankerCookie(t *testing.T) {
	url := "https://google.com"
	s := SitePrivacy{
		Url:    &url,
		Cookie: true,
	}
	assert.Equal(t, 1000, Rank(s))
}
