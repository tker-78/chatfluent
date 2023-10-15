package main

import (
	"encoding/json"
	"fmt"
	"html/template"
	"log"
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
	config := Configulation{}
	err = decoder.Decode(&config)
	if err != nil {
		log.Fatalln("Cannot get configlation from file", err)
	}
}

// parse HTML templates
// 各route_authから呼び出す
func parseTemplateFiles(filenames ...string) (t *template.Template) {
	var files []string
	t = template.New("layout")
	for _, file := range filenames {
		files = append(files, fmt.Sprintf("templates/%s.html", file))
	}

	t = template.Must(t.ParseFiles(files...))
	return
}
