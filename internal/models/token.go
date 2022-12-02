package models

import "encoding/json"

// Token is a private struct containing information
// about the user's token.
type Token struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	Login        string `json:"-"`
	UserID       string `json:"-"`
}

func (t *Token) UnmarshalFromString(payload string) error {
	return json.Unmarshal([]byte(payload), t)
}
