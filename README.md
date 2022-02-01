# flyer-collector


### 1, 概要
- チラシ自動収集ツール（Webスクレイピング+ OCRスキャン）


### 2, スクレイピング対象のWebサイト
- [tokubai](https://tokubai.co.jp/)(チラシ共有サイト)
- ただ今回は、検索条件を`福岡県`&&`ディスカウントドラッグコスモス(チェーン)`に絞る


### 3, 使用するOCR API
- [GCP Vision API](https://cloud.google.com/vision?authuser=1)
  - 有料サービス（ただ、GCPへ新規登録してから一定期間は無料枠で使用できる）


### 4, 開発言語と主に使用するパッケージ
- 開発言語
  -  [Go@v1.17](https://go.dev/)
    - バージョン・アップデートあり次第、更新する予定
- 主に使用するパッケージ
  - [github.com/PuerkitoBio/goquery@v1.7.1](https://pkg.go.dev/github.com/PuerkitoBio/goquery@v1.7.1)
    - HTMLドキュメントからの読み取り
  - [cloud.google.com/go/vision@v1.0.0](https://pkg.go.dev/cloud.google.com/go/vision@v1.0.0)
    - チラシ画像へのOCRスキャン
