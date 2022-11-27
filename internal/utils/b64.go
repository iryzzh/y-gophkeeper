package utils

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"strings"
)

// ToBase64 encodes the source in base64.
func ToBase64(src []byte) []byte {
	return []byte(base64.StdEncoding.EncodeToString(src))
}

func FromBase64(src []byte) []byte {
	dec, _ := base64.StdEncoding.DecodeString(string(src))

	return dec
}

func FromBinaryBase64(src []byte, dst interface{}) error {
	return json.NewDecoder(base64.NewDecoder(base64.StdEncoding, strings.NewReader(string(src)))).Decode(dst)
}

// AppendBase64 appends base64-encoded src to dst and returns the dst.
func AppendBase64(e *base64.Encoding, dst, src []byte) []byte {
	buf := make([]byte, e.EncodedLen(len(src)))
	e.Encode(buf, src)
	return append(dst, buf...)
}

// InterfaceToBase64 encodes the source in base64.
func InterfaceToBase64(src interface{}) ([]byte, error) {
	var buf bytes.Buffer
	enc := base64.NewEncoder(base64.StdEncoding, &buf)
	err := json.NewEncoder(enc).Encode(src)
	if err != nil {
		return nil, err
	}
	_ = enc.Close()

	return buf.Bytes(), nil
}
