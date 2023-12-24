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

userに関するテストは、`user_test.go`に記述する。
```go
func Test_UserCreate(t *testing.T) {
  setup()
  if err := users[0].Create(); err != nil {
    t.Error(err, "Cannot create user.")
  }
  if users[0].Id == 0 {
    t.Errorf("No id or created_at in user")
  }
  ...
}
```


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

### staticファイルの取り扱い処理
config.Static = "public"に設定しているので、publicフォルダに静的ファイルを格納し、
main.goでそれらを読み込む設定をする。  

```go[main.go]
files := http.FileServer(http.Dir(config.Static))
mux.Handle("/static/", http.StripPrefix("/static/", files))
```



## ログイン状態の判定
`route_main.go`に記述する.  

オリジナルのソースコードは下記のようになっている。
```go
func index(w http.ResponseWriter, r *http.Request) {
  threads, err := data.Threads()
  if err != nil {
    error_message(w, r, "error")
  } else {
    _, err := session(w, r)
    if err != nil {
      generateHTML(w, threads, "layout", "public.navbar", "index")
    } else {
      generateHTML(w, threads, "layout", "private.navbar", index)
    }
  }
}
```
session()ヘルパーは、`config.go`に下記のように定義する。  

```go
// Userのログイン判定
func session(w http.ResponseWriter, r *http.Request) (data.Session, error) {
	cookie, err := r.Cookie("_cookie")
	if err == nil {
		sess := data.Session{Uuid: cookie.Value}
		if ok, _ := sess.Check(); !ok {
			err = errors.New("invalid session")
		}

	}
	return data.Session{}, err
}
```


## templateからの構造体のメソッド呼び出し
例えば、`thread.User()`というメソッドをtemplateから呼び出したい場合、
HandlerFuncには構造体のアドレスを渡す必要がある。

```go
func threadRead(w http.ResponseWriter, r *http.Request) {
  ...中略...
  th, err := data.ThreadByUuid(uuid)
  t.ExecuteTemplate(w, "layout", &th)
}
```

## Delete メソッドの実装
formは`method="delete"`をサポートしていないので、
下記のようにformを実装する。

```html
  <div class="container delete-section">
    <form action="/thread/delete" method="post">
      <input type="hidden" name="uuid" value="{{.Uuid}}">
      <input type="hidden" name="_method" value="delete">
      <button type="submit" class="btn btn-danger">このThreadを削除</button>
    </form>
  </div>
```


## デプロイ

Goodle Cloud Runを使う。(App Engineとの接続がうまくできなかったので、Cloud Runで。)

[ここ](https://github.com/tker-78/cloudrun-sample)の手順でアプリケーションを作成.  


base.goに`cloud.google.com/go/cloudsqlconn`をインポートする。 

## projectの作成と設定

プロジェクトを指定する
```
gcloud config set project chatfluent-409007
```

API権限を与える
```
gcloud services enable compute.googleapis.com sqladmin.googleapis.com \
  run.googleapis.com artifactregistry.googleapis.com \
  cloudbuild.googleapis.com servicenetworking.googleapis.com
```

サービスアカウントを作成する
```
$ gcloud iam service-accounts create chatfluent-service-account \ 
  --display-name="Chatfluent Service Account" 
```

適切なロールをサービスアカウントに与える
```
$ gcloud projects add-iam-policy-binding chatfluent-409007 \
 --member="serviceAccount:chatfluent-service-account@chatfluent-409007.iam.gserviceaccount.com" \
  --role="roles/cloudsql.client"

```

```
$ gcloud projects add-iam-policy-binding chatfluent-409007 \
--member="serviceAccount:chatfluent-service-account@chatfluent-409007.iam.gserviceaccount.com" \
--role="roles/cloudsql.instanceUser"
```

```
$ gcloud projects add-iam-policy-binding chatfluent-409007 \
--member="serviceAccount:chatfluent-service-account@chatfluent-409007.iam.gserviceaccount.com" \
--role="roles/logging.logWriter"
```



### Cloud SQLの作成

インスタンスの作成
```
gcloud sql instances create chatfluent-instance \
 --database-version=POSTGRES_14 \
 --cpu=1 \
 --memory=4GB \
 --region=us-central1\
 --database-flags=cloudsql.iam_authentication=on \
 --root-password=chat
```

```
gcloud sql databases create chatfluent_db --instance=chatfluent-instance
```

```
 gcloud sql users create chatfluent-service-account@chatfluent-409007.iam \
 --instance=chatfluent-instance \
 --type=cloud_iam_service_account
```

Cloud SQLに`setup.sql`を流し込む。

このとき、アプリケーションのサービスアカウントを選択する(ここ大事)。  


```
gcloud run deploy chatfluent \
  --region=us-central1 \
  --source=. \
  --set-env-vars INSTANCE_CONNECTION_NAME="chatfluent-409007:us-central1:chatfluent-instance" \
  --set-env-vars DB_NAME="chatfluent_db" \
  --set-env-vars DB_USER="chatfluent-service-account@chatfluent-409007.iam" \
  --service-account="chatfluent-service-account@chatfluent-409007.iam.gserviceaccount.com" \
  --allow-unauthenticated
```


### ローカル実行

ローカルでのテスト環境構築のため、
設定ファイルで`production`と`development`を切り替えられるようにする。

 config.jsonで設定する。