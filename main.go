package main

import (
	"html/template"
	"log"
	"net/http"

	"github.com/tker-78/chatfluent/data"
)

func main() {

	//マルチプレクサを使ってリクエストを受け付ける
	mux := http.NewServeMux()

	mux.HandleFunc("/", home)
	mux.HandleFunc("/login", login)
	mux.HandleFunc("/authenticate", authenticate)
	mux.HandleFunc("/signup", signup)
	mux.HandleFunc("/signup_account", signupAccount)
	mux.HandleFunc("/threads", threads)
	http.ListenAndServe("0.0.0.0:5050", mux)

}

// Handlerの定義
func home(w http.ResponseWriter, r *http.Request) {
	// t, _ := template.ParseFiles("templates/layout.html", "templates/public.home.html")
	// t.ExecuteTemplate(w, "layout", nil)

	// login状態で表示するページを変える
	_, err := session(w, r)
	if err != nil {
		log.Println(err, "error in checking session")
		t, _ := template.ParseFiles("templates/layout.html", "templates/public.home.html")
		t.ExecuteTemplate(w, "layout", nil)
	} else {
		t, _ := template.ParseFiles("templates/layout.html", "templates/private.home.html")
		t.ExecuteTemplate(w, "layout", nil)
	}

}

// GET /login
func login(w http.ResponseWriter, r *http.Request) {
	t, _ := template.ParseFiles("templates/layout.html", "templates/login.html")
	t.ExecuteTemplate(w, "layout", nil)
}

// GET /signup
// signupページの表示
func signup(w http.ResponseWriter, r *http.Request) {
	t, _ := template.ParseFiles("templates/layout.html", "templates/signup.html")
	t.ExecuteTemplate(w, "layout", nil)

}

// POST /signup_account
func signupAccount(w http.ResponseWriter, r *http.Request) {
	// todo:
	err := r.ParseForm()
	if err != nil {
		log.Println(err)
	}
	user := data.User{
		Name:     r.PostFormValue("name"),
		Email:    r.PostFormValue("email"),
		Password: r.PostFormValue("password"),
	}
	err = user.Create()
	if err != nil {
		log.Println(err, "cannot create user")
	}
	http.Redirect(w, r, "/login", http.StatusFound)

}

// POST /authenticate
func authenticate(w http.ResponseWriter, r *http.Request) {
	// todo: session管理がうまくいかないので、今はパスする(2023/11/17 kinoshita)
	err := r.ParseForm()
	if err != nil {
		log.Println(err)
	}

	user, err := data.UserByEmail(r.PostFormValue("email"))
	if err != nil {
		log.Println(err, "cannot find the user")
	}
	if user.Password == data.Encrypt(r.PostFormValue("password")) {
		session, err := user.CreateSession()
		if err != nil {
			log.Println(err, "cannot create session")
		}
		cookie := http.Cookie{
			Name:     "_cookie",
			Value:    session.Uuid,
			HttpOnly: true,
		}
		http.SetCookie(w, &cookie)
		log.Println("ログインしました。")
		log.Println(r.Cookie("_cookie"))
		http.Redirect(w, r, "/", http.StatusFound)
	} else {
		log.Println("ログインできませんでした。")
		http.Redirect(w, r, "/login", http.StatusFound)
	}
}

// GET /logout
func logout(w http.ResponseWriter, r *http.Request) {

}

// threadsの一覧ページ
func threads(w http.ResponseWriter, r *http.Request) {
	t, _ := template.ParseFiles("templates/layout.html", "templates/threads.html")
	t.ExecuteTemplate(w, "layout", nil) // ここにthread構造体を渡してレンダリングする
}
