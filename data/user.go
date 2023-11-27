package data

import (
	"log"
	"time"
)

type User struct {
	Id        int
	Uuid      string
	Name      string
	Email     string
	Password  string
	CreatedAt time.Time
}

// Create a new user
// not yet tested
func (user *User) Create() error {
	cmd := `INSERT INTO users 
					(uuid, name, email, password, created_at) 
					VALUES ($1, $2, $3, $4, $5)
					returning id, uuid, created_at` // オブジェクトに値を返して格納する

	// _, err := DbConnection.Exec(cmd, createUUID(), user.Name, user.Email, Encrypt(user.Password), time.Now())
	row := DbConnection.QueryRow(cmd, createUUID(), user.Name, user.Email, Encrypt(user.Password), time.Now())
	err := row.Scan(&user.Id, &user.Uuid, &user.CreatedAt)
	return err
}

// Get a single user given the email
// not yet tested
func UserByEmail(email string) (user User, err error) {
	cmd := `SELECT id, uuid, name, email, password, created_at
					FROM users
					WHERE email = $1`
	row := DbConnection.QueryRow(cmd, email)
	user = User{}
	err = row.Scan(&user.Id, &user.Uuid, &user.Name, &user.Email, &user.Password, &user.CreatedAt)
	return
}

// Get all users from database
// not yet tested
func Users() ([]User, error) {
	cmd := "SELECT id, uuid, name, email, password, created_at FROM users"
	rows, err := DbConnection.Query(cmd)
	if err != nil {
		log.Fatalln(err)
	}
	users := []User{}

	for rows.Next() {
		user := User{}
		err = rows.Scan(&user.Id, &user.Uuid, &user.Name, &user.Email, &user.Password, &user.CreatedAt)
		if err != nil {
			return users, err
		}
		users = append(users, user)
	}
	return users, err
}

// Update user information in the database
// userの更新を保存する
// not yet tested
func (user *User) Update() error {
	cmd := "UPDATE users set name= $2, email = $3 WHERE id = $1"
	_, err := DbConnection.Exec(cmd, user.Id, user.Name, user.Email)
	if err != nil {
		log.Fatalln(err)
	}
	return err
}

// Delete a single user given the email
// not yet tested
func DeleteByEmail(email string) error {
	cmd := "SELECT id, uuid, name, email, password, created_at FROM users WHERE email = $1"
	row := DbConnection.QueryRow(cmd, email)
	user := User{}
	row.Scan(&user.Id, &user.Uuid, &user.Name, &user.Email, &user.Password, &user.CreatedAt)

	cmd2 := "DELETE FROM users WHERE id = $1"
	_, err := DbConnection.Exec(cmd2, user.Id)
	if err != nil {
		log.Fatalln(err)
	}
	return err
}

// Delete all users from database
// not yet tested
func DeleteAllUsers() error {
	cmd := "DELETE FROM users"
	_, err := DbConnection.Exec(cmd)
	return err
}
