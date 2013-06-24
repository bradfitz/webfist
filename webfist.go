/*
Copyright 2013 WebFist AUTHORS

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

     http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

// Package webfist implements WebFist.
package webfist

import (
	"time"
)

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

	PutEncryptedEmail(addrKey, sha1 string, data []byte) error

	// RecentMeta returns the recently-received encrypted emails.
	RecentMeta() ([]*RecentMeta, error)
}

// RecentMeta describes an encrypted email in the storage system.
type RecentMeta struct {
	AddrHexKey string
	EncSHA1    string
	AddTime    time.Time
}

// Defined in: http://tools.ietf.org/html/draft-ietf-appsawg-webfinger
type Link struct {
	Rel string `json:"rel"`
	Type string `json:"type,omitempty"`
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
	WebFinger(string) (*WebFingerResponse, error)
}
