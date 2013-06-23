package main

import (
	"flag"
	"log"
	"net/http"

	"github.com/bradfitz/go-smtpd/smtpd"
	"github.com/bradfitz/runsit/listen"
	// "github.com/bradfitz/webfist"
)

var (
	webAddr  = listen.NewFlag("web", ":8080", "Web port")
	smtpAddr = listen.NewFlag("smtp", ":2500", "SMTP port")
)

type server struct {
	httpServer http.Server
	smtpServer *smtpd.Server
}

func (s *server) runSMTP(ln net.Listener) {
	err := s.smtpServer.Serve(ln)
	log.Fatalf("SMTP failure: %v", err)
}

func main() {
	flag.Parse()
	webln, err := webAddr.Listen()
	if err != nil {
		log.Fatalf("web listen: %v", err)
	}
	smtpln, err := smtpAddr.Listen()
	if err != nil {
		log.Fatalf("SMTP listen: %v", err)
	}

	srv := &server{
		smtpServer: &smtpd.Server{
			ReadTimeout:  5 * time.Minute,
			WriteTimeout: 5 * time.Minute,
		},
	}

	// TODO: Actually hook up the lookup
	// lookup := &lookupHandler {
	// 	lookup: NewLookup(&DummyStorage{}),
	// }
	// http.Handle("/.well-known/webfinger", lookup)

	log.Printf("Server up. web %s, smtp %s", webAddr, smtpAddr)
	go srv.runSMTP(smtpln)

	log.Fatal(srv.httpServer.Serve(webln))
}
