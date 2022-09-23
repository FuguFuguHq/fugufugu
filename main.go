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
	var checkForCookie bool
	var maxPages uint64
	flag.StringVar(&checkUrl, "url", "", "url to check")
	flag.Uint64Var(&maxPages, "max", 10000, "max pages to check")
	flag.BoolVar(&verbose, "verbose", false, "verbose mode")
	flag.BoolVar(&checkForCookie, "cookie", false, "check for cookies")
	flag.Parse()

	if len(checkUrl) == 0 {
		fmt.Println("No url specified.")
		os.Exit(-1)
	}

	externals := make(map[string]fugu.Privacy)
	scanner := fugu.NewCollector(maxPages, checkForCookie, checkUrl, externals, verbose)

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
		fmt.Printf("Summary %s: %d pages\n", checkUrl, *scanner.Pages)
		fmt.Println("Cool, no external resources found!")
	} else {
		t := table.NewWriter()
		t.SetOutputMirror(os.Stdout)
		if checkForCookie {
			t.AppendHeader(table.Row{"Site", "Company", "Country", "Script", "Image", "Css", "Cookie"})
		} else {
			t.AppendHeader(table.Row{"Site", "Company", "Country", "Script", "Image", "Css"})
		}

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
			if checkForCookie {
				t.AppendRows([]table.Row{
					{*p.Url, company, country, s, i, css, c},
				})
			} else {
				t.AppendRows([]table.Row{
					{*p.Url, company, country, s, i, css},
				})
			}
		}
		if len(*scanner.PrivacyPages) > 0 {
			p := table.NewWriter()
			p.SetOutputMirror(os.Stdout)
			p.AppendHeader(table.Row{"Privacy Page", "Title"})
			for _, page := range *scanner.PrivacyPages {
				p.AppendRows([]table.Row{
					{*page.URL, *page.Title},
				})
			}
			p.Render()
		}

		t.Render()
		fmt.Printf("Summary %s: %d pages | %d scripts | %d images | %d css\n", checkUrl, *scanner.Pages, scripts, images, css)
	}
}
