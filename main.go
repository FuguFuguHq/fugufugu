package main

import (
	"flag"
	"fmt"
	"fugufugu/fugu"
	"github.com/jedib0t/go-pretty/v6/table"
	"os"
	"sort"
)

var Products map[string]fugu.Product

type Report struct {
	Url         string
	ImageCount  int
	CssCount    int
	ScriptCount int
	PageCount   uint64

	Externals    []ReportExternal
	PrivacyPages []ReportPrivacyPage
}

type ReportPrivacyPage struct {
	Url   string
	Title string
}

type ReportExternal struct {
	Url     string
	Company string
	Product string
	Country string
	Script  bool
	Image   bool
	Css     bool
	Cookie  bool
}

func ReportFromScanner(checkUrl string, privacies map[string]fugu.SitePrivacy, scanner fugu.Scanner) Report {
	v := make([]fugu.SitePrivacy, 0, len(privacies))

	for _, value := range privacies {
		value.Rank = fugu.Rank(value)
		v = append(v, value)
	}

	sort.Slice(v, func(i, j int) bool {
		return v[i].Rank > v[j].Rank
	})

	scripts := 0
	images := 0
	css := 0

	r := Report{
		Url:          checkUrl,
		PageCount:    *scanner.Pages,
		Externals:    make([]ReportExternal, 0),
		PrivacyPages: make([]ReportPrivacyPage, 0),
	}

	for _, p := range v {
		scripts += p.ScriptCount
		images += p.ImgCount
		css += p.CssCount
		country := ""
		company := ""
		product := ""
		if p.Product != nil && p.Product.Company != nil {
			product = p.Product.Name
			country = p.Product.Company.Country
			company = p.Product.Company.Name
		}

		rl := ReportExternal{
			Url:     *p.Url,
			Company: company,
			Product: product,
			Country: country,
			Script:  p.ScriptCount > 0,
			Image:   p.ImgCount > 0,
			Css:     p.CssCount > 0,
			Cookie:  p.Cookie,
		}
		r.Externals = append(r.Externals, rl)
	}
	r.ImageCount = images
	r.ScriptCount = scripts
	r.CssCount = css

	if len(*scanner.PrivacyPages) > 0 {
		for _, page := range *scanner.PrivacyPages {
			rpp := ReportPrivacyPage{
				Url:   *page.URL,
				Title: *page.Title,
			}
			r.PrivacyPages = append(r.PrivacyPages, rpp)
		}
	}
	return r
}

func main() {
	Products = fugu.Products()

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

	privacies := fugu.FromExternals(Products, externals)

	r := ReportFromScanner(checkUrl, privacies, scanner)

	if len(privacies) == 0 {
		fmt.Printf("Summary %s: %d pages\n", checkUrl, *scanner.Pages)
		fmt.Println("Cool, no external resources found!")
	} else {

		// PRINT REPORT

		if len(r.PrivacyPages) > 0 {
			p := table.NewWriter()
			p.SetOutputMirror(os.Stdout)
			p.AppendHeader(table.Row{"Privacy Page", "Title"})
			for _, page := range r.PrivacyPages {
				p.AppendRows([]table.Row{
					{page.Url, page.Title},
				})
			}
			p.Render()
		}

		t := table.NewWriter()
		t.SetOutputMirror(os.Stdout)
		if checkForCookie {
			t.AppendHeader(table.Row{"Site", "Company", "Product", "Country", "Script", "Image", "Css", "Cookie"})
		} else {
			t.AppendHeader(table.Row{"Site", "Company", "Product", "Country", "Script", "Image", "Css"})
		}

		for _, e := range r.Externals {
			script := ""
			image := ""
			css := ""
			cookie := ""
			if e.Image {
				image = "Yes"
			}
			if e.Script {
				script = "Yes"
			}
			if e.Css {
				css = "Yes"
			}
			if e.Cookie {
				cookie = "Yes"
			}
			if checkForCookie {
				t.AppendRows([]table.Row{
					{e.Url, e.Company, e.Product, e.Country, script, image, css, cookie},
				})
			} else {
				t.AppendRows([]table.Row{
					{e.Url, e.Company, e.Product, e.Country, script, image, css},
				})
			}
		}
		t.Render()

		fmt.Printf("Summary %s: %d pages - External resources: %d scripts | %d images | %d css\n", r.Url, r.PageCount, r.ScriptCount,
			r.ImageCount, r.CssCount)
	}
}
