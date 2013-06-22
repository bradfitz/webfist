// Package webfist implements WebFist verification.
package webfist

import (
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

// EmailKey returns the key used to store WebFist claims.
func EmailKey(email string) string {
	email = CanonicalEmail(email)
	key, err := scrypt.Key([]byte(email), fistSalt, 16384 * 8, 8, 1, 32)
	if err != nil {
		panic(err)
	}
	return fmt.Sprintf("%x", key)
}
