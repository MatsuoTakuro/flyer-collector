package main

import (
	"flag"
	"log"
)

const (
	// TODO: #1 検索条件はコマンドラインから入力する（県名("fukuoka")は番号へ変換 or コマンドライン上で選択 等）
	storeName string = "ディスカウントドラッグコスモス"
	prefName  string = "fukuoka"
)

var (
	maxStores *int
)

func main() {
	maxStores = flag.Int("store", 0, "Define the limit to stores that will be scraped ")
	flag.Parse()
	if *maxStores < 0 || *maxStores > 15 {
		log.Fatal("Do not set a value less than 0 or more than 16, as \"store\" command line argument")
	}
	scrapeTokubai(storeName, prefName)
}
