package fugu

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestCss(t *testing.T) {
	css := "   @import url('https://fonts.googleapis.com/css?family=Muli')    "
	imports := ImportsFromCss(css, false)

	assert.Equal(t, 1, len(imports))
	assert.Equal(t, "https://fonts.googleapis.com/css?family=Muli", imports[0])
}
