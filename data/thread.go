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

// Create a new thread
func (user *User) CreateThread(topic string) (Thread, error) {
	cmd := `INSERT INTO threads (uuid, topic, user_id, created_at) 
					VALUES ($1, $2, $3, $4)
					returning id, uuid, topic, user_id, created_at`

	th := Thread{}
	row := DbConnection.QueryRow(cmd, createUUID(), topic, user.Id, time.Now())
	err := row.Scan(&th.Id, &th.Uuid, &th.Topic, &th.UserId, &th.CreatedAt)
	if err != nil {
		log.Println(err)
		return th, err
	}
	return th, err
}

// Create a new post to a thread
func (user *User) CreatePost(th Thread, body string) (Post, error) {
	cmd := `INSERT INTO posts (uuid, body, user_id, thread_id, created_at)
					VALUES ($1, $2, $3, $4, $5)
					returning id, uuid, body, user_id, thread_id, created_at
					`
	row := DbConnection.QueryRow(cmd, createUUID(), body, user.Id, th.Id, time.Now())
	post := Post{}
	err := row.Scan(&post.Id, &post.Uuid, &post.Body, &post.UserId, &post.ThreadId, &post.CreatedAt)
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