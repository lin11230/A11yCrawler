package main

import (
	"fmt"
	"net/http"
	"net/url"
	"os"
	"strings"

	"golang.org/x/net/html"
)

var (
	visitedUrl = make(map[string]bool)
)

// Helper function to pull the href attribute from a Token
func getHref(t html.Token) (ok bool, href string) {
	// Iterate over all of the Token's attributes until we find an "href"
	for _, a := range t.Attr {
		if a.Key == "href" {
			href = a.Val
			ok = true
		}
	}

	// "bare" return will return the variables (ok, href) as defined in
	// the function definition
	return
}

// Helper function to pull the href attribute from a Token
func get_IMG_alt(t html.Token) (ok bool, alt string) {
	// Iterate over all of the Token's attributes until we find an "alt"
	for _, a := range t.Attr {
		if a.Key == "alt" {
			alt = a.Val
			ok = true
		}
	}

	// "bare" return will return the variables (ok, ) as defined in
	// the function definition
	return
}

// Helper function to pull the src attribute from a Token
func get_IMG_src(t html.Token) (ok bool, src string) {
	// Iterate over all of the Token's attributes until we find an "alt"
	for _, s := range t.Attr {
		if s.Key == "src" {
			src = s.Val
			ok = true
		}
	}

	// "bare" return will return the variables (ok, ) as defined in
	// the function definition
	return
}

// Extract all http** links from a given webpage
func crawl(url string, ch map[string]bool, chalt map[string]bool) {
	fmt.Println("\n抓取網頁:", url)
	resp, err := http.Get(url)
	if err != nil {
		fmt.Println("錯誤: Failed to crawl \"" + url + "\"")
		return
	}

	b := resp.Body
	defer b.Close() // close Body when the function returns

	z := html.NewTokenizer(b)

	for {
		tt := z.Next()

		switch {
		case tt == html.ErrorToken:
			// End of the document, we're done
			return
		case tt == html.StartTagToken:
			t := z.Token()

			isImg := t.Data == "img"
			if isImg {
				//log.Println("img tag data:", t.Data)
				ok, src := get_IMG_src(t)
				if ok {
					fmt.Println("  圖檔名：", src)
					aok, alt := get_IMG_alt(t)
					if aok {
						fmt.Println("說明文字：", alt)
						if alt == "" {
							fmt.Println("*** 沒有說明文字! ***")
						}
					} else {
						fmt.Println("*** 沒有說明文字標籤! ***")
						alt = "無圖"
					}
					chalt[src+": "+alt] = true
				}
			}

			// Check if the token is an <a> tag
			isAnchor := t.Data == "a"
			if !isAnchor {
				//log.Println("found anchor:", t.Data)
				continue
			}

			// Extract the href value, if there is one
			ok, url := getHref(t)
			if !ok {
				continue
			} else {
				//log.Println("url is ", url)
			}

			// Make sure the url begines in http**
			hasAction := strings.Index(url, ".action") > 0
			if hasAction {
				ch[url] = true
			}
		}
	}
}

func main() {
	foundUrls := make(map[string]bool)
	foundAlts := make(map[string]bool)
	seedUrl := os.Args[1]

	u, err := url.Parse(seedUrl)
	if err != nil {
		fmt.Println("\n", err)
	}
	p := strings.SplitAfterN(u.Path, "/", 3)
	fmt.Println(u.Scheme, u.Host, p[1])
	hostname := u.Scheme + "://" + u.Host + "/" + p[1]

	crawl(seedUrl, foundUrls, foundAlts)

	fmt.Println("\n檢查網址:", seedUrl)
	fmt.Println("\nFound", len(foundUrls), "unique urls:\n")

	for url, _ := range foundUrls {
		//fmt.Println("\n - " + url + "\n")
		//fmt.Println(strings.Index(url, "Fpage.action"))
		if strings.Index(url, "Fpage.action") == 0 {
			suburl := hostname + url
			crawl(suburl, foundUrls, foundAlts)
		}
		//Index.action
		if strings.Index(url, "Index.action") == 0 {
			suburl := hostname + url
			crawl(suburl, foundUrls, foundAlts)
		}
		//Sitemap.action
		if strings.Index(url, "Sitemap.action") == 0 {
			suburl := hostname + url
			crawl(suburl, foundUrls, foundAlts)
		}
		//News!one.action
		if strings.Index(url, "News!one.action") == 0 {
			suburl := hostname + url
			crawl(suburl, foundUrls, foundAlts)
		}
		//Ad.action
		if strings.Index(url, "Ad.action") == 0 {
			suburl := hostname + url
			crawl(suburl, foundUrls, foundAlts)
		}
		//News.action
		if strings.Index(url, "News.action") == 0 {
			suburl := hostname + url
			crawl(suburl, foundUrls, foundAlts)
		}
		//News!link.action
		if strings.Index(url, "News!link.action") == 0 {
			suburl := hostname + url
			crawl(suburl, foundUrls, foundAlts)
		}
		//Other links contains vghtpe.gov.tw
		if strings.Index(url, "vghtpe.gov.tw") > 0 {
			crawl(url, foundUrls, foundAlts)
		}
	}

	//fmt.Println("\nFound", len(foundAlts), "unique alts:\n")

	//for alt, _ := range foundAlts {
	//fmt.Println(" - " + alt)
	//}

}
