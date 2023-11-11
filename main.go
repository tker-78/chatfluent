package main

import (
	"html/template"
	"net/http"
)

func main() {

	//マルチプレクサを使ってリクエストを受け付ける
	mux := http.NewServeMux()

	mux.HandleFunc("/", home)
	mux.HandleFunc("/login", login)
	mux.HandleFunc("/signup", signupAccount)
	http.ListenAndServe("0.0.0.0:8080", mux)

}

// Handlerの定義
func home(w http.ResponseWriter, r *http.Request) {
	// fmt.Fprintf(w, "hello world")
	t, _ := template.ParseFiles("templates/layout.html", "templates/home.html")
	t.ExecuteTemplate(w, "layout", nil)

}

// GET /login
func login(w http.ResponseWriter, r *http.Request) {
	t, _ := template.ParseFiles("templates/layout.html", "templates/login.html")
	t.ExecuteTemplate(w, "layout", nil)
}

// POST /signup
func signupAccount(w http.ResponseWriter, r *http.Request) {
	t, _ := template.ParseFiles("templates/layout.html", "templates/signup.html")
	t.ExecuteTemplate(w, "layout", nil)

}

// POST /authenticate
func authenticate(w http.ResponseWriter, r *http.Request) {

}

// GET /logout
func logout(w http.ResponseWriter, r *http.Request) {

}
