// Package webfist implements WebFist.
package webfist

// Storage is the interface implemented by backends.
type Storage interface {
	PutEmail(*EmailAddr, *Email) error
	Emails(*EmailAddr) ([]*Email, error)

	// StatEncryptedBlob returns the size of the encrypted blob on
	// disk. addrKey (the Email's HexKey) and encSHA1 (the SHA-1
	// of the encrypted email) are lowercase hex. The err will be
	// os.ErrNotExist if the file is doesn't exist.
	StatEncryptedEmail(addrKey, encSHA1 string) (size int, err error)

	// EncryptedEmail returns the encrypted email with for the
	// addrKey (the Email's HexKey) and encSHA1 (the SHA-1 of |
	// fi, err := os.Stat(s.hexPath(sha1)) the encrypted
	// email). Both are lowercase hex.  The err will be
	// os.ErrNotExist if the file is doesn't exist.
	EncryptedEmail(addrKey, sha1 string) ([]byte, error)
}

// Payload from: http://tools.ietf.org/html/draft-ietf-appsawg-webfinger
// TODO: Make this type pretty and more native, not just a bag of properties.
type WebFingerResponse struct {
  JSON map[string]interface{}
}

// Lookup performs a WebFinger query for an email address and returns all known
// data for that address. Implementations may do standard WebFinger lookups over
// the network, fallback to using the WebFist network, or use local storage to
// map email address to WebFinger response.
type Lookup interface {
  WebFinger(emailAddr string) (*WebFingerResponse, error)
}
