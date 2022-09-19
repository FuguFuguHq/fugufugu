package main

import (
	"flag"
	"fmt"
	"fugufugu/fugu"
	"github.com/gocolly/colly/v2"
	"github.com/jedib0t/go-pretty/v6/table"
	"log"
	"net/http"
	"net/url"
	"os"
	"sort"
	"strings"
)

var Companies map[string]fugu.Company

func main() {
	Companies = fugu.Companies()

	var checkUrl string
	flag.StringVar(&checkUrl, "url", "", "url to check")
	flag.Parse()

	if len(checkUrl) == 0 {
		fmt.Println("No url specified.")
		os.Exit(-1)
	}

	u, err := url.Parse(checkUrl)
	if err != nil {
		panic(err)
	}
	c := colly.NewCollector(
		colly.AllowedDomains(u.Host),
	)

	externals := make(map[string]fugu.Privacy)

	c.OnRequest(func(r *colly.Request) {
	})

	c.OnResponse(func(r *colly.Response) {
		//	fmt.Printf("Link found: %v\n", r.Headers)
	})

	// On every a element which has href attribute call callback
	c.OnHTML("a[href]", func(e *colly.HTMLElement) {
		link := e.Attr("href")
		// Print link
		// fmt.Printf("Link found: %q -> %s\n", e.Text, link)
		// Visit link found on page
		// Only those links are visited which are in AllowedDomains
		c.Visit(e.Request.AbsoluteURL(link))
	})

	c.OnHTML("img[src]", func(e *colly.HTMLElement) {
		link := e.Attr("src")
		if _, ok := externals[link]; !ok {
			if !strings.HasPrefix(link, checkUrl) && strings.HasPrefix(link, "https://") {
				resp, err := http.Get(link)
				if err != nil {
					log.Fatalln(err)
				}

				externals[link] = fugu.Privacy{
					Typ:    "Image",
					Cookie: len(resp.Header.Get("Set-Cookie")) > 0,
				}
			}
		}
	})

	c.OnHTML("script[src]", func(e *colly.HTMLElement) {
		link := e.Attr("src")
		if _, ok := externals[link]; !ok {
			if !strings.HasPrefix(link, checkUrl) && strings.HasPrefix(link, "https://") {
				resp, err := http.Get(link)
				if err != nil {
					log.Fatalln(err)
				}
				externals[link] = fugu.Privacy{
					Typ:    "Script",
					Cookie: len(resp.Header.Get("Set-Cookie")) > 0,
				}

			}
		}
	})

	c.OnHTML("link[href]", func(e *colly.HTMLElement) {
		link := e.Attr("href")
		rel := e.Attr("rel")
		if len(rel) > 0 {
			if !strings.HasPrefix(link, checkUrl) && strings.HasPrefix(link, "https://") {
				resp, err := http.Get(link)
				if err != nil {
					log.Fatalln(err)
				}
				externals[link] = fugu.Privacy{
					Typ:    "Css",
					Cookie: len(resp.Header.Get("Set-Cookie")) > 0,
				}
			}
		}
	})

	// Before making a request print "Visiting ..."
	c.OnRequest(func(r *colly.Request) {
	})

	c.Visit(checkUrl)

	privacies := fugu.FromExternals(Companies, externals)

	v := make([]fugu.SitePrivacy, 0, len(privacies))

	for _, value := range privacies {
		value.Rank = fugu.Rank(value)
		v = append(v, value)
	}

	sort.Slice(v, func(i, j int) bool {
		return v[i].Rank > v[j].Rank
	})

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
