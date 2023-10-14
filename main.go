package main

import (
	"fmt"

	"github.com/tker-78/chatfluent/data"
)

func main() {
	user := data.User{
		Name:     "yyyy",
		Email:    "yyyy@gmail.com",
		Password: "password",
	}

	user.Create()

	fmt.Println(data.Users())

}
