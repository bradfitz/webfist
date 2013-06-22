// Package webfist implements WebFist.
package webfist

import (
	"crypto/sha1"
	"io"
	"sync"

	"crypto/aes"
	"crypto/cipher"
	"fmt"

	"code.google.com/p/go.crypto/scrypt"
)

var (
	dummyIV  = make([]byte, 16) // all zeros
	fistSalt = []byte("WebFist salt.")
)

// Email provides utility functions on a wrapped email address.
type Email struct {
	email string // canonical

	keyOnce sync.Once
	lazyKey []byte // scrypt
}

// NewEmail returns a Email wrapper around an email address string.
// The incoming email address does not need to be canonicalized.
func NewEmail(email string) *Email {
	return &Email{
		email: canonicalEmail(email),
	}
}

// Canonical returns the canonical version of the email address.
func (e *Email) Canonical() string {
	return e.email
}

// HexKey returns the human-readable, lowercase hex version of
// the email address's key.
func (e *Email) HexKey() string {
	return fmt.Sprintf("%x", e.getKey())
}

func (e *Email) getKey() []byte {
	e.keyOnce.Do(e.initLazyKey)
	return e.lazyKey
}

func (e *Email) initLazyKey() {
	key, err := scrypt.Key([]byte(e.Canonical()), fistSalt, 16384*8, 8, 1, 32)
	if err != nil {
		panic(err)
	}
	e.lazyKey = key
}

func (e *Email) block() cipher.Block {
	block, err := aes.NewCipher(e.encryptionKey())
	if err != nil {
		panic(err)
	}
	return block
}

// encryptionKey returns the AES-128 key for this email address.
func (e *Email) encryptionKey() []byte {
	s1 := sha1.New()
	io.WriteString(s1, e.email)
	s1.Write(e.getKey())
	return s1.Sum(nil)[:16]
}

// canonicalEmail returns the canonicalized version of the provided
// email address.
func canonicalEmail(email string) string {
	// TODO
	return email
}

func (e *Email) Encrypter(w io.Writer) io.Writer {
	return cipher.StreamWriter{
		S: cipher.NewCTR(e.block(), dummyIV),
		W: w,
	}
}

func (e *Email) Decrypter(r io.Reader) io.Reader {
	return cipher.StreamReader{
		S: cipher.NewCTR(e.block(), dummyIV),
		R: r,
	}
}
