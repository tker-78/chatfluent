package main

import (
	"log"

	"github.com/tker-78/chatfluent/data"
)

func main() {
	user := data.User{
		Name:     "kktak",
		Email:    "kktak02@gmail.com",
		Password: "password",
	}

	err := user.Create()
	if err != nil {
		log.Println(err)
	}

}
