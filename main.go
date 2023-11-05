package main

import (
	"fmt"
	"html/template"
	"net/http"
)

func main() {

	//マルチプレクサを使ってリクエストを受け付ける
	mux := http.NewServeMux()

	mux.HandleFunc("/", index)
	mux.HandleFunc("/login", login)
	http.ListenAndServe("0.0.0.0:8080", mux)

}

// Handlerの定義
func index(w http.ResponseWriter, r *http.Request) {
	// fmt.Fprintf(w, "hello world")
	t, _ := template.ParseFiles("templates/index.html")
	t.Execute(w, "helloooo!!!")

}

// GET /login
func login(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "login")
}

// POST /signup
func signupAccount(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "signup")

}

// POST /authenticate
func authenticate(w http.ResponseWriter, r *http.Request) {

}

// GET /logout
func logout(w http.ResponseWriter, r *http.Request) {

}
