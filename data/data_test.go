package data

// test data

var users = []User{
	{
		Name:     "takuya kinoshita",
		Email:    "aaa@gmail.com",
		Password: "password",
	},
	{
		Name:     "tomomo yakayama",
		Email:    "bbb@gmail.com",
		Password: "password",
	},
}

func setup() {
	DeleteAllUsers()
	DeleteAllSessions()
}
