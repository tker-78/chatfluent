package main

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"os"

	"github.com/tker-78/chatfluent/data"
)

type Configulation struct {
	Adress       string
	ReadTimeout  int64
	WriteTimeout int64
	Static       string
}

var config Configulation

func init() {
	loadConfig()
	// Todo: ログの書き出し

}

// config.jsonの読み出し
// config structに情報を格納する
// Fatalするので、errorの返り値は必要ない
func loadConfig() {
	file, err := os.Open("config.json")
	if err != nil {
		log.Fatalln("Cannot open config file", err)
	}
	decoder := json.NewDecoder(file)
	config = Configulation{}
	err = decoder.Decode(&config)
	if err != nil {
		log.Fatalln("Cannot get configlation from file", err)
	}
}

// Userのログイン判定
func session(w http.ResponseWriter, r *http.Request) (sess data.Session, err error) {
	cookie, err := r.Cookie("_cookie")
	if err == nil {
		sess := data.Session{Uuid: cookie.Value}
		if ok, _ := sess.Check(); !ok {
			err = errors.New("invalid session")
		}
	}
	return
}
