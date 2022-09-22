package fugu

import (
	"fmt"
	"github.com/gocolly/colly/v2"
	"io"
	"log"
	"net/http"
	"net/url"
	"strings"
	"sync/atomic"
	"time"
)

type Scanner struct {
	Collector *colly.Collector
	Pages     *uint64
}

func NewCollector(maxPages uint64, checkForCookie bool, checkUrl string, externals map[string]Privacy, verbose bool) Scanner {

	checkedInternalCss := make(map[string]bool)

	u, err := url.Parse(checkUrl)
	if err != nil {
		panic(err)
	}

	c := colly.NewCollector(
		colly.AllowedDomains(u.Hostname()),
		colly.UserAgent("FuguFugu"),
		colly.MaxDepth(1),
	)

	scanner := Scanner{
		Collector: c,
		Pages:     new(uint64),
	}

	c.OnRequest(func(r *colly.Request) {
		if verbose {
			fmt.Println("Checking " + r.URL.String())
		}
	})

	c.OnResponse(func(r *colly.Response) {
		atomic.AddUint64(scanner.Pages, 1)
		//	fmt.Printf("Link found: %v\n", r.Headers)
	})

	// On every a element which has href attribute call callback
	c.OnHTML("a[href]", func(e *colly.HTMLElement) {
		link := e.Attr("href")
		// Print link
		// fmt.Printf("Link found: %q -> %s\n", e.Text, link)
		// Visit link found on page
		// Only those links are visited which are in AllowedDomains
		if *scanner.Pages < maxPages {
			c.Visit(e.Request.AbsoluteURL(link))
		}
	})

	c.OnHTML("img[src]", func(e *colly.HTMLElement) {
		link := e.Attr("src")
		if _, ok := externals[link]; !ok {
			if !strings.HasPrefix(link, checkUrl) && strings.HasPrefix(link, "https://") {
				hasCookie := false
				if checkForCookie {
					client := http.Client{
						Timeout: 5 * time.Second,
					}
					resp, err := client.Get(link)
					if verbose {
						fmt.Println("Downloading " + link)
					}
					if err != nil {
						log.Println(err)
					} else {
						hasCookie = len(resp.Header.Get("Set-Cookie")) > 0
					}
					externals[link] = Privacy{
						Typ:    "Image",
						Cookie: hasCookie,
					}
				}
			}
		}
	})

	c.OnHTML("script[src]", func(e *colly.HTMLElement) {
		link := e.Attr("src")
		if _, ok := externals[link]; !ok {
			if !strings.HasPrefix(link, checkUrl) && strings.HasPrefix(link, "https://") {
				hasCookie := false
				if checkForCookie {
					client := http.Client{
						Timeout: 5 * time.Second,
					}
					resp, err := client.Get(link)
					if verbose {
						fmt.Println("Downloading " + link)
					}
					if err != nil {
						log.Println(err)
					} else {
						hasCookie := len(resp.Header.Get("Set-Cookie")) > 0
						externals[link] = Privacy{
							Typ:    "Script",
							Cookie: hasCookie,
						}
					}
				} else {
					externals[link] = Privacy{
						Typ:    "Script",
						Cookie: hasCookie,
					}
				}
			}
		}
	})

	c.OnHTML("link[href]", func(e *colly.HTMLElement) {
		link := e.Attr("href")
		rel := e.Attr("rel")
		if _, ok := externals[link]; !ok {
			if len(rel) > 0 && rel == "stylesheet" {
				hasCookie := false
				if checkForCookie {
					cssLink := ""
					if !strings.HasPrefix(link, "http") {
						cssLink = checkUrl + "/" + link
					} else {
						cssLink = link
					}
					if _, ok := checkedInternalCss[cssLink]; !ok {
						checkedInternalCss[cssLink] = true
						client := http.Client{
							Timeout: 5 * time.Second,
						}
						if verbose {
							fmt.Println("Downloading " + cssLink)
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
								hasCookie = len(resp.Header.Get("Set-Cookie")) > 0
							}
						}
					}
					if !strings.HasPrefix(link, checkUrl) && strings.HasPrefix(link, "https://") {
						externals[link] = Privacy{
							Typ:    "Css",
							Cookie: hasCookie,
						}
					}

				}
			}
		}
	})

	// Before making a request print "Visiting ..."
	c.OnRequest(func(r *colly.Request) {
	})

	return scanner
}
