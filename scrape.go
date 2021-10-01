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
	id     int
	desc   string
	imgURL string
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
	fmt.Println("\nStarted to scrape for stores list page that is the search results")
	doc, err := requestHTMLDoc(sc_url)
	if err != nil {
		log.Fatal(err)
	}

	var stores []Store
	var id int
	_, stores = addStores(&id, doc, sc_url)
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

	// TODO: #3 OCRでスキャンする(GCP Vision APIを使用)
	visionScann("/Users/user/go/src/flyer-collector/files/001_ディスカウントドラッグコスモス 久留米南店/00101_2021年9月25日〜10月1日までのチラシ_20210929185250.jpg")

	// TODO: #4 スキャンされた情報を整形し、ファイルに保存する

	fmt.Println("\nStarted to save images of gotten flyers")
	saveFlyImgsFrom(stores)

	fmt.Println("\nScraping and OCR-Scanning are finished.")
	fmt.Println()
}
