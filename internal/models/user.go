package models

import jsoniter "github.com/json-iterator/go"

// User contains information about the user.
type User struct {
	ID           string `json:"id"`
	Login        string `json:"login"`
	Password     string `json:"password,omitempty"`
	PasswordHash string `json:"-"`
}

// Sanitize clears the password field.
func (u *User) Sanitize() {
	u.Password = ""
}

// Marshal returns the JSON encoding of user.
func (u *User) Marshal() ([]byte, error) {
	var json = jsoniter.ConfigFastest
	return json.Marshal(&u)
}
