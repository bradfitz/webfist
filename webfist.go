// Package webfist implements WebFist.
package webfist

type WebFingerResponse struct {
  JSON map[string]interface{}
}

// Storage is the interface implemented by backends.
type Storage interface {
	PutEmail(*EmailAddr, *Email) error
	Emails(*EmailAddr) ([]*Email, error)
  WebFinger(emailAddr *EmailAddr) *WebFingerResponse
}
