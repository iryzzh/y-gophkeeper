package models

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
)

type Card struct {
	Type        string `json:"type"`
	Number      string `json:"number"`
	ExpiryMonth int    `json:"expiryMonth"`
	ExpiryYear  int    `json:"expiryYear"`
	CCV         int    `json:"ccv"`
}

func (c *Card) EncodeToBase64() ([]byte, error) {
	var buf bytes.Buffer
	enc := base64.NewEncoder(base64.StdEncoding, &buf)
	err := json.NewEncoder(enc).Encode(c)
	if err != nil {
		return nil, err
	}
	_ = enc.Close()

	return buf.Bytes(), nil
}
