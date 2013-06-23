// Package webfist implements WebFist.
package webfist

// Payload from: tools.ietf.org/html/draft-ietf-appsawg-webfinger-14‎
// TODO: Consider restricting to only delegation to another WebFinger server.
type WebFingerResponse struct {
  JSON map[string]interface{}
}

// Storage is the interface implemented by backends.
type Storage interface {
	PutEmail(*EmailAddr, *Email) error
	Emails(*EmailAddr) ([]*Email, error)
  WebFinger(emailAddr *EmailAddr) *WebFingerResponse
}
