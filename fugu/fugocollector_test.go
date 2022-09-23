package fugu

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestForPrivacyPage(t *testing.T) {
	url := "/what/privacy/ever/"
	assert.True(t, isPrivacyPage(url, "Whatever"))
}
