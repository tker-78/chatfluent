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

type Session struct {
	Id        int
	Uuid      string
	Email     string
	UserId    string
	CreatedAt time.Time
}

// ### Sesson ###
// Todo: Create a new session for an existing user
func (user *User) CreateSession() (Session, error) {
	cmd := "INSERT INTO sessions (uuid, email, user_id, created_at) values($1, $2, $3, $4) returning id , uuid, email, user_id, created_at"

	row := DbConnection.QueryRow(cmd, createUUID(), user.Email, user.Id, time.Now())
	session := Session{}
	err := row.Scan(&session.Id, &session.Uuid, &session.Email, &session.UserId, &session.CreatedAt)

	return session, err
}

// Todo: Get the session for an existing user
func (user *User) Session() (Session, error) {
	cmd := "SELECT (uuid, email, user_id, created_at) FROM sessions WHERE user_id = $1"
	row := DbConnection.QueryRow(cmd, user.Id)
	session := Session{}
	err := row.Scan(&session.Uuid, &session.Email, &session.UserId, &session.CreatedAt)
	if err != nil {
		return Session{}, err
	}
	return session, nil
}

// Todo: Check if session is valid in the database
func (session *Session) Check() (bool, error) {
	cmd := "SELECT (id, uuid, email, user_id, created_at) FROM sessions WHERE uuid = $1"
	// stmt, err := DbConnection.Prepare(cmd)
	// if err != nil {
	// 	return false, err
	// }
	// defer stmt.Close()

	row := DbConnection.QueryRow(cmd, session.Uuid)
	row.Scan(&session.Id, &session.Uuid, &session.Email, &session.UserId, &session.CreatedAt)

	if session.Id != 0 {
		return true, nil
	} else {
		return false, nil
	}
}

// Todo: Delete session from database
func (session *Session) DeleteByUUID() error {
	cmd := "DELETE FROM sessions WHERE uuid = $1"
	stmt, err := DbConnection.Prepare(cmd)
	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(session.Uuid)
	if err != nil {
		return err
	}

	return err
}

// Get the user from the session
func (session *Session) User() (User, error) {
	cmd := "SELECT (id, uuid, name, email, created_at) FROM users WHERE ID = $1"
	stmt, err := DbConnection.Prepare(cmd)
	if err != nil {
		return User{}, err
	}
	defer stmt.Close()

	row := stmt.QueryRow(session.UserId)
	user := User{}
	row.Scan(&user.Id, &user.Uuid, &user.Name, &user.Email, &user.CreatedAt)

	return user, err
}

// Delete all sessions from database
func SessionDeleteAll() error {
	cmd := "DELETE FROM sessions"
	_, err := DbConnection.Exec(cmd)
	if err != nil {
		return err
	}
	return err
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
	rows := stmt.QueryRow(createUUID(), user.Name, user.Email, Encrypt(user.Password), time.Now())
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
	cmd := "DELETE FROM users;"

	_, err := DbConnection.Exec(cmd)
	if err != nil {
		return err
	}

	return err
}

// Todo: Get all users in the database and returns it
func Users() ([]User, error) {
	cmd := "SELECT id, uuid, name, email, password, created_at FROM users;"
	stmt, err := DbConnection.Prepare(cmd)
	if err != nil {
		return []User{}, err
	}
	defer stmt.Close()

	rows, err := stmt.Query()
	if err != nil {
		return []User{}, err
	}
	defer rows.Close()

	var users []User
	for rows.Next() {
		var user User
		// Todo: passwordは暗号化が必要
		rows.Scan(&user.Id, &user.Uuid, &user.Name, &user.Email, &user.Password, &user.CreatedAt)
		users = append(users, user)
	}

	return users, nil
}

// Todo: Get a single user given the email
func UserByEmail(email string) (User, error) {
	cmd := "select id, uuid, email, password, created_at FROM users WHERE email = $1"
	row := DbConnection.QueryRow(cmd, email)
	user := User{}
	err := row.Scan(&user.Id, &user.Uuid, &user.Email, &user.Password, &user.CreatedAt)
	if err != nil {
		log.Fatalln(err)
	}
	return user, err
}

// Todo: Get a single user given the UUID
func UserByUUID(uuid string) (User, error) {
	cmd := "SELECT id, uuid, name, email, password, created_at FROM users WHERE uuid = $1"
	stmt, err := DbConnection.Prepare(cmd)
	if err != nil {
		return User{}, err
	}
	defer stmt.Close()
	row := stmt.QueryRow(uuid)
	user := User{}
	row.Scan(&user.Id, &user.Uuid, &user.Name, &user.Email, &user.Password, &user.CreatedAt)

	return user, err

}
