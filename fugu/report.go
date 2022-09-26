package fugu

import (
	"fmt"
	"github.com/jedib0t/go-pretty/v6/table"
	"os"
	"sort"
)

type Report struct {
	Url            string
	ImageCount     int
	CssCount       int
	ScriptCount    int
	PageCount      uint64
	CheckForCookie bool

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

func PrintReport(r Report) {
	if len(r.Externals) == 0 {
		fmt.Printf("Summary %s: %d pages\n", r.Url, r.PageCount)
		fmt.Println("Cool, no external resources found!")
	} else {
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
		if r.CheckForCookie {
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
			if r.CheckForCookie {
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

func ReportFromScanner(checkUrl string, checkForCookie bool, privacies map[string]SitePrivacy, scanner Scanner) Report {
	v := make([]SitePrivacy, 0, len(privacies))

	for _, value := range privacies {
		value.Rank = Rank(value)
		v = append(v, value)
	}

	sort.Slice(v, func(i, j int) bool {
		return v[i].Rank > v[j].Rank
	})

	r := Report{
		Url:            checkUrl,
		PageCount:      *scanner.Pages,
		Externals:      make([]ReportExternal, 0),
		PrivacyPages:   make([]ReportPrivacyPage, 0),
		CheckForCookie: checkForCookie,
	}

	for _, p := range v {
		r.ScriptCount += p.ScriptCount
		r.ImageCount += p.ImgCount
		r.CssCount += p.CssCount
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
