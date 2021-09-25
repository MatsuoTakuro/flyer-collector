package main

import (
	"log"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
)

func request(sc_url string) (*goquery.Document, error) {
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

	// Load the HTML document
	defer res.Body.Close()
	doc, err := goquery.NewDocumentFromReader(res.Body)

	return doc, err
}

// addStores add each store's names, urls and flyers, from searched stores list pages, to store collection
func addStores(id *int, doc *goquery.Document, sc_url string) (string, []Store) {
	// Get the title of this page
	title := doc.Find("title").Text()

	// Find the store items
	var stores []Store
	doc.Find("a.shop_index_card").Each(func(i int, s *goquery.Selection) {
		// For each store item found, get it's name, url and flyers
		*id++
		name := strings.TrimSpace(s.Find("div.name_text").Text())
		href, _ := s.Attr("href")
		url := toAbsUrl(sc_url, href)
		// flyers := scrapeFlyers(url)
		// st := Store{i, name, url, flyers}
		st := Store{*id, name, url, nil}
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

func scrapeFlyers(sc_url string) []Flyer {
	// Request the HTML page and Load the HTML document
	doc, err := request(sc_url)
	if err != nil {
		log.Fatal(err)
	}
	// Find the flyer items
	var flyers []Flyer
	doc.Find("a.shop_index_card").Each(func(i int, s *goquery.Selection) {
		// For each flyer item found, get the flyer's id, title and image
		title := strings.TrimSpace(s.Find("div.name_text").Text())
		image := "image"
		fly := Flyer{i, title, image}
		flyers = append(flyers, fly)
	})
	return flyers
}
