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
	// TODO:コマンドラインから入力するように変更する（"福岡県"は40へ自動変換 or コマンドライン上で選択）
	baseURL string = "https://tokubai.co.jp"
	rawName string = "ディスカウントドラッグコスモス"
	fukuoka int    = 40
)

func scrapeTokubai() {
	// Set the target url
	storeName := url.QueryEscape(rawName)
	sc_url := fmt.Sprintf("%v/%v/prefectures/%d", baseURL, storeName, fukuoka)

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
	title, stores = scrapeThisPage(doc, sc_url)
	// TODO: チラシ画像を取得する
	// TODO: OCRでスキャンする
	// TODO: スキャンされた情報を整形する
	// TODO: ファイルに保存する

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

		_, tmpStores := scrapeThisPage(doc, sc_url)
		stores = append(stores, tmpStores...)
		href, exists = doc.Find("span.next a").Attr("href")
	}

	fmt.Printf("Page title: %v\n", title)
	for i, s := range stores {
		fmt.Printf("Store name #%d : %v\n", i+1, s.name)
		fmt.Printf("Store URL  #%d : %v\n\n", i+1, s.storeURL)
	}

}

func main() {
	scrapeTokubai()
}
