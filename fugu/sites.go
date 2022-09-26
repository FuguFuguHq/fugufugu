package fugu

import (
	"net/url"
)

type SitePrivacy struct {
	Product     *Product
	Url         *string
	ImgCount    int
	CssCount    int
	ScriptCount int
	Cookie      bool
	Rank        int
}

func FromExternals(products map[string]Product, externals map[string]Privacy) map[string]SitePrivacy {
	sites := make(map[string]SitePrivacy)

	for urlString, p := range externals {
		u, err := url.Parse(urlString)
		if err != nil {
			panic(err)
		}
		if _, ok := sites[u.Host]; !ok {
			product := ProductForUrl(products, u.Host)
			sP := SitePrivacy{
				Product: product,
				Url:     &u.Host,
				Cookie:  p.Cookie,
			}
			if p.Typ == "Script" {
				sP.ScriptCount += 1
			} else if p.Typ == "Image" {
				sP.ImgCount += 1
			} else if p.Typ == "Css" {
				sP.CssCount += 1
			}
			sites[u.Host] = sP
		} else {
			sP := sites[u.Host]
			sP.Cookie = sP.Cookie && p.Cookie
			if p.Typ == "Script" {
				sP.ScriptCount += 1
			} else if p.Typ == "Image" {
				sP.ImgCount += 1
			} else if p.Typ == "Css" {
				sP.CssCount += 1
			}
			sites[u.Host] = sP
		}
	}
	return sites
}
