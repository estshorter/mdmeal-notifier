# mdmeal
天津飯とかのマズメシが出そうになったらLINE or Win10トーストに通知できる。

## 動作環境
- Windows 10
- Chrome
- [ChromeDriver](https://chromedriver.chromium.org/)

## 使い方
まず、同じフォルダに`configs.json`として
``` json
{
	"triggerMenus": [
		"天津飯"
	],
	"mdmealURL": "MDMEAL_LOGIN_ID",
	"mdmealAcount": {
		"id": "YOUR_ID",
		"password": "YOUR_PASS"
	},
	"lineNotifyToken": "YOUR_TOKEN",
	"appID": "{1AC14E77-02E7-4E5D-B744-2EB1AE5198B7}\\notepad.exe"
}
```
のようなファイルを作る。`lineNotifyToken`は https://qiita.com/ken_yoshi/items/7879b3117d298a143101 などを参考に作れるが、
なくてもアプリ起動時に`-line=false`と指定すれば不要。このときはWin10の通知機能が使用される。

次に、`mdmeal-notifer.exe PATH_TO_JSON` とすれば起動できる。`triggerMenus`のいずれかがメニューにある時はLINE or Win10に通知される。
コマンドライン引数が与えられていない場合は、カレントディレクトリの`configs.json`を読もうとする。

## コマンドライン引数
``` ps1
PS > .\mdmeal-notifier.exe -h 
Usage of mdmeal-notifier.exe:
  -line
        use LINE Notify (default true)
```
