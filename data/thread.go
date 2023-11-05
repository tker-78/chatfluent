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

// format the CreatedAt date to display nicely on the screen
func (thread *Thread) CreatedAtDate() string {
	return thread.CreatedAt.Format(time.RFC3339)
}

// threadの中のpostsの数を返す
func (thread *Thread) NumsReplies() int {
	cmd := "select count(*) from posts where thread_id = %1"

	rows, err := DbConnection.Query(cmd, thread.Id)
	if err != nil {
		log.Println(err)
	}

	var count int
	for rows.Next() {
		err = rows.Scan(&count)
		if err != nil {
			log.Println(err)
		}
	}
	return count

}

// 新しいthreadを作成する
func (user *User) CreateThread(topic string) (thread Thread, err error) {
	cmd := "INSER INTO threads (uuid, topic, user_id, created_at) VALUES ($1, $2, $3, $4)"
	row := DbConnection.QueryRow(cmd, createUUID(), topic, user.Id, time.Now())
	err = row.Scan(&thread.Id, &thread.Uuid, &thread.Topic, &thread.UserId, &thread.CreatedAt)
	if err != nil {
		log.Fatalln(err)
	}
	return thread, err
}

// threadに新しいpostを追加する
func (user *User) CreatePost(thread Thread, body string) (post Post, err error) {
	// todo:
	cmd := "INSER INTO posts (uuid, body, user_id, thread_id, created_at) VALUES ($1, $2, $3, $4, $5) returning (uuid, body, user_id, thread_id, created_at)"
	row := DbConnection.QueryRow(cmd, createUUID(), body, user.Id, thread.Id, time.Now())
	err = row.Scan(&post.Uuid, &post.Body, &post.UserId, &post.ThreadId, &post.CreatedAt)
	if err != nil {
		log.Fatalln(err)
	}
	return
}

// すべてのthreadを取得する
func Threads() (threads []Thread, err error) {
	cmd := "SELECT id, uuid, topc, user_id, created_at FROM threads;"
	rows, err := DbConnection.Query(cmd)
	if err != nil {
		log.Fatalln(err)
	}

	var thread Thread

	for rows.Next() {
		err = rows.Scan(&thread)
		if err != nil {
			log.Println(err)
		}
		threads = append(threads, thread)
	}
	return threads, err
}
