package data

import (
	"database/sql"
	_ "database/sql"
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
