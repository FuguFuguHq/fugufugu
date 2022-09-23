package fugu

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestForPrivacyPage(t *testing.T) {
	url1 := "/what/privacy/ever/"
	assert.True(t, isPrivacyPage(url1, "Whatever"))
	url2 := "/what/ever/"
	assert.True(t, isPrivacyPage(url2, "Datenschutz"))
}
