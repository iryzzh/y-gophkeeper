package models

import "encoding/base64"

type Entry struct {
	Name      string `json:"name"`
	Value     string `json:"value"`
	EntryType string `json:"value_type"`
}

const (
	EntryTypeText  = "text"
	EntryTypeFile  = "file"
	EntryTypeImage = "image"
	EntryTypeCard  = "card"
)

func (e *Entry) EncodeBytes() []byte {
	buf := make([]byte, base64.StdEncoding.EncodedLen(len(e.Value)))
	base64.StdEncoding.Encode(buf, []byte(e.Value))

	return buf
}
