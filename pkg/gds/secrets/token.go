package secrets

import (
	"math/rand"
	"strings"
	"time"
)

var (
	chars        = []rune("ABCDEFGHIJKLMNPQRSTUVWXYZabcdefghjkmnpqrstuvwxyz1234567890")
	specialChars = []rune("#%&()*+-<=>?@[]^_{}")
)

// CreateToken creates a variable length random token that can be used for passwords or API keys.
func CreateToken(length int) string {
	random := rand.New(rand.NewSource(time.Now().UnixNano()))
	var b strings.Builder
	for i := 0; i < length; i++ {
		if random.Float64() <= 0.90 {
			b.WriteRune(chars[random.Intn(len(chars))])
		} else {
			b.WriteRune(specialChars[random.Intn(len(specialChars))])
		}
	}
	return b.String()
}
