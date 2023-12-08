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

	mux.HandleFunc("/top", top)
	mux.HandleFunc("/", index)

	mux.HandleFunc("/login", login)
	mux.HandleFunc("/authenticate", authenticate)
	mux.HandleFunc("/signup", signup)
	mux.HandleFunc("/signup_account", signupAccount)
	mux.HandleFunc("/logout", logout)
	mux.HandleFunc("/thread/read", threadRead)

	server := &http.Server{
		Addr:           config.Address,
		Handler:        mux,
		ReadTimeout:    time.Duration(config.ReadTimeout * int64(time.Second)),
		WriteTimeout:   time.Duration(config.WriteTimeout * int64(time.Second)),
		MaxHeaderBytes: 1 << 20,
	}

	server.ListenAndServe()

}

func top(w http.ResponseWriter, r *http.Request) {
	// loginしている場合は、indexページにリダイレクトする
	if _, err := session(w, r); err == nil {
		http.Redirect(w, r, "/", http.StatusFound)
		log.Println("redirected to index page.")
	}
	t, err := template.ParseFiles("templates/layout.html", "templates/public.navbar.html", "templates/top.html")
	if err != nil {
		log.Println(err)
	}
	t.ExecuteTemplate(w, "layout", nil)
}

func index(w http.ResponseWriter, r *http.Request) {
	// loginしていない場合は、topページにリダイレクトする。
	_, err := session(w, r)
	if err != nil {
		http.Redirect(w, r, "/top", http.StatusFound)
		log.Println("redirected to top page.")
	}

	// index.htmlをパースしてテンプレートオブジェクトに変換する
	t, err := template.ParseFiles("templates/layout.html", "templates/index.html", "templates/private.navbar.html")
	if err != nil {
		log.Println(err)
	}

	threads, err := data.Threads()

	// nameがindexのテンプレートに対して、helloの文字列を渡して実行する。
	t.ExecuteTemplate(w, "layout", threads)
}

func login(w http.ResponseWriter, r *http.Request) {
	t, err := template.ParseFiles("templates/layout.html", "templates/public.navbar.html", "templates/login.html")
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
	t, err := template.ParseFiles("templates/layout.html", "templates/public.navbar.html", "templates/signup.html")
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

// logout
func logout(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie("_cookie")
	if err != http.ErrNoCookie {
		log.Println(err)
		session := data.Session{Uuid: cookie.Value}
		session.DeleteByUuid()
	}
	http.Redirect(w, r, "/top", http.StatusFound)
}

// threadRead
func threadRead(w http.ResponseWriter, r *http.Request) {
	vals := r.URL.Query()
	uuid := vals.Get("id")
	thread, err := data.ThreadByUuid(uuid)
	if err != nil {
		log.Println(err)
	}
	_, err = session(w, r)
	if err != nil {
		t, err := template.ParseFiles("templates/layout.html", "templates/public.navbar.html", "templates/public.thread.html")
		if err != nil {
			log.Println(err)
		}
		t.ExecuteTemplate(w, "layout", thread)
	} else {
		t, err := template.ParseFiles("templates/layout.html", "templates/private.navbar.html", "templates/private.thread.html")
		if err != nil {
			log.Println(err)
		}
		t.ExecuteTemplate(w, "layout", thread)
	}

}
