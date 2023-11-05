# CharFluent

## 実装メモ

### Userの作成

**データベースの準備**  

`setup.sql`にデータベース定義を記述しておく.  

```bash
$ psql postgres
> CREATE DATABASE chatfluent;
$ psql -f data/setup.sql -d chatfluent
```

**データベースとの接続**  
`data.go`にデータベースとの接続処理を記述する.  
(DbConnection.Exex, DbConnection.Prepareの使い分けがわからない。)  

uuidを生成する関数の定義もここ.  


**Userの作成**  
`user.go`にuser構造体とsession構造体、それに関連するメソッドを記述する.  


**uuidの生成**  
`github.com/nu7hatch/gouuid`を参考に実装する。  
uuidはセッションの管理に使用する。  


### config設定

`config.go`に記述する.  
- os.Openで設定ファイルを開く
- json.NewDecoder()でデコーダーを生成する
- decoder.Decode(&config)でconfig構造体に読み込む

### サーバーの生成
`main.go`に記述する.  

```go
mux := http.NewServeMux()
mux.HandleFunc("/", index) // rootパスにアクセスした際にindexのハンドラを呼び出す
http.ListenAndServe("0.0.0.0:8080", mux) //マルチプレクサを使ってサーバーを起動

func index(w http.ResponseWriter, r *http.Request) {
  t, _ := template.ParseFiles("templates/index.html")
  t.Execute(w, nil)
  // templates/index.htmlのテンプレートをレンダリングする
}

```

### レイアウトの作成

`templates/layout.html`に記述する.  

```html
{{ define "layout" }}
<html>
  <body>
    {{ template "navbar" . }}
    <div>
      {{ template "content" . }}
    </div>
  </body>
</html>
{{ end }}
```


homeページの作成  
```html
{{ define "home" }}
<div class="container">
  {{ . }}
</div>

{{ end }}
```


