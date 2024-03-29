package secrets

import (
	"crypto/rand"
	"encoding/binary"
	"fmt"
	"strings"
)

const (
	alphabet = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	alphanum = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	idxbits  = 6
	idxmask  = 1<<idxbits - 1
	idxmax   = 63 / idxbits
)

// CreateToken creates a variable length random token that can be used for passwords or API keys.
func CreateToken(length int) string {
	return generate(length, alphanum)
}

// Alpha generates a random string of n characters that only includes upper and
// lowercase letters (no symbols or digits).
func Alpha(n int) string {
	return generate(n, alphabet)
}

// AlphaNumeric generates a random string of n characters that includes upper and
// lowercase letters and the digits 0-9.
func AlphaNumeric(n int) string {
	return generate(n, alphanum)
}

// generate is a helper function to create a random string of n characters from the
// character set defined by chars. It uses as efficient a method of generation as
// possible, using a string builder to prevent multiple allocations and a 6 bit mask
// to select 10 random letters at a time to add to the string. This method would be far
// faster if it used math/rand src and the Int63() function, but for API key generation
// it is important to use a cryptographically random generator.
//
// See: https://stackoverflow.com/questions/22892120/how-to-generate-a-random-string-of-a-fixed-length-in-go
func generate(n int, chars string) string {
	if n <= 0 {
		return ""
	}

	sb := strings.Builder{}
	sb.Grow(n)

	for i, cache, remain := n-1, CryptoRandInt(), idxmax; i >= 0; {
		if remain == 0 {
			cache, remain = CryptoRandInt(), idxmax
		}

		if idx := int(cache & idxmask); idx < len(chars) {
			sb.WriteByte(chars[idx])
			i--
		}

		cache >>= idxbits
		remain--
	}

	return sb.String()
}

func CryptoRandInt() uint64 {
	buf := make([]byte, 8)
	if _, err := rand.Read(buf); err != nil {
		panic(fmt.Errorf("cannot generate random number: %w", err))
	}
	return binary.BigEndian.Uint64(buf)
}
