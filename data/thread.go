package data

import (
	"log"
	"time"
)

type Thread struct {
	Id        int
	Uuid      string
	Topic     string
	UserId    int
	CreatedAt time.Time
}

type Post struct {
	Id        int
	Uuid      string
	Body      string
	UserId    int
	ThreadId  int
	CreatedAt time.Time
}

func (thread *Thread) CreatedAtDate() string {
	return thread.CreatedAt.Format("2006年01月02日 03:04pm")
}

func (post *Post) CreatedAtDate() string {
	return post.CreatedAt.Format("2006年01月02日 03:04pm")
}

// Create a new thread
func (user *User) CreateThread(topic string) (th Thread, err error) {
	if topic == "" {
		log.Println("topic should not be empty")
		return
	}
	cmd := `INSERT INTO threads (uuid, topic, user_id, created_at) 
					VALUES ($1, $2, $3, $4)
					returning id, uuid, topic, user_id, created_at`

	th = Thread{}
	err = DbConnection.QueryRow(cmd, createUUID(), topic, user.Id, time.Now()).
		Scan(&th.Id, &th.Uuid, &th.Topic, &th.UserId, &th.CreatedAt)

	if err != nil {
		log.Println(err)
		return th, err
	}
	return th, err
}

// Create a new post to a thread
func (user *User) CreatePost(th Thread, body string) (post Post, err error) {
	if body == "" {
		log.Println("post body should not be empty.")
		return
	}
	cmd := `INSERT INTO posts (uuid, body, user_id, thread_id, created_at)
					VALUES ($1, $2, $3, $4, $5)
					returning id, uuid, body, user_id, thread_id, created_at
					`
	row := DbConnection.QueryRow(cmd, createUUID(), body, user.Id, th.Id, time.Now())
	post = Post{}
	err = row.Scan(&post.Id, &post.Uuid, &post.Body, &post.UserId, &post.ThreadId, &post.CreatedAt)
	if err != nil {
		log.Println(err)
		return post, err
	}
	return post, err
}

// Get all threads in the database
func Threads() ([]Thread, error) {
	threads := []Thread{}
	rows, err := DbConnection.Query("SELECT * FROM threads")
	if err != nil {
		log.Println(err)
		return threads, err
	}
	for rows.Next() {
		var thread Thread
		rows.Scan(&thread.Id, &thread.Uuid, &thread.Topic, &thread.UserId, &thread.CreatedAt)
		threads = append(threads, thread)
	}
	return threads, err
}

func ThreadByUuid(uuid string) (thread Thread, err error) {
	cmd := "SELECT id, uuid, topic, user_id, created_at FROM threads WHERE uuid = $1"
	thread = Thread{}
	err = DbConnection.QueryRow(cmd, uuid).Scan(&thread.Id, &thread.Uuid, &thread.Topic, &thread.UserId, &thread.CreatedAt)
	return
}

func PostByUuid(uuid string) (post Post, err error) {
	cmd := "SELECT id, uuid, body, user_id, thread_id, created_at FROM posts WHERE uuid = $1"
	post = Post{}
	err = DbConnection.QueryRow(cmd, uuid).Scan(&post.Id, &post.Uuid, &post.Body, &post.UserId, &post.ThreadId, &post.CreatedAt)
	return
}

// delete all threads
func DeleteAllThreads() error {
	cmd := "DELETE FROM threads"
	_, err := DbConnection.Exec(cmd)
	if err != nil {
		log.Println(err)
		return err
	}
	return err
}

// threadのUserを返す
func (thread *Thread) User() (user *User) {
	cmd := "SELECT id, uuid, name, email, created_at FROM users WHERE id = $1"
	user = new(User)
	DbConnection.QueryRow(cmd, thread.UserId).Scan(&user.Id, &user.Uuid, &user.Name, &user.Email, &user.CreatedAt)
	return
}

// postのUserを返す
func (post *Post) User() (user *User) {
	cmd := "SELECT id, uuid, name, email, created_at FROM users WHERE id = $1"
	user = new(User)
	DbConnection.QueryRow(cmd, post.UserId).Scan(&user.Id, &user.Uuid, &user.Name, &user.Email, &user.CreatedAt)
	return
}

// threadに属する全てのPostを返す
func (thread *Thread) Posts() ([]Post, error) {
	cmd := "SELECT id, uuid, body, user_id, thread_id, created_at FROM posts WHERE thread_id = $1"
	posts := []Post{}
	rows, err := DbConnection.Query(cmd, thread.Id)
	if err != nil {
		log.Println(err, "cannot get query for posts")
	}
	defer rows.Close()
	for rows.Next() {
		post := Post{}
		rows.Scan(&post.Id, &post.Uuid, &post.Body, &post.UserId, &post.ThreadId, &post.CreatedAt)
		posts = append(posts, post)
	}
	return posts, err
}

// delete thread
func (thread *Thread) Delete() error {
	_, err := DbConnection.Exec("DELETE FROM threads WHERE id = $1", thread.Id)
	return err
}

func (post *Post) Delete() error {
	_, err := DbConnection.Exec("DELETE FROM posts WHERE id = $1", post.Id)
	return err
}

// delete all posts
func DeleteAllPosts() error {
	cmd := "DELETE FROM posts"
	_, err := DbConnection.Exec(cmd)
	if err != nil {
		log.Println(err)
		return err
	}
	return err
}
