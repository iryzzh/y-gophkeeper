package models

// User contains information about the user.
type User struct {
	ID           string `json:"id"`
	Login        string `json:"login"`
	Password     string `json:"password,omitempty"`
	PasswordHash string `json:"-"`
}
