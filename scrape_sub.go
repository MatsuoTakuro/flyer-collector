package main

import (
	"fmt"
	"log"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
)

func requestHTMLDoc(sc_url string) (*goquery.Document, error) {
	// Reduce the burden on the target server
	time.Sleep(1 * time.Second)

	// Request the HTML page.
	res, err := http.Get(sc_url)
	fmt.Printf("Sent http Get request to : %v\n", sc_url)
	if err != nil {
		log.Fatal(err)
	}
	if res.StatusCode != 200 {
		log.Fatalf("status code error: %d %s", res.StatusCode, res.Status)
	}
	fmt.Printf("- The respose status of the Get request sent is : %v (%v)\n", res.StatusCode, res.Status)

	// Load the HTML document
	defer res.Body.Close()
	doc, err := goquery.NewDocumentFromReader(res.Body)

	return doc, err
}

// addStores add each store's names, urls and flyers, from searched stores list pages, to store collection
func addStores(id *int, doc *goquery.Document, sc_url string) []Store {
	// Get the title of this page
	title := doc.Find("title").Text()
	fmt.Printf("\nStarted to scrape for stores list page's title : %v\n", title)

	// Find the store items
	var stores []Store
	doc.Find("a.shop_index_card").Each(func(i int, s *goquery.Selection) {
		// If specified "-store" as a command line argument
		if *maxStores > 0 && i >= *maxStores {
			return
		}
		// For each store item found, get it's id, name, url and flyers
		*id++
		name := strings.TrimSpace(s.Find("div.name_text").Text())
		href, _ := s.Attr("href")
		url := toAbsUrl(sc_url, href)
		flyers := scrapeFlyers(url)
		st := Store{*id, name, url, flyers}
		stores = append(stores, st)
	})
	return stores
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
	// 1st; Request the HTML page and Load the HTML document, for store detail page
	fmt.Println("\nStarted to scrape for store detail page")
	doc, err := requestHTMLDoc(sc_url)
	if err != nil {
		log.Fatal(err)
	}

	href, _ := doc.Find("li.shop_header_tab:nth-child(2) a").Attr("href")
	sc_url = toAbsUrl(sc_url, href)
	// 2nd; Request the HTML page and Load the HTML document, for the 1st flyer detail page
	fmt.Println("For the 1st flyer detail page")
	doc, err = requestHTMLDoc(sc_url)
	if err != nil {
		log.Fatal(err)
	}
	var flyers []Flyer
	var id = 1
	desc, _ := doc.Find("img.leaflet").Attr("alt")
	imgURL, _ := doc.Find("img.leaflet").Attr("src")
	fly1st := Flyer{id, desc, imgURL}
	flyers = append(flyers, fly1st)

	sc_url = makeFly2ndURL(sc_url)
	// 3nd; Request the HTML page and Load the HTML document, for the 2nd flyer detail page
	fmt.Println("For the 2nd flyer detail page")
	doc, err = requestHTMLDoc(sc_url)
	if err != nil {
		log.Fatal(err)
	}
	id++
	desc, _ = doc.Find("img.leaflet").Attr("alt")
	imgURL, _ = doc.Find("img.leaflet").Attr("src")
	fly2nd := Flyer{id, desc, imgURL}
	flyers = append(flyers, fly2nd)

	return flyers
}

func makeFly2ndURL(fly1stURL string) string {
	splitUrl := strings.Split(fly1stURL, "/")
	intFly2nd, err := strconv.Atoi(splitUrl[len(splitUrl)-1])
	if err != nil {
		log.Fatal(err)
	}
	splitUrl[len(splitUrl)-1] = strconv.Itoa(intFly2nd + 1)
	sc_url := strings.Join(splitUrl, "/")
	return sc_url
}
