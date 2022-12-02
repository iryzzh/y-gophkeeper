package file

import (
	"encoding/base64"
	"fmt"
	"os"

	"github.com/iryzzh/y-gophkeeper/internal/config"
)

// Encode encodes the received file in base64.
func Encode(path string) ([]byte, error) {
	stat, err := os.Stat(path)
	if err != nil {
		return nil, err
	}

	if stat.Size() > config.MaxFileSize {
		return nil, fmt.Errorf("file is too big")
	}

	bytes, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	buf := make([]byte, base64.StdEncoding.EncodedLen(len(bytes)))
	base64.StdEncoding.Encode(buf, bytes)

	return buf, nil
}

// Decode decodes the received source from base64 and
// saves it to the specified file.
func Decode(src []byte, dest string) error {
	buf := make([]byte, base64.StdEncoding.DecodedLen(len(src)))
	_, err := base64.StdEncoding.Decode(buf, src)
	if err != nil {
		return err
	}

	file, err := os.Create(dest)
	if err != nil {
		return err
	}
	defer func() { _ = file.Close() }()

	if _, err = file.Write(buf); err != nil {
		return err
	}

	return file.Sync()
}
