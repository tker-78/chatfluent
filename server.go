package main

import (
	"log"
	"net/http"
	"text/template"
	"time"
)

func startServer() {
	mux := http.NewServeMux()
	mux.HandleFunc("/", index)

	server := &http.Server{
		Addr:           config.Address,
		Handler:        mux,
		ReadTimeout:    time.Duration(config.ReadTimeout * int64(time.Second)),
		WriteTimeout:   time.Duration(config.WriteTimeout * int64(time.Second)),
		MaxHeaderBytes: 1 << 20,
	}

	server.ListenAndServe()

}

func index(w http.ResponseWriter, r *http.Request) {
	// fmt.Fprintf(w, "This is the %s page", "xxxx")

	// index.htmlをパースしてテンプレートオブジェクトに変換する
	t, err := template.ParseFiles("templates/index.html")
	if err != nil {
		log.Println(err)
	}

	// nameがindexのテンプレートに対して、helloの文字列を渡して実行する。
	t.ExecuteTemplate(w, "index", "hello")
}
