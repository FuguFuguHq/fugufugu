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

func main() {
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
		if !strings.HasPrefix(link, checkUrl) && strings.HasPrefix(link, "https://") {
			resp, err := http.Get(link)
			if err != nil {
				log.Fatalln(err)
			}
			externals[link] = fugu.Privacy{
				Typ:    "Link",
				Cookie: len(resp.Header.Get("Set-Cookie")) > 0,
			}
		}
	})

	// Before making a request print "Visiting ..."
	c.OnRequest(func(r *colly.Request) {
	})

	c.Visit(checkUrl)

	privacies := fugu.FromExternals(externals)

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
	t.AppendHeader(table.Row{"Site", "Script", "Image", "Cookie"})

	for _, p := range v {
		c, s, i := "", "", ""
		if p.Cookie {
			c = "Yes"
		}
		if p.ScriptCount > 0 {
			s = "Yes"
		}
		if p.ImgCount > 0 {
			i = "Yes"
		}
		t.AppendRows([]table.Row{
			{*p.Url, s, i, c},
		})
	}
	t.Render()
}
