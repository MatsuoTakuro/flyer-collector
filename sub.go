package main

import (
	"io"
	"log"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
)

func request(sc_url string) io.ReadCloser {
	// Reduce the burden on the target server
	time.Sleep(1 * time.Second)

	// Request the HTML page.
	res, err := http.Get(sc_url)
	if err != nil {
		log.Fatal(err)
	}
	if res.StatusCode != 200 {
		log.Fatalf("status code error: %d %s", res.StatusCode, res.Status)
	}
	return res.Body
}

// scrapeStoresList scrapes store names and urls from searched stores list pages
func scrapeStoresList(doc *goquery.Document, sc_url string) (string, []Store) {
	// Get the title of this page
	title := doc.Find("title").Text()

	// Find the items
	var stores []Store

	doc.Find("a.shop_index_card").Each(func(i int, s *goquery.Selection) {
		// For each item found, get the store url
		name := strings.TrimSpace(s.Find("div.name_text").Text())
		href, _ := s.Attr("href")
		url := toAbsUrl(sc_url, href)
		st := Store{name, url, nil}
		stores = append(stores, st)
	})
	return title, stores
}

func toAbsUrl(sc_url string, weburl string) string {
	baseurl, _ := url.Parse(sc_url)
	relurl, err := url.Parse(weburl)
	if err != nil {
		return ""
	}
	absurl := baseurl.ResolveReference(relurl)
	return absurl.String()
}
