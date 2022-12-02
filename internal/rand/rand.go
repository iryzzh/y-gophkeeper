//nolint:gochecknoglobals
package rand

// Reader is the standard crypto/rand.Reader with added buffering.
var Reader = defaultSecureSource

// randomCharset contains the characters that can make up a rand.String().
const randomCharset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ1234567890"

var (
	// defaultSecureSource is a concurrency-safe, cryptographically secure
	// math/rand.Reader.
	defaultSecureSource = newSecureSource()
)

// randString generates a random string of a given length.
func randString(n int, strSet string) string {
	b := make([]byte, n)
	for i := range b {
		b[i] = strSet[Reader.Int63()%int64(len(strSet))]
	}
	return string(b)
}

// String returns a cryptographically secure random string.
func String(n int) string {
	return randString(n, randomCharset)
}

// Uint64 returns a cryptographically secure strongly random uint64.
func Uint64() uint64 {
	return defaultSecureSource.Uint64()
}
