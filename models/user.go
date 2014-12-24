package models

import (
	"database/sql"

	"github.com/albrow/go-data-parser"
)

type User struct {
	Id             int    `json:"id"`
	Name           string `json:"name"`
	HashedPassword string `json:"hashedPassword"`
	Email          string `json:"email"`
}

func (u *User) Save() error {
	return Db.Insert(u)
}

// IsUserUnique checks if email and name are unique, i.e. not already in the database.
// If they are not unique, adds an error to val with a detailed message. An error will
// be returned if there was a problem connecting to the database.
func ValidateUserUnique(val *data.Validator, email string, name string) (err error) {
	values := map[string]string{"email": email, "name": name}
	return mValidateUnique(val, "users", values, "That %s is already taken.")
}

func FindUserByEmail(email string) (*User, error) {
	u := new(User)
	if err := Db.SelectOne(u, "SELECT * FROM users WHERE email=$1", email); err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		} else {
			return nil, err
		}
	} else {
		return u, nil
	}
}

func FindUserById(id int32) (*User, error) {
	u := new(User)
	if err := Db.SelectOne(u, "SELECT * FROM users WHERE id=$1", id); err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		} else {
			return nil, err
		}
	} else {
		return u, nil
	}
}
