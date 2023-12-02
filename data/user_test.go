package data

import (
	"database/sql"
	"testing"
)

func Test_UserCreate(t *testing.T) {
	setup()
	if err := users[0].Create(); err != nil {
		t.Error(err, "Cannot create user.")
	}
	if users[0].Id == 0 {
		t.Errorf("No id or created_at in user")
	}
	u, err := UserByEmail(users[0].Email)
	if err != nil {
		t.Error(err, "User not created.")
	}
	if users[0].Email != u.Email {
		t.Errorf("User retrieved is not the same as the one created.")
	}
}

func Test_UserDelete(t *testing.T) {
	setup()
	if err := users[0].Create(); err != nil {
		t.Error(err, "Cannot create user")
	}
	if err := users[0].Delete(); err != nil {
		t.Error(err, "Cannot delete user")
	}

	_, err := UserByEmail(users[0].Email)
	if err != sql.ErrNoRows {
		t.Error(err, "User not deleted.")
	}
}

func Test_UserDeleteByEmail(t *testing.T) {
	setup()
	if err := users[0].Create(); err != nil {
		t.Error(err, "Cannot create user")
	}
	u, err := UserByEmail(users[0].Email)
	if err != nil {
		t.Errorf("cannot find the user")
	}
	if err := DeleteByEmail(u.Email); err != nil {
		t.Error(err, "cannot delete user")
	}
	_, err = UserByEmail(u.Email)
	if err != sql.ErrNoRows {
		t.Error(err, "user not deleted")
	}
}

func Test_Users(t *testing.T) {
	setup()
	for _, user := range users {
		if err := user.Create(); err != nil {
			t.Error(err, "cannot crate user")
		}
	}

	us, err := Users()
	if err != nil {
		t.Error(err, "cannot get users")
	}

	if len(us) != len(users) {
		t.Error(err, "wrong numbers of users retrieved")
	}
	if us[0].Email != users[0].Email {
		t.Error(us[0], users[0], "wrong user retrieved")
	}
}

func Test_UserUpdate(t *testing.T) {
	// todo:
	setup()
	if err := users[0].Create(); err != nil {
		t.Error(err, "cannot create user")
	}

	users[0].Name = "updatedUser"

	if err := users[0].Update(); err != nil {
		t.Error(err, "cannot update user")
	}

	u, err := UserByEmail(users[0].Email)
	if err != nil {
		t.Error(err, "cannot find the user")
	}
	if u.Name != "updatedUser" {
		t.Error(err, "user not updated")
	}

}

func Test_UserDeleteAll(t *testing.T) {
	setup()
	if err := users[0].Create(); err != nil {
		t.Error(err, "cannot create user")
	}
	if err := DeleteAllUsers(); err != nil {
		t.Error(err, "cannot delete all users")
	}
}

func Test_SessionCreate(t *testing.T) {
	setup()
	if err := users[0].Create(); err != nil {
		t.Error(err, "cannot create user")
	}
	session, err := users[0].SessionCreate()
	if err != nil {
		t.Error(err, "cannot create session")
	}
	if session.UserId != users[0].Id {
		t.Error(err, "User not linked with session")
	}
}

func Test_GetSession(t *testing.T) {
	setup()
	if err := users[0].Create(); err != nil {
		t.Error(err, "cannot create user")
	}
	session, err := users[0].SessionCreate()
	if err != nil {
		t.Error(err, "cannot create session")
	}
	sess, err := users[0].Session()
	if err != nil {
		t.Error(err, "cannot get session")
	}
	if sess.Id != session.Id {
		t.Error(err, "cannot get correct user's session")
	}
}

func Test_DeleteSession(t *testing.T) {
	setup()
	if err := users[0].Create(); err != nil {
		t.Error(err, "cannot create user")
	}
	_, err := users[0].SessionCreate()
	if err != nil {
		t.Error(err, "cannot create session")
	}

	if err = users[0].SessionDelete(); err != nil {
		t.Error(err, "cannot delete session")
	}

	_, err = users[0].Session()
	if err != sql.ErrNoRows {
		t.Error(err, "session not deleted")
	}

}
