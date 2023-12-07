package data

import "testing"

func Test_CreateThread(t *testing.T) {
	setup()
	if err := users[0].Create(); err != nil {
		t.Error(err, "cannot create user.")
	}
	th, err := users[0].CreateThread("a new thread")
	if err != nil {
		t.Error(err, "cannot create thread.")
	}
	if th.UserId != users[0].Id {
		t.Error(err, "cannot create thread.")
	}
}

func Test_CreatePost(t *testing.T) {
	setup()
	if err := users[0].Create(); err != nil {
		t.Error(err, "cannot create user")
	}
	th, err := users[0].CreateThread("a new thread")
	if err != nil {
		t.Error(err, "cannot create user")
	}
	post, err := users[0].CreatePost(th, "a new post")
	if err != nil {
		t.Error(err, "cannot create post")
	}
	if post.UserId != users[0].Id {
		t.Error(err, "cannot create for correct user")
	}
	if post.ThreadId != th.Id {
		t.Error(err, "cannot create for correct thread")
	}

}
