// Package webfist implements WebFist.
package webfist

// Storage is the interface implemented by backends.
type Storage interface {
	PutEmail(*EmailAddr, *Email) error
	Emails(*EmailAddr) ([]*Email, error)
}

// Defined in: http://tools.ietf.org/html/draft-ietf-appsawg-webfinger
type Link struct {
	Rel string `json:"rel"`
	Type string `json:"type"`
	Href string `json:"href"`
	Titles []string `json:"titles,omitempty"`
	Properties map[string]string `json:"properties,omitempty"`
}

type WebFingerResponse struct {
	Subject string `json:"subject"`
	Aliases []string `json:"aliases,omitempty"`
	Properties map[string]string `json:"properties,omitempty"`
	Links []Link `json:"links,omitempty"`
}

// Lookup performs a WebFinger query for an email address and returns all known
// data for that address. Implementations may do standard WebFinger lookups over
// the network, fallback to using the WebFist network, or use local storage to
// map email address to WebFinger response.
type Lookup interface {
  WebFinger(emailAddr string) (*WebFingerResponse, error)
}
