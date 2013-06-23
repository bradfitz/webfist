// Package webfist implements WebFist.
package webfist

// Storage is the interface implemented by backends.
type Storage interface {
	PutEmail(*EmailAddr, *Email) error
	Emails(*EmailAddr) ([]*Email, error)
}

// Payload from: tools.ietf.org/html/draft-ietf-appsawg-webfinger-14‎
// TODO: Consider restricting to only delegation to another WebFinger server.
type WebFingerResponse struct {
  JSON map[string]interface{}
}

type Lookup interface {
  WebFinger(emailAddr *EmailAddr) *WebFingerResponse
}
