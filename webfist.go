// Package webfist implements WebFist verification.
package webfist

import (
	"crypto/sha1"
	"io"

	"crypto/aes"
	"crypto/cipher"
	"fmt"

	"code.google.com/p/go.crypto/scrypt"
)

// CanonicalEmail returns the canonicalized version of the provided
// email address.
func CanonicalEmail(email string) string {
	// TODO
	return email
}

var fistSalt = []byte("WebFist salt.")

// EmailKey returns the human-readable hex version of EmailKey.
func EmailKeyString(email string) string {
	return fmt.Sprintf("%x", EmailKey(email))
}

// EmailKey returns the one-way slow hash of email.
// This result key is the key that the encrypted blobs are stored by.
func EmailKey(email string) []byte {
	// TODO: optional cache.
	email = CanonicalEmail(email)
	key, err := scrypt.Key([]byte(email), fistSalt, 16384*8, 8, 1, 32)
	if err != nil {
		panic(err)
	}
	return key
}

// EncryptionKey returns the AES-128 key used to encrypt
// and decrypt the payload blobs.
func EncryptionKey(email string, emailKey []byte) []byte {
	s1 := sha1.New()
	io.WriteString(s1, CanonicalEmail(email))
	s1.Write(emailKey)
	return s1.Sum(nil)[:16]
}

var dummyIV = make([]byte, 16) // all zeros

func emailAESBlock(email string) cipher.Block {
	block, err := aes.NewCipher(EncryptionKey(email, EmailKey(email)))
	if err != nil {
		panic(err)
	}
	return block
}

func NewEncrypter(email string, w io.Writer) io.Writer {
	return cipher.StreamWriter{
		S: cipher.NewCTR(emailAESBlock(email), dummyIV),
		W: w,
	}
}

func NewDecrypter(email string, r io.Reader) io.Reader {
	return cipher.StreamReader{
		S: cipher.NewCTR(emailAESBlock(email), dummyIV),
		R: r,
	}
}
