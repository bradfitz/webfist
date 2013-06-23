package main

import (
	"flag"
	"log"
	"net/http"

	"github.com/bradfitz/go-smtpd/smtpd"
	"github.com/bradfitz/runsit/listen"
)

var (
	webAddr  = listen.NewFlag("web", ":8080", "Web port")
	smtpAddr = listen.NewFlag("smtp", ":2500", "SMTP port")
)

type server struct {
	httpServer http.Server
	smtpServer *smtpd.Server
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

	var srv server
	srv.initSMTPServer()
	log.Printf("Server up. web %s, smtp %s", webAddr, smtpAddr)
	go srv.runSMTP(smtpln)

	// TODO: Actually hook up the lookup
	// lookup := &lookupHandler {
	// 	lookup: NewLookup(&DummyStorage{}),
	// }
	// http.Handle("/.well-known/webfinger", lookup)

	log.Fatal(srv.httpServer.Serve(webln))
}
