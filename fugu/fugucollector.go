package fugu

import (
	"github.com/gocolly/colly/v2"
	"io"
	"log"
	"net/http"
	"net/url"
	"strings"
	"time"
)

func NewCollector(checkUrl string, externals map[string]Privacy) *colly.Collector {

	u, err := url.Parse(checkUrl)
	if err != nil {
		panic(err)
	}

	c := colly.NewCollector(
		colly.AllowedDomains(u.Hostname()),
		colly.MaxDepth(1),
	)

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
					log.Println(err)
				} else {
					externals[link] = Privacy{
						Typ:    "Image",
						Cookie: len(resp.Header.Get("Set-Cookie")) > 0,
					}
				}
			}
		}
	})

	c.OnHTML("script[src]", func(e *colly.HTMLElement) {
		link := e.Attr("src")
		if _, ok := externals[link]; !ok {
			if !strings.HasPrefix(link, checkUrl) && strings.HasPrefix(link, "https://") {
				client := http.Client{
					Timeout: 5 * time.Second,
				}
				resp, err := client.Get(link)
				if err != nil {
					log.Println(err)
				} else {
					externals[link] = Privacy{
						Typ:    "Script",
						Cookie: len(resp.Header.Get("Set-Cookie")) > 0,
					}
				}
			}
		}
	})

	c.OnHTML("link[href]", func(e *colly.HTMLElement) {
		link := e.Attr("href")
		rel := e.Attr("rel")
		if len(rel) > 0 && rel == "stylesheet" {
			cssLink := ""
			if !strings.HasPrefix(link, "http") {
				cssLink = checkUrl + "/" + link
			} else {
				cssLink = link
			}
			client := http.Client{
				Timeout: 5 * time.Second,
			}
			resp, err := client.Get(cssLink)
			if err != nil {
				log.Println(err)
			} else {
				b, err := io.ReadAll(resp.Body)
				if err != nil {
					log.Println(err)
				} else {
					imports := ImportsFromCss(string(b))
					for _, i := range imports {
						externals[i] = Privacy{
							Typ: "Css",
							// @TODO check @import files for cookies
							Cookie: false,
						}
					}
				}
			}

			if !strings.HasPrefix(link, checkUrl) && strings.HasPrefix(link, "https://") {
				externals[link] = Privacy{
					Typ:    "Css",
					Cookie: len(resp.Header.Get("Set-Cookie")) > 0,
				}
			}
		}
	})

	// Before making a request print "Visiting ..."
	c.OnRequest(func(r *colly.Request) {
	})

	return c
}
