package main

import (
	"fmt"
	"html/template"
	"net/http"
)

func main() {

	//マルチプレクサを使ってリクエストを受け付ける
	mux := http.NewServeMux()

	mux.HandleFunc("/", home)
	mux.HandleFunc("/login", login)
	http.ListenAndServe("0.0.0.0:8080", mux)

}

// Handlerの定義
func home(w http.ResponseWriter, r *http.Request) {
	// fmt.Fprintf(w, "hello world")
	t, _ := template.ParseFiles("templates/layout.html", "templates/navbar.html", "templates/home.html")
	t.ExecuteTemplate(w, "layout", "helloooo!!!")

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
