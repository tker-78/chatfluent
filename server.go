package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"text/template"
	"time"

	"example.com/tker-78/chatfluent/data"
	"github.com/joho/godotenv"
)

func startServer() {
	mux := http.NewServeMux()
	files := http.FileServer(http.Dir(Config.Static))
	mux.Handle("/static/", http.StripPrefix("/static/", files))

	mux.HandleFunc("/top", top)
	mux.HandleFunc("/", index)

	mux.HandleFunc("/login", login)
	mux.HandleFunc("/authenticate", authenticate)
	mux.HandleFunc("/signup", signup)
	mux.HandleFunc("/signup_account", signupAccount)
	mux.HandleFunc("/logout", logout)
	mux.HandleFunc("/thread/read", threadRead)
	mux.HandleFunc("/thread/new", threadNew)
	mux.HandleFunc("/thread/create", threadCreate)
	mux.HandleFunc("/thread/post", threadPost)
	mux.HandleFunc("/thread/delete", threadDelete)
	mux.HandleFunc("/post/delete", postDelete)

	server := &http.Server{
		Addr:           portNumber(),
		Handler:        mux,
		ReadTimeout:    time.Duration(Config.ReadTimeout * int64(time.Second)),
		WriteTimeout:   time.Duration(Config.WriteTimeout * int64(time.Second)),
		MaxHeaderBytes: 1 << 20,
	}

	fmt.Println(server.ListenAndServe())

}

func portNumber() (port string) {
	err := godotenv.Load()
	if err != nil {
		log.Println(err, "cannot read .env file")
	}
	if os.Getenv("environment") == "production" {
		port = ":" + os.Getenv("PORT")
	} else if os.Getenv("environment") == "development" {
		port = fmt.Sprintf(":%s", os.Getenv("localport"))
	}
	return

}

func top(w http.ResponseWriter, r *http.Request) {
	// loginしている場合は、indexページにリダイレクトする
	if _, err := session(w, r); err == nil {
		http.Redirect(w, r, "/", http.StatusFound)
		log.Println("redirected to index page.")
	} else {
		t, err := template.ParseFiles("templates/layout.html", "templates/public.navbar.html", "templates/top.html")
		if err != nil {
			log.Println(err)
		}
		t.ExecuteTemplate(w, "layout", nil)
	}
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
	t, err := template.ParseFiles("templates/login.layout.html", "templates/public.navbar.html", "templates/login.html")
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
	th, err := data.ThreadByUuid(uuid)
	if err != nil {
		log.Println(err)
	}
	_, err = session(w, r)
	if err != nil {
		t, err := template.ParseFiles("templates/layout.html", "templates/public.navbar.html", "templates/public.thread.html")
		if err != nil {
			log.Println(err)
		}
		t.ExecuteTemplate(w, "layout", &th)
	} else {
		t, err := template.ParseFiles("templates/layout.html", "templates/private.thread.html", "templates/private.navbar.html")
		if err != nil {
			log.Println(err)
		}
		t.ExecuteTemplate(w, "layout", &th)
	}

}

// threadNew
func threadNew(w http.ResponseWriter, r *http.Request) {
	// loginユーザーのみが新規スレッドを作成できる
	_, err := session(w, r)
	if err != nil {
		http.Redirect(w, r, "/login", http.StatusFound)
	} else {
		t, err := template.ParseFiles("templates/layout.html", "templates/private.navbar.html", "templates/new.thread.html")
		if err != nil {
			log.Println(err)
		}
		t.ExecuteTemplate(w, "layout", nil)
	}

}

func threadCreate(w http.ResponseWriter, r *http.Request) {
	sess, err := session(w, r)
	if err != nil {
		http.Redirect(w, r, "/login", http.StatusFound)
	}
	user, err := sess.User()
	if err != nil {
		log.Println(err)
	}
	topic := r.PostFormValue("topic")
	if _, err := user.CreateThread(topic); err != nil {
		log.Println(err)
	}
	http.Redirect(w, r, "/", http.StatusFound)

}

func threadPost(w http.ResponseWriter, r *http.Request) {
	sess, err := session(w, r)
	if err != nil {
		http.Redirect(w, r, "/login", http.StatusFound)
	}
	user, err := sess.User()
	if err != nil {
		log.Println(err, "cannot find user")
	}

	body := r.PostFormValue("body")
	uuid := r.PostFormValue("uuid")

	thread, err := data.ThreadByUuid(uuid)
	if err != nil {
		log.Println(err, "incorrect uuid")
		http.Redirect(w, r, "/", http.StatusFound)
	}

	if _, err := user.CreatePost(thread, body); err != nil {
		log.Println(err, "cannot create a post")
	} else {
		url := fmt.Sprintf("/thread/read?id=%s", uuid)
		http.Redirect(w, r, url, http.StatusFound)
	}

}

func threadDelete(w http.ResponseWriter, r *http.Request) {
	uuid := r.PostFormValue("uuid")
	thread, err := data.ThreadByUuid(uuid)
	if err != nil {
		log.Println(err)
	}
	err = thread.Delete()
	if err != nil {
		log.Println(err)
	}
	http.Redirect(w, r, "/", http.StatusFound)

}

func postDelete(w http.ResponseWriter, r *http.Request) {
	uuid := r.PostFormValue("uuid")
	tuuid := r.PostFormValue("tuuid")
	post, err := data.PostByUuid(uuid)
	if err != nil {
		log.Println(err)
	}
	err = post.Delete()
	if err != nil {
		log.Println(err)
	}
	url := fmt.Sprintf("/thread/read?id=%s", tuuid)
	http.Redirect(w, r, url, http.StatusFound)
}
