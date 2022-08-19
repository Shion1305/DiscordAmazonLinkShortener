# DiscordLinkShortener

チャットに送られた長いリンクを検出し、短縮するDiscordBotです。

現在は以下のリンクに対応をしています。

- [Amazon Japanのリンク](#Amazonのリンク構成)


# それぞれのリンク構成
## Amazonのリンク構成
Amazon.co.jpのリンクの一例を次に挙げる
> https://www.amazon.co.jp/-/en/%E3%82%A2%E3%82%A4%E3%83%86%E3%83%83%E3%82%AFIT%E4%BA%BA%E6%9D%90%E6%95%99%E8%82%B2%E7%A0%94%E7%A9%B6%E9%83%A8/dp/4865752536

通常のリンクには[パーセントエンコーディング](https://ja.wikipedia.org/wiki/%E3%83%91%E3%83%BC%E3%82%BB%E3%83%B3%E3%83%88%E3%82%A8%E3%83%B3%E3%82%B3%E3%83%BC%E3%83%87%E3%82%A3%E3%83%B3%E3%82%B0)で変換された商品名にdpが続き、そのあとに[ASNIコード](https://en.wikipedia.org/wiki/Amazon_Standard_Identification_Number)が続く。

このうちASINコードのみを抽出して以下のように変換することリンクの短縮を実現することが可能である。

> https://www.amazon.co.jp/dp/(ASINコード)

ASINコードは10文字の英数字で構成されている