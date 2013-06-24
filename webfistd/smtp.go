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

// SMTP-related parts of webfistd.

package main

import (
	"bytes"
	"errors"
	"flag"
	"log"
	"net"
	"strings"
	"time"

	"github.com/bradfitz/go-smtpd/smtpd"
	"github.com/bradfitz/webfist"
)

var hostName = flag.String("hostname", "webfist.org", "Hostname to announce over SMTP")

var (
	lf            = []byte("\n")
	crlf          = []byte("\r\n")
	dkimSignature = []byte("DKIM-Signature")
)

func (s *server) initSMTPServer() {
	s.smtpServer = &smtpd.Server{
		ReadTimeout:  5 * time.Minute,
		WriteTimeout: 5 * time.Minute,
		Hostname:     *hostName,
		OnNewMail:    s.onNewMail,
	}
}

func (s *server) onNewMail(conn smtpd.Connection, from smtpd.MailAddress) (smtpd.Envelope, error) {
	log.Printf("smtp: new mail from %s", from.Email())
	return &env{s: s, from: webfist.NewEmailAddr(from.Email())}, nil
}

func (s *server) runSMTP(ln net.Listener) {
	err := s.smtpServer.Serve(ln)
	log.Fatalf("SMTP failure: %v", err)
}

type env struct {
	from       *webfist.EmailAddr
	buf         bytes.Buffer
	s           *server
	headersDone bool
}

// hasSignatureHeaders true if e likely contains a signed email.
// False positives are okay.
func (e *env) hasSignatureHeader() bool {
	if bytes.Contains(e.buf.Bytes(), dkimSignature) {
		return true
	}
	if strings.Contains(strings.ToLower(e.buf.String()), "dkim-signature") {
		return true
	}
	return false
}

func (e *env) AddRecipient(rcpt smtpd.MailAddress) error { return nil }
func (e *env) BeginData() error                          { return nil }

func (e *env) Write(line []byte) error {
	if e.buf.Len()+len(line) > webfist.MaxEmailSize {
		return errors.New("email too large for webfist")
	}
	e.buf.Write(line)
	if !e.headersDone && (bytes.Equal(line, lf) || bytes.Equal(line, crlf)) {
		e.headersDone = true
		if !e.hasSignatureHeader() {
			log.Printf("Rejecting email that isn't signed.")
			return errors.New("not a signed email")
		}
	}
	return nil
}

func (e *env) Close() error {
	em, err := webfist.NewEmail(e.buf.Bytes())
	if err != nil {
		return err
	}
	verified := em.Verify()
	log.Printf("email from %v; verified = %v", e.from.Canonical(), verified)
	if !verified {
		return errors.New("DKIM verification failed")
	}

	_, err = em.WebFist()
	if err != nil {
		return errors.New("Invalid or missing WebFist commands in email.")
	}
	from, err := em.From()
	if err != nil {
		return errors.New("Bogus From header")
	}

	return e.s.storage.PutEmail(from, em)
}
