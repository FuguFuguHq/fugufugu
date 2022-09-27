package fugu

import (
	"log"
	"net/url"
	"strings"
)

type Product struct {
	Name    string
	Url     string
	Company *Company
}

type Company struct {
	Name    string
	Country string
	Privacy string
}

func ProductForUrl(products map[string]Product, u string) *Product {
	var product Product

	host, err := url.Parse("https://" + u)
	if err != nil {
		log.Fatalln(err)
	}
	// used publicprefix but it's a mess so
	// we create another mess :-)
	// will not work for e.g. co.uk
	parts := strings.Split(host.Hostname(), ".")
	if len(parts) < 2 {
		return nil
	}

	i := 0
	found := false
	for i <= len(parts)-2 && !found {
		product, found = products[strings.Join(parts[i:], ".")]
		i += 1
	}
	return &product
}

type ProductRaw struct {
}

func Products() map[string]Product {
	rawProducts := [][]string{
		[]string{"Stripe", "US", "stripe.com", "Payment", "https://stripe.com/de/privacy"},
		[]string{"Mailchimp", "US", "mailchimp.com", "Newsletter", "https://mailchimp.com/en/legal/"},
		[]string{"sendinblue", "EU", "sibforms.com", "Newsletter", "https://www.sendinblue.com/legal/privacypolicy/"},
		[]string{"Simple Analytics", "EU", "simpleanalyticscdn.com", "Analytics", "https://simpleanalytics.com/privacy-policy"},
		[]string{"rapidmail", "EU", "emailsys1a.net", "Newsletter", "https://www.rapidmail.com/data-security"},
		[]string{"Google", "US", "fonts.googleapis.com", "Google Fonts", ""},
		[]string{"Google", "US", "google.com", "Google", ""},
		[]string{"Google", "US", "googleapis.com", "Google", ""},
		[]string{"Google", "US", "gstatic.com", "Google", ""},
		[]string{"Google", "US", "google-analytics.com", "Google Analytics", ""},
		// probably also analytics, how to detect?
		// https://www.googletagmanager.com/gtag/js
		// Always GA included?
		[]string{"Google", "US", "googletagmanager.com", "Tagmanager", ""},
		[]string{"Squarespace", "US", "squarespace.com", "CMS", "https://www.squarespace.com/privacy"},
		[]string{"Twitter", "US", "twitter.com", "Social Media", "https://twitter.com/en/privacy"},
		[]string{"Tilda", "UK", "tildacdn.com", "CMS", "https://tilda.cc/privacy/"},
		[]string{"Calendly", "US", "calendly.com", "Calendar", "https://calendly.com/privacy"},
		[]string{"Cloudflare", "US", "cdnjns.cloudflare.com", "CDN", "https://www.cloudflare.com/de-de/privacypolicy/"},
		[]string{"Cloudflare", "US", "cloudflareinsights.com", "Analytics", "https://www.cloudflare.com/de-de/privacypolicy/"},
		[]string{"Cloudflare", "US", "cloudfront.net", "CDN", "https://www.cloudflare.com/de-de/privacypolicy/"},
		[]string{"Webflow", "US", "webflow.com", "CMS", "https://webflow.com/legal/privacy"},
		[]string{"Amazon", "US", "s3.amazonaws.com", "Cloud Storage", "https://aws.amazon.com/privacy/"},
		[]string{"storyblok", "EU", "storyblok.com", "CMS", "https://www.storyblok.com/legal/privacy-policy"},
		[]string{"Zendesk", "US", "zdassets.com", "Helpdesk", "https://www.zendesk.com/company/agreements-and-terms/privacy-policy-2021-12-06/"},
	}

	products := make(map[string]Product)
	for _, p := range rawProducts {
		company := &Company{
			Name:    p[0],
			Country: p[1],
			Privacy: p[4],
		}
		product := Product{
			Name:    p[3],
			Url:     p[2],
			Company: company,
		}
		products[p[2]] = product

	}
	return products
}
