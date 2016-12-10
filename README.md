# sacloud-balus-old

[バルス](https://ja.wikipedia.org/wiki/飛行石#.E3.81.8A.E3.81.BE.E3.81.98.E3.81.AA.E3.81.84)をさくらのクラウド上で実装するためのソフトウェアパーツ

Arduino + さくらのIoT 通信モジュールでマイク入力を収集し、「バルス」と言ったかを判定するバルス判定サーバーの試験実装です。

**当プロジェクトは実験的なコードを含む未完成なものです。**

当プロジェクトの詳細は[はてなブログ](http://febc-yamamoto.hatenablog.com/entry/sacloud-balus)を参照ください。

## 概要

このプロジェクトは、さくらのクラウド上に「バルス」を実装するためのソフトウェアパーツです。

以下のような処理を行っています。

  - さくらのIoT PlatformからのWebhookを待ち受け
  - マイク入力を受け取り、waveファイル作成
  - Azure cognitive servicesを呼び出し、音声認識
  - バルスと言われたか判定


IoT Platformに繋ぐためのArduino用プロジェクトは以下のリポジトリです。

[sakura_mic_input](https://github.com/yamamoto-febc/sakura_mic_input)

## License

 `sacloud-balus` Copyright (C) 2016 Kazumichi Yamamoto.

  This project is published under [Apache 2.0 License](LICENSE.txt).
  
## Author

  * Kazumichi Yamamoto ([@yamamoto-febc](https://github.com/yamamoto-febc))

