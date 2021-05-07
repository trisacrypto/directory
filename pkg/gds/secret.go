package trisads

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/hmac"
	crand "crypto/rand"
	"crypto/sha256"
	"errors"
	"io"
	"math/rand"
	"strings"
	"time"
)

const nonceSize = 12

var chars = []rune("ABCDEFGHIJKLMNPQRSTUVWXYZabcdefghjkmnpqrstuvwxyz1234567890")

// CreateToken creates a variable length random token that can be used for passwords or API keys.
func CreateToken(length int) string {
	random := rand.New(rand.NewSource(time.Now().UnixNano()))
	var b strings.Builder
	for i := 0; i < length; i++ {
		b.WriteRune(chars[random.Intn(len(chars))])
	}
	return b.String()
}

// Encrypt is a helper utility to encrypt a plain text string with the server's secret
// token, returns a cipher string which is the base64 encoded
func (s *Server) Encrypt(plaintext string) (ciphertext, signature []byte, err error) {
	// Create a 32 byte signature of the key
	hash := sha256.New()
	hash.Write([]byte(s.conf.SecretKey))
	key := hash.Sum(nil)

	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, nil, err
	}

	aesgcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, nil, err
	}

	nonce, err := genNonce(nonceSize)
	if err != nil {
		return nil, nil, err
	}

	ciphertext = aesgcm.Seal(nil, nonce, []byte(plaintext), nil)
	if len(ciphertext) == 0 {
		return nil, nil, errors.New("could not encrypt secret with aes gcm")
	}

	sig, err := createHMAC(key, ciphertext)
	if err != nil {
		return nil, nil, err
	}

	// Concatenate the ciphertext and the nonce to facilitate decryption
	ciphertext = append(ciphertext, nonce...)
	return ciphertext, sig, nil
}

// Decrypt the ciphertext with the server's secret key and verify the HMAC.
func (s *Server) Decrypt(ciphertext, signature []byte) (plaintext string, err error) {
	if len(ciphertext) == 0 {
		return "", errors.New("empty cipher text")
	}

	// Create a 32 byte signature of the key
	hash := sha256.New()
	hash.Write([]byte(s.conf.SecretKey))
	key := hash.Sum(nil)

	// Separate the data from the nonce
	data := ciphertext[:len(ciphertext)-nonceSize]
	nonce := ciphertext[len(ciphertext)-nonceSize:]

	// Validate HMAC signature
	if err = validateHMAC(key, data, signature); err != nil {
		return "", err
	}

	block, err := aes.NewCipher(key)
	if err != nil {
		return "", err
	}

	aesgcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", err
	}

	plainbytes, err := aesgcm.Open(nil, nonce, data, nil)
	if err != nil {
		return "", err
	}

	return string(plainbytes), nil
}

func genNonce(n int) ([]byte, error) {
	// Never use more than 2^32 random nonces with a given key because of the risk of a repeat.
	nonce := make([]byte, n)
	if _, err := io.ReadFull(crand.Reader, nonce); err != nil {
		return nil, err
	}
	return nonce, nil
}

func createHMAC(key, data []byte) ([]byte, error) {
	if len(data) == 0 {
		return nil, errors.New("cannot sign empty data")
	}
	hm := hmac.New(sha256.New, key)
	hm.Write(data)
	return hm.Sum(nil), nil
}

func validateHMAC(key, data, sig []byte) error {
	hmac, err := createHMAC(key, data)
	if err != nil {
		return err
	}

	if !bytes.Equal(sig, hmac) {
		return errors.New("HMAC mismatch")
	}
	return nil
}
