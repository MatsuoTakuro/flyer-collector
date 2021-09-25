package main

import (
	"log"
	"net/http"
	"net/url"
	"strconv"
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
		// For each store item found, get it's id, name, url and flyers
		*id++
		name := strings.TrimSpace(s.Find("div.name_text").Text())
		href, _ := s.Attr("href")
		url := toAbsUrl(sc_url, href)
		flyers := scrapeFlyers(url)
		st := Store{*id, name, url, flyers}
		// st := Store{*id, name, url, nil}
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
	// 1st Request the HTML page and Load the HTML document, for store detail page
	doc, err := request(sc_url)
	if err != nil {
		log.Fatal(err)
	}

	href, _ := doc.Find("li.shop_header_tab:nth-child(2) a").Attr("href")
	sc_url = toAbsUrl(sc_url, href)
	// 2nd Request the HTML page and Load the HTML document, for the 1st flyer detail page
	doc, err = request(sc_url)
	if err != nil {
		log.Fatal(err)
	}
	var flyers []Flyer
	var id = 1
	desc, _ := doc.Find("img.leaflet").Attr("alt")
	image, _ := doc.Find("img.leaflet").Attr("src")
	fly1st := Flyer{id, desc, image}
	flyers = append(flyers, fly1st)
	if len(flyers) == 1 {
		s := strings.Split(sc_url, "/")
		intFlyNum, err := strconv.Atoi(s[len(s)-1])
		if err != nil {
			log.Fatal(err)
		}
		flyNum := strconv.Itoa(intFlyNum + 1)
		s[len(s)-1] = flyNum
		sc_url = strings.Join(s, "/")
		// 3nd Request the HTML page and Load the HTML document, for the 2nd flyer detail page
		doc, err = request(sc_url)
		if err != nil {
			log.Fatal(err)
		}
		id++
		desc, _ = doc.Find("img.leaflet").Attr("alt")
		image, _ = doc.Find("img.leaflet").Attr("src")
		fly2nd := Flyer{id, desc, image}
		flyers = append(flyers, fly2nd)
	}
	return flyers
}
