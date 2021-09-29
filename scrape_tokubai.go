package main

import (
	"fmt"
	"log"
	"net/url"
)

type Store struct {
	id     int
	name   string
	url    string
	flyers []Flyer
}
type Flyer struct {
	id    int
	desc  string
	image string
}

const (
	tokubaiBaseURL string = "https://tokubai.co.jp"
)

func getPrefsList() map[string]int {
	return map[string]int{"fukuoka": 40}
}

func scrapeTokubai(rawStoreName string, prefName string) {
	// Set the target url
	storeName := url.QueryEscape(rawStoreName)
	prefsList := getPrefsList()
	sc_url := fmt.Sprintf("%v/%v/prefectures/%d", tokubaiBaseURL, storeName, prefsList[prefName])

	// Request the HTML page and Load the HTML document, for stores list page that is the search result
	fmt.Println("\nStarted to scrape for stores list page that is the search result")
	doc, err := request(sc_url)
	if err != nil {
		log.Fatal(err)
	}

	var title string
	var stores []Store
	var id int
	title, stores = addStores(&id, doc, sc_url)
	// TODO: #2 チラシ画像を取得する
	// TODO: #3 OCRでスキャンする(GCP Vision APIを使用、コストは要検討)
	// TODO: #4 スキャンされた情報を整形し、ファイルに保存する

	// // Check if next page exists
	// href, exists := doc.Find("span.next a").Attr("href")
	// // Scrape the next page if it exists
	// for exists {
	// 	// Set the target url
	// 	next_sc_url := toAbsUrl(sc_url, href)

	// 	// Request the HTML page and Load the HTML document
	// 	doc, err := request(next_sc_url)
	// 	if err != nil {
	// 		log.Fatal(err)
	// 	}

	// 	// For each store item found, get it's name, url and flyers
	// 	_, tmpStores := addStores(&id, doc, sc_url)
	// 	stores = append(stores, tmpStores...)

	// 	// Check if next page exists, recursively
	// 	href, exists = doc.Find("span.next a").Attr("href")
	// }

	fmt.Printf("\ntitle: %v\n\n", title)
	for _, st := range stores {
		fmt.Printf("store no.%d\n", st.id)
		fmt.Printf("  name    : %v\n", st.name)
		fmt.Printf("  url     : %v\n", st.url)
		for _, fly := range st.flyers {
			fmt.Printf("  flyer #%d\n", fly.id)
			fmt.Printf("    desc  : %v\n", fly.desc)
			fmt.Printf("    image : %v\n", fly.image)
		}
		fmt.Println()
	}
}
