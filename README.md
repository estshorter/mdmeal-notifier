# mdmeal
## 動作環境
Windows 10

## How to use
まず、同じフォルダに`account.json`として
``` json
{
    "id": "ID",
    "password": "PASS"
}
```
のようなファイルを作る。

次に、`mdmeal.exe PATH_TO_JSON` とすれば起動できる。天津飯がメニューにある時はWindowsのトーストで通知される。
コマンドライン引数が与えられていない場合は、カレントディレクトリの`account.json`を読もうとする。
