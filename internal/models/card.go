//nolint:gomnd
package models

import (
	"bytes"
	"encoding/gob"
	"regexp"

	validation "github.com/go-ozzo/ozzo-validation"
)

type Card struct {
	Type   string `json:"type"`
	Number string `json:"number"`
	Month  string `json:"month"`
	Year   string `json:"year"`
	CVV    string `json:"cvv"`
}

func (c *Card) Sanitize() {
	c.Type, c.Number, c.Month, c.Year, c.CVV = "unknown", "0000000000000000", "01", "1970", "0000"
}

func (c *Card) Encode() ([]byte, error) {
	var buf bytes.Buffer
	if err := gob.NewEncoder(&buf).Encode(c); err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

func (c *Card) Decode(data []byte) error {
	var buf bytes.Buffer
	buf.Write(data)

	return gob.NewDecoder(&buf).Decode(&c)
}

func (c *Card) Validate() error {
	return validation.ValidateStruct(
		c,
		validation.Field(&c.Number, validation.Required, validation.By(IsValidLuhn)),
		validation.Field(&c.Month, validation.Required, validation.Match(regexp.MustCompile("^(0?[1-9]|1[012])$"))),
		validation.Field(&c.Year, validation.Required, validation.Match(regexp.MustCompile(`^\d{2}$`))),
		validation.Field(&c.CVV, validation.Required, validation.Length(3, 4)),
	)
}
