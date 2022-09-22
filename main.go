package main

import (
	"flag"
	"fmt"
	"fugufugu/fugu"
	"github.com/jedib0t/go-pretty/v6/table"
	"os"
	"sort"
)

var Companies map[string]fugu.Company

func main() {
	Companies = fugu.Companies()

	var checkUrl string
	var verbose bool
	flag.StringVar(&checkUrl, "url", "", "url to check")
	flag.BoolVar(&verbose, "verbose", false, "verbose mode")
	flag.Parse()

	if len(checkUrl) == 0 {
		fmt.Println("No url specified.")
		os.Exit(-1)
	}

	externals := make(map[string]fugu.Privacy)
	scanner := fugu.NewCollector(checkUrl, externals, verbose)

	scanner.Collector.Visit(checkUrl)

	privacies := fugu.FromExternals(Companies, externals)

	v := make([]fugu.SitePrivacy, 0, len(privacies))

	for _, value := range privacies {
		value.Rank = fugu.Rank(value)
		v = append(v, value)
	}

	sort.Slice(v, func(i, j int) bool {
		return v[i].Rank > v[j].Rank
	})

	if len(v) == 0 {
		fmt.Println("Cool, no external resources found!")
	} else {
		t := table.NewWriter()
		t.SetOutputMirror(os.Stdout)
		t.AppendHeader(table.Row{"Site", "Company", "Country", "Script", "Image", "Css", "Cookie"})

		scripts := 0
		images := 0
		css := 0

		for _, p := range v {
			scripts += p.ScriptCount
			images += p.ImgCount
			css += p.CssCount
			country := ""
			company := ""
			if p.Company != nil {
				country = p.Company.Country
				company = p.Company.Name
			}
			c, s, i, css := "", "", "", ""
			if p.Cookie {
				c = "Yes"
			}
			if p.CssCount > 0 {
				css = "Yes"
			}
			if p.ScriptCount > 0 {
				s = "Yes"
			}
			if p.ImgCount > 0 {
				i = "Yes"
			}
			t.AppendRows([]table.Row{
				{*p.Url, company, country, s, i, css, c},
			})
		}
		fmt.Printf("Summary %s: %d scripts | %d images | %d css\n", checkUrl, scripts, images, css)
		t.Render()
	}
}
