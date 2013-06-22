package webfist

import "errors"

// MaxEmailSize is the maxium size of an RFC 822 email, including
// both its headers and body.
const MaxEmailSize = 16 << 10

// Email wraps a signed email.
type Email struct {
	all []byte
}

func NewEmail(all []byte) (*Email, error) {
	if len(all) > MaxEmailSize {
		return nil, errors.New("email too large")
	}
	e := &Email{
		all: all,
	}
	return e, nil
}
