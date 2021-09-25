package main

import (
	"fmt"
	"log"
	"net/url"

	"github.com/PuerkitoBio/goquery"
)

type Store struct {
	name      string
	storeURL  string
	flyerImgs []string
}

const (
	baseURL string = "https://tokubai.co.jp"
)

func getPrefsList() map[string]int {
	return map[string]int{"fukuoka": 40}
}

func scrapeTokubai(rawStoreName string, prefName string) {
	// Set the target url
	storeName := url.QueryEscape(rawStoreName)
	prefsList := getPrefsList()
	sc_url := fmt.Sprintf("%v/%v/prefectures/%d", baseURL, storeName, prefsList[prefName])

	// Request the HTML page.
	resBody := request(sc_url)

	// Load the HTML document
	doc, err := goquery.NewDocumentFromReader(resBody)
	resBody.Close()
	if err != nil {
		log.Fatal(err)
	}

	var title string
	var stores []Store
	title, stores = scrapePage(doc, sc_url)
	// TODO: #2 チラシ画像を取得する
	// TODO: #3 OCRでスキャンする(GCP Vision APIを使用、コストは要検討)
	// TODO: #4 スキャンされた情報を整形し、ファイルに保存する

	// Go to next page
	href, exists := doc.Find("span.next a").Attr("href")
	for exists {
		// Set the target url
		next_sc_url := toAbsUrl(sc_url, href)

		// Request the HTML page.
		resBody := request(next_sc_url)

		// Load the HTML document
		doc, err := goquery.NewDocumentFromReader(resBody)
		resBody.Close()
		if err != nil {
			log.Fatal(err)
		}

		_, tmpStores := scrapePage(doc, sc_url)
		stores = append(stores, tmpStores...)
		href, exists = doc.Find("span.next a").Attr("href")
	}

	fmt.Printf("Page title: %v\n\n", title)
	for i, s := range stores {
		fmt.Printf("Store name #%d : %v\n", i+1, s.name)
		fmt.Printf("Store URL  #%d : %v\n\n", i+1, s.storeURL)
	}

}
