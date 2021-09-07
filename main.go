package main

import (
	"fmt"
	"log"
	"net/http"
	"net/url"

	"github.com/PuerkitoBio/goquery"
)

const (
	// TODO:コマンドラインから入力するように変更する（"福岡県"は40へ自動変換 or コマンドライン上で選択）
	fukuoka int    = 40
	rawName string = "ディスカウントドラッグコスモス"
)

func ExampleScrape() {
	// Set the target url
	storeName := url.QueryEscape(rawName)
	sc_url := fmt.Sprintf("https://tokubai.co.jp/%v/prefectures/%d", storeName, fukuoka)

	// Request the HTML page.
	res, err := http.Get(sc_url)
	if err != nil {
		log.Fatal(err)
	}
	defer res.Body.Close()
	if res.StatusCode != 200 {
		log.Fatalf("status code error: %d %s", res.StatusCode, res.Status)
	}

	// Load the HTML document
	doc, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		log.Fatal(err)
	}

	// Get the title of this page
	title := doc.Find("title").Text()
	fmt.Println("Page title: " + title)

	// Find the review items
	doc.Find("section.shop_index_cards a").Each(func(i int, s *goquery.Selection) {
		// For each item found, get the store url
		// TODO: ページネーションのリンクはさらに展開して店舗URLを取得するようにする
		href, _ := s.Attr("href")
		burl, _ := url.Parse(sc_url)

		var full_url = toAbsUrl(burl, href)
		fmt.Printf("In-page link URL #%d : %v\n", i+1, full_url)
		// TODO: チラシ画像を取得する
		// TODO: OCRでスキャンする、整理する

		// TODO: ファイルに保存する
	})

}

func toAbsUrl(baseurl *url.URL, weburl string) string {
	relurl, err := url.Parse(weburl)
	if err != nil {
		return ""
	}
	absurl := baseurl.ResolveReference(relurl)
	return absurl.String()
}

func main() {
	ExampleScrape()
}
