// SMTP-related parts of webfistd.

package main

import (
	"log"
	"net"
	"time"

	"github.com/bradfitz/go-smtpd/smtpd"
)

func (s *server) initSMTPServer() {
	s.smtpServer = &smtpd.Server{
		ReadTimeout:  5 * time.Minute,
		WriteTimeout: 5 * time.Minute,
	}
}

func (s *server) runSMTP(ln net.Listener) {
	err := s.smtpServer.Serve(ln)
	log.Fatalf("SMTP failure: %v", err)
}
