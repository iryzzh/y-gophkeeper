//nolint:gochecknoglobals
package rand

import (
	cryptoRand "crypto/rand"
	"io"
	mathRand "math/rand"
	"unsafe"
)

// Reader is the standard crypto/rand.Reader with added buffering.
var Reader = defaultSecureSource

func Read(p []byte) (int, error) {
	return io.ReadFull(defaultSecureSource, p)
}

// randomCharset contains the characters that can make up a rand.String().
const randomCharset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ1234567890"

const (
	letterIdxBits = 6                    // 6 bits to represent a letter index
	letterIdxMask = 1<<letterIdxBits - 1 // All 1-bits, as many as letterIdxBits
	letterIdxMax  = 63 / letterIdxBits   // # of letter indices fitting in 63 bits
)

var (
	// defaultSecureSource is a concurrency-safe, cryptographically secure
	// math/rand.Reader.
	defaultSecureSource = newSecureSource()

	// defaultSecureRand is a math/rand.Rand based on the secure source.
	defaultSecureRand = mathRand.New(defaultSecureSource)
)

// randStringBytesMaskImprSrcUnsafe generates a random string of a given length.
// https://stackoverflow.com/questions/22892120/how-to-generate-a-random-string-of-a-fixed-length-in-go/31832326#31832326
func randStringBytesMaskImprSrcUnsafe(n int, strSet string) string {
	b := make([]byte, n)
	// A Reader.Int63() generates 63 random bits, enough for letterIdxMax characters!
	for i, cache, remain := n-1, Reader.Int63(), letterIdxMax; i >= 0; {
		if remain == 0 {
			cache, remain = Reader.Int63(), letterIdxMax
		}
		if idx := int(cache & letterIdxMask); idx < len(strSet) {
			b[i] = strSet[idx]
			i--
		}
		cache >>= letterIdxBits
		remain--
	}

	return *(*string)(unsafe.Pointer(&b))
}

// String returns a cryptographically secure random string.
func String(n int) string {
	return randStringBytesMaskImprSrcUnsafe(n, randomCharset)
}

// Bytes returns an arbitrary set of bytes.
func Bytes(n uint32) ([]byte, error) {
	b := make([]byte, n)
	_, err := cryptoRand.Read(b)
	if err != nil {
		return nil, err
	}

	return b, nil
}

// Uint64 returns a cryptographically secure strongly random uint64.
func Uint64() uint64 {
	return defaultSecureSource.Uint64()
}

// Intn returns, as an int, a cryptographically secure non-negative
// random number in [0,n). It panics if n <= 0.
func Intn(n int) int {
	return defaultSecureRand.Intn(n)
}
