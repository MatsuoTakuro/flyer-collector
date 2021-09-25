package main

const (
	// TODO: #1 検索条件はコマンドラインから入力する（県名("fukuoka")は番号へ変換 or コマンドライン上で選択 等）
	storeName string = "ディスカウントドラッグコスモス"
	prefName  string = "fukuoka"
)

func main() {
	scrapeTokubai(storeName, prefName)
}
