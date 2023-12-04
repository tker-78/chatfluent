package main

import (
	"log"
	"net/http"
	"text/template"
	"time"

	"example.com/tker-78/chatfluent/data"
)

func startServer() {
	mux := http.NewServeMux()
	files := http.FileServer(http.Dir(config.Static))
	mux.Handle("/static/", http.StripPrefix("/static/", files))
	mux.HandleFunc("/", index)

	mux.HandleFunc("/login", login)
	mux.HandleFunc("/authenticate", authenticate)
	mux.HandleFunc("/signup", signup)
	mux.HandleFunc("/signup_account", signupAccount)

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
	t, err := template.ParseFiles("templates/layout.html", "templates/index.html", "templates/navbar.html")
	if err != nil {
		log.Println(err)
	}

	users, err := data.Users()
	if err != nil {
		log.Println(err)
	}

	// nameがindexのテンプレートに対して、helloの文字列を渡して実行する。
	t.ExecuteTemplate(w, "layout", users)
}

func login(w http.ResponseWriter, r *http.Request) {
	t, err := template.ParseFiles("templates/layout.html", "templates/navbar.html", "templates/login.html")
	if err != nil {
		log.Println(err)
	}
	t.ExecuteTemplate(w, "layout", nil)
}

func authenticate(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		log.Println(err)
	}
	user, err := data.UserByEmail(r.PostFormValue("email"))
	if err != nil {
		log.Println(err, "cannot find user")
	}
	if user.Password == data.Encrypt(r.PostFormValue("password")) {
		session, err := user.SessionCreate()
		if err != nil {
			log.Println(err, "cannot create session")
		}
		cookie := http.Cookie{
			Name:     "_cookie",
			Value:    session.Uuid,
			HttpOnly: true,
		}
		http.SetCookie(w, &cookie)
		http.Redirect(w, r, "/", http.StatusFound)
		log.Println("login successful.")
	} else {
		http.Redirect(w, r, "/login", http.StatusFound)
	}
}

// userの新規登録
// GET
func signup(w http.ResponseWriter, r *http.Request) {
	t, err := template.ParseFiles("templates/layout.html", "templates/navbar.html", "templates/signup.html")
	if err != nil {
		log.Println(err)
	}
	t.ExecuteTemplate(w, "layout", nil)
}

// POST
func signupAccount(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		log.Println(err)
	}
	user := data.User{
		Name:     r.PostFormValue("name"),
		Email:    r.PostFormValue("email"),
		Password: r.PostFormValue("password"),
	}
	if err := user.Create(); err != nil {
		log.Println(err)
	}
	http.Redirect(w, r, "/login", http.StatusFound)
	log.Println("signup successful, please login.")
}
