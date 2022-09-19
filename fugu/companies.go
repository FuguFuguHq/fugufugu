package fugu

import (
	"golang.org/x/net/publicsuffix"
	"strings"
)

type Company struct {
	Name    string
	Country string
	Urls    []string
	Privacy string
}

func CompanyForUrl(companies map[string]Company, url string) *Company {
	var company Company

	// publicsuffix doesn't work for "d3e54v103j8qbb.cloudfront.net" ?
	if strings.HasSuffix(url, "cloudfront.net") {
		domain := "cloudfront.net"
		company = companies[domain]
	} else {
		domain, err := publicsuffix.EffectiveTLDPlusOne(url)
		if err == nil {
			company = companies[domain]
		}
	}
	return &company
}

func Companies() map[string]Company {
	raw := []Company{
		{
			Name:    "Simple Analytics",
			Country: "EU",
			Urls:    []string{"simpleanalyticscdn.com"},
			Privacy: "https://simpleanalytics.com/privacy-policy",
		},
		{
			Name:    "rapidmail",
			Country: "EU",
			Urls:    []string{"emailsys1a.net"},
			Privacy: "https://www.rapidmail.com/data-security",
		},
		{
			Name:    "Google",
			Country: "US",
			Urls:    []string{"googleapis.com", "gstatic.com"},
			Privacy: "",
		},
		{
			Name:    "Tilda",
			Country: "UK",
			Urls:    []string{"tildacdn.com"},
			Privacy: "https://tilda.cc/privacy/",
		},
		{
			Name:    "Calendly",
			Country: "US",
			Urls:    []string{"calendly.com"},
			Privacy: "https://calendly.com/privacy",
		},
		{
			Name:    "Cloudflare",
			Country: "US",
			Urls:    []string{"cloudflareinsights.com", "cloudfront.net"},
			Privacy: "https://www.cloudflare.com/de-de/privacypolicy/",
		},
		{
			Name:    "Webflow",
			Country: "US",
			Urls:    []string{"webflow.com"},
			Privacy: "https://webflow.com/legal/privacy",
		},
	}

	companies := make(map[string]Company)
	for _, c := range raw {
		for _, u := range c.Urls {
			companies[u] = c
		}
	}
	return companies
}
