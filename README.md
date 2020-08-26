# mdmeal-notifier
天津飯とかのマズメシが出そうになったらLINE or Win10トーストに通知できる。

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
	"appID": "{1AC14E77-02E7-4E5D-B744-2EB1AE5198B7}\\notepad.exe"
}
```
のようなファイルを作る。`lineNotifyToken`は https://qiita.com/ken_yoshi/items/7879b3117d298a143101 などを参考に生成、入力すればよいが、
アプリ起動時に`-line=false`と指定すれば不要。このときはWin10の通知機能が使用される。

次に、`mdmeal-notifier.exe PATH_TO_CONFIGS_JSON` とすれば起動できる。`triggerMenus`のいずれかがメニューにある時はLINE or Win10に通知される。
コマンドライン引数が与えられていない場合は、カレントディレクトリの`configs.json`を読もうとする。

## コマンドライン引数
``` ps1
PS > .\mdmeal-notifier.exe -h 
Usage of mdmeal-notifier.exe:
  -line
        use LINE Notify (default true)
```
