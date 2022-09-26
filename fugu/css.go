package fugu

import (
	"fmt"
	"regexp"
)

// find @import url('https://fonts.googleapis.com/css?family=Muli&display=swap');
func ImportsFromCss(css string, verbose bool) []string {
	imports := make([]string, 0)

	regex := regexp.MustCompile("@import[ ]+url\\('(.*)'\\)")
	res := regex.FindAllStringSubmatch(css, -1)
	for i := range res {
		imports = append(imports, res[i][1])
		if verbose {
			fmt.Println("Importing CSS: " + res[i][1])
		}
	}
	return imports
}
