package main

import (
	"encoding/json"
	"log"
	"net/http"
	"os"
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

// parse template files
// Todo: will make if it is necessary to be made.

// generateHTML

func generateHTML(w http.ResponseWriter, data interface{}, filenames ...string) {
	var files []string

}
