package main

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"errors"
	"fmt"

	"github.com/urfave/cli/v2"
)

//===========================================================================
// Cipher Actions
//===========================================================================

const nonceSize = 12

func cipherDecrypt(c *cli.Context) (err error) {
	if c.NArg() != 2 {
		return cli.Exit("must specify ciphertext and hmac arguments", 1)
	}

	var secret string
	if secret = c.String("key"); secret == "" {
		return cli.Exit("cipher key required", 1)
	}

	var ciphertext, signature []byte
	if ciphertext, err = base64.RawStdEncoding.DecodeString(c.Args().Get(0)); err != nil {
		return cli.Exit(fmt.Errorf("could not decode ciphertext: %s", err), 1)
	}
	if signature, err = base64.RawStdEncoding.DecodeString(c.Args().Get(1)); err != nil {
		return cli.Exit(fmt.Errorf("could not decode signature: %s", err), 1)
	}

	if len(ciphertext) == 0 {
		return cli.Exit("empty cipher text", 1)
	}

	// Create a 32 byte signature of the key
	hash := sha256.New()
	hash.Write([]byte(secret))
	key := hash.Sum(nil)

	// Separate the data from the nonce
	data := ciphertext[:len(ciphertext)-nonceSize]
	nonce := ciphertext[len(ciphertext)-nonceSize:]

	// Validate HMAC signature
	if err = validateHMAC(key, data, signature); err != nil {
		return cli.Exit(err, 1)
	}

	block, err := aes.NewCipher(key)
	if err != nil {
		return cli.Exit(err, 1)
	}

	aesgcm, err := cipher.NewGCM(block)
	if err != nil {
		return cli.Exit(err, 1)
	}

	plainbytes, err := aesgcm.Open(nil, nonce, data, nil)
	if err != nil {
		return cli.Exit(err, 1)
	}

	fmt.Println(string(plainbytes))
	return nil
}

//===========================================================================
// Cipher Helper Functions
//===========================================================================

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
