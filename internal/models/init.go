package models

type Init struct {
	Remote string `json:"remote"`
	User   *User
}
