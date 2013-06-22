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

// MaxEmailSize is the maxium size of an RFC 822 email, including
// both its headers and body.
const MaxEmailSize = 16 << 10

var (
	dummyIV  = make([]byte, 16) // all zeros
	fistSalt = []byte("WebFist salt.")
)

// EmailAddr provides utility functions on a wrapped email address.
type EmailAddr struct {
	email string // canonical

	keyOnce sync.Once
	lazyKey []byte // scrypt
}

// NewEmailAddr returns a EmailAddr wrapper around an email address string.
// The incoming email address does not need to be canonicalized.
func NewEmailAddr(addr string) *EmailAddr {
	return &EmailAddr{
		email: canonicalEmail(addr),
	}
}

// Canonical returns the canonical version of the email address.
func (e *EmailAddr) Canonical() string {
	return e.email
}

// HexKey returns the human-readable, lowercase hex version of
// the email address's key.
func (e *EmailAddr) HexKey() string {
	return fmt.Sprintf("%x", e.getKey())
}

func (e *EmailAddr) getKey() []byte {
	e.keyOnce.Do(e.initLazyKey)
	return e.lazyKey
}

func (e *EmailAddr) initLazyKey() {
	key, err := scrypt.Key([]byte(e.Canonical()), fistSalt, 16384*8, 8, 1, 32)
	if err != nil {
		panic(err)
	}
	e.lazyKey = key
}

func (e *EmailAddr) block() cipher.Block {
	block, err := aes.NewCipher(e.encryptionKey())
	if err != nil {
		panic(err)
	}
	return block
}

// encryptionKey returns the AES-128 key for this email address.
func (e *EmailAddr) encryptionKey() []byte {
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

func (e *EmailAddr) Encrypter(w io.Writer) io.Writer {
	return cipher.StreamWriter{
		S: cipher.NewCTR(e.block(), dummyIV),
		W: w,
	}
}

func (e *EmailAddr) Decrypter(r io.Reader) io.Reader {
	return cipher.StreamReader{
		S: cipher.NewCTR(e.block(), dummyIV),
		R: r,
	}
}
