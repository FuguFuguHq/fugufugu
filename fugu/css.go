package fugu

import (
	"github.com/davecgh/go-spew/spew"
	"regexp"
)

// find @import url('https://fonts.googleapis.com/css?family=Muli&display=swap');
func ImportsFromCss(css string) []string {
	imports := make([]string, 0)

	regex := regexp.MustCompile("@import[ ]+url\\('(.*)'\\)")
	res := regex.FindAllStringSubmatch(css, -1)
	spew.Dump(css)
	for i := range res {
		imports = append(imports, res[i][1])
	}
	return imports
}
