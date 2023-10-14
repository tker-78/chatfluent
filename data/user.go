package data

import (
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

type Session struct {
	Id       int
	Uuid     string
	Email    string
	UserId   string
	CreateAt time.Time
}

// ### Sesson ###
// Create a new session for an existing user
func (user *User) CreateSession() (Session, error) {

	return Session{}, nil
}

// Get the session for an existing user
func (user *User) Session() (Session, error) {

	return Session{}, nil
}

// Check if session is valid in the database
func (session *Session) Check() (bool, error) {

	return false, nil
}

// Delete session from database
func (session *Session) DeleteByUUID() error {

	return nil
}

// Get the user from the session
func (session *Session) User() (User, error) {

	return User{}, nil
}

// Delete all sessions from database
func SessionDeleteAll() error {

	return nil
}

// ### User ###
// Create user
func (user *User) Create() error {
	// Idは自動で付与されるので、指定する必要なし
	cmd := "INSERT INTO users (uuid, name, email, password, created_at) values ($1, $2, $3, $4, $5)"
	stmt, err := DbConnection.Prepare(cmd)
	if err != nil {
		return err
	}
	defer stmt.Close()

	// UUIDを生成して、データベースに保存する
	rows := stmt.QueryRow(createUUID(), user.Name, user.Email, user.Password, time.Now())
	// user構造体にスキャンする
	err = rows.Scan(&user.Uuid, &user.Name, &user.Email, &user.Password, &user.CreatedAt)

	return err
}

// Delete user from database
func (user *User) Delete() error {
	cmd := "DELETE FROM users WHERE uuid = $1"
	_, err := DbConnection.Exec(cmd, user.Id)
	if err != nil {
		return err
	}
	return err
}

// update user information in the database
func (user *User) Update() error {
	cmd := "UPDATE users SET name = $2, email = $3, where id = $1"
	stmt, err := DbConnection.Prepare(cmd)
	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(user.Id, user.Name, user.Email)

	return err
}

// Todo: Delete all users from database
func UserDeleteAll() error {

	return nil
}

// Todo: Get all users in the database and returns it
func Users() ([]User, error) {

	return []User{}, nil
}

// Todo: Get a single user given the email
func UserByEmail() (User, error) {

	return User{}, nil
}

// Todo: Get a single user given the UUID
func UserByUUID() (User, error) {

	return User{}, nil
}
