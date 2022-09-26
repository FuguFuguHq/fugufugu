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
	Collector    *colly.Collector
	Pages        *uint64
	PrivacyPages *[]PrivacyPage
}

type PrivacyPage struct {
	URL   *string
	Title *string
}

func CheckCookie(checkUrl string, verbose bool) bool {
	hasCookie := false
	client := http.Client{
		Timeout: 5 * time.Second,
	}
	resp, err := client.Get(checkUrl)
	if verbose {
		fmt.Println("Downloading " + checkUrl)
	}
	if err != nil {
		log.Println(err)
	} else {
		hasCookie = len(resp.Header.Get("Set-Cookie")) > 0
	}
	return hasCookie
}

func isPrivacyPage(url string, title string) bool {
	toLowerUrl := strings.ToLower(url)
	toLowerTitle := strings.ToLower(title)
	return strings.Contains(toLowerUrl, "privacy") || strings.Contains(toLowerTitle, "datenschutz")
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
		colly.MaxDepth(5),
	)

	scanner := Scanner{
		Collector:    c,
		Pages:        new(uint64),
		PrivacyPages: &[]PrivacyPage{},
	}

	privacyPages := make(chan PrivacyPage, 5)
	go func() {
		for {
			page, more := <-privacyPages
			if more {
				*scanner.PrivacyPages = append(*scanner.PrivacyPages, page)
			}
		}
	}()

	c.OnRequest(func(r *colly.Request) {
		if verbose {
			fmt.Println("Checking " + r.URL.String())
		}
	})

	c.OnResponse(func(r *colly.Response) {
		atomic.AddUint64(scanner.Pages, 1)
		//	fmt.Printf("Link found: %v\n", r.Headers)
	})

	c.OnHTML("title", func(e *colly.HTMLElement) {
		if isPrivacyPage(e.Request.URL.Path, e.Text) {
			privacyPages <- PrivacyPage{
				URL:   &e.Request.URL.Path,
				Title: &e.Text,
			}
		}
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
					hasCookie = CheckCookie(checkUrl, verbose)
				}
				externals[link] = Privacy{
					Typ:    "Image",
					Cookie: hasCookie,
				}
				if verbose {
					fmt.Println("External image " + link)
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
					hasCookie = CheckCookie(checkUrl, verbose)
				}
				externals[link] = Privacy{
					Typ:    "Script",
					Cookie: hasCookie,
				}
				if verbose {
					fmt.Println("External script " + link)
				}
			}
		}
	})

	c.OnHTML("link[href]", func(e *colly.HTMLElement) {
		link := e.Attr("href")
		rel := e.Attr("rel")
		if _, ok := externals[link]; !ok {
			if len(rel) > 0 && strings.ToLower(rel) == "stylesheet" {
				hasCookie := false
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
							imports := ImportsFromCss(string(b), verbose)
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
					if verbose {
						fmt.Println("External CSS: " + link)
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
