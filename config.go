package main

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"os"
	"strings"

	"example.com/tker-78/chatfluent/data"
)

/*
config.goでは下記の役割を持たせています

1. config.jsonでの設定
2. Webサーバーの接続情報の定義と設定
3. logの設定
4. エラーメッセージの設定
5. ヘルパーメソッドの定義(session, generateHTMLなど)

*/

type Configuration struct {
	Address      string
	ReadTimeout  int64
	WriteTimeout int64
	Static       string
}

// アプリケーション全体で変数にアクセスできるようにグローバル宣言
var config Configuration
var logger *log.Logger

func init() {
	loadConfig()
	file, err := os.OpenFile("chatfluent.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalln("Failed to open the log file", err)
	}
	logger = log.New(file, "INFO ", log.Ldate|log.Ltime|log.Lshortfile)
}

// config.jsonからの読み出しの設定
func loadConfig() {
	file, err := os.Open("config.json")
	if err != nil {
		log.Fatalln("Failed to open config file", err)
	}
	decoder := json.NewDecoder(file)
	config = Configuration{}
	err = decoder.Decode(&config)
	if err != nil {
		log.Fatalln("Cannot get configuration from file", err)
	}

}

// エラーメッセージのページへのリダイレクト
func error_message(w http.ResponseWriter, r *http.Request, msg string) {
	url := []string{"/err?msg=", msg}
	http.Redirect(w, r, strings.Join(url, ""), 302)

}

// session
func session(w http.ResponseWriter, r *http.Request) (sess data.Session, err error) {
	cookie, err := r.Cookie("_cookie")
	if err == nil {
		sess = data.Session{Uuid: cookie.Value}
		if ok, _ := sess.Check(); !ok {
			err = errors.New("invalid session")
		}
	}
	return

}

// HTMLの生成
func generateHTML(w http.ResponseWriter, r *http.Request) {
	// todo:
}

// ログファイル
func info(args ...interface{}) {
	logger.SetPrefix("INFO")
	logger.Println(args...)
}
