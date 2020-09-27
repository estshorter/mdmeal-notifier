# mdmeal-notifier
天津飯とかのマズメシが出そうになったらLINEに通知できる。

## 動作環境
- Windows 10
- Chrome
- [ChromeDriver](https://chromedriver.chromium.org/)

## 使い方
まず、本ソフトのexeファイルがあるフォルダに`configs.json`として
``` json
{
	"triggerMenus": [
		"天津飯",
		"うどん"
	],
	"mdmealURL": "OUR_MDMEAL_LOGIN_URL",
	"mdmealAcount": {
		"id": "YOUR_ID",
		"password": "YOUR_PASS"
	},
	"lineNotifyToken": "YOUR_TOKEN",
}
```
のようなファイルを作る。`lineNotifyToken`は https://qiita.com/ken_yoshi/items/7879b3117d298a143101 などを参考に生成よい。

次に、`mdmeal-notifier.exe PATH_TO_CONFIGS_JSON` とすれば起動できる。`triggerMenus`のいずれかがメニューにある時はLINEに通知される。
コマンドライン引数が与えられていない場合は、カレントディレクトリの`configs.json`を読もうとする。

