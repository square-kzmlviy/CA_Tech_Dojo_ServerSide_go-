## 解説する命令とか
### http.HandleFunc()
今回は、レスポンスの処理を別ファイル化してある為分かり易く代入する。
```go
http.HandleFunc("/user/create",page.Create)
```
とは、下記です。
```go
http.HandleFunc("/user/create",func (w http.ResponseWriter, r *http.Request) )
```
```func (w http.ResponseWriter, r *http.Request)```ですが、```http.HandlerFunc```関数を無名関数にしたものだそうです。
#### http.ResponseWriterインタフェース
ハンドラ関数の第一引数には、http.ResponseWriterインタフェース型のオブジェクトが渡されます。 このオブジェクトはHTTPレスポンスヘッダへの値セットや、HTTPレスポンスボディへの出力に使用します。

#### http.Request構造体
ハンドラ関数の第二引数には、http.Request構造体型のオブジェクトが渡されます。 この構造体にはHTTPリクエストの内容が格納されています。<br>
Requestの中身は[こちらのDocument](https://golang.org/pkg/net/http/#Request)にありました。ほええ



### HTTPメソッド（GET/POST/PUTなど）

```go
switch r.Method {
    // func Create(w http.ResponseWriter, r *http.Request)
    case http.MethodPost:
    // func Get(w http.ResponseWriter, r *http.Request)
    case http.MethodGet:
    // func Update(w http.ResponseWriter, r *http.Request)
    case http.MethodPut
}
```
```r.Method```は、[http.Request構造体](####http.Request構造体)のMethodパラメータです。
公式Documentによると、<br>
```Methodは、HTTPメソッド（GET、POST、PUTなど）を指定します。クライアント要求の場合、空の文字列はGETを意味します。```(google翻訳)
と記載されていました。<br>

今回私が作成したAPIは、HTTPメソッドごとにURLが異なる為```switch cace文```を利用するべきか悩みましたが、明示的に表現したかったので記載しました。

### HTTPリクエストメッセージボディ

#### HTTPリクエストメッセージボディってなに
HTTPリクエストのBody部。あるサイトには、メモ書き程度と記載されていましたが要は受け渡すパラメータなんかが記載されているようです。

```go
body := r.Body
defer body.Close()
```
二行目ですが、[公式Document](https://golang.org/pkg/net/http/#overView)にあるように```The client must close the response body when finished with it(クライアントは、応答が終了したら応答本文を閉じる必要があります。)```とのことなので```body.Close()```しています。ちなみに```defer```は遅延実行です。

クローズを怠ると、TCPコネクションがクローズされないために、そのまま続けてHttpリクエストを発行していき、ファイルディスクリプタが枯渇してしまいます。


### sql.Open

```go
sql.Open("mysql", "root@/ca_tech_dojo")
```

外部ライブラリ：Mysqlのドライバを導入

```cmd
go get "github.com/go-sql-driver/mysql"
```

import文

```
import _ "github.com/go-sql-driver/mysql"
```

先頭の```_```がないと、パッケージのメンバを明示的に利用するコードが無いことからビルドエラーとなる。

#### 役割
ドライバの名前を指定して、データベースに接続する。
| パラメータ 	| 内容                            	|
|------------	|---------------------------------	|
| 第一引数   	| ドライバ名                      	|
| 第二引数   	| 接続情報（ユーザ名 パスワード） 	|
| 第一戻り値 	| DBハンドル                      	|
| 第二戻り値 	| Errorハンドル                   	|

### JSONの取り扱い

```go
// JSONを格納するための構造体の宣言
type UserName struct {
	Name string `json:"name"`
}

buf := new(bytes.Buffer)
io.Copy(buf, body)

var UserData UserName
json.Unmarshal(buf.Bytes(), &UserData)
```





