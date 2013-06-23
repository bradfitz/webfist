package main

import (
	"flag"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"runtime"

	"github.com/bradfitz/go-smtpd/smtpd"
	"github.com/bradfitz/runsit/listen"
	"github.com/bradfitz/webfist"
)

var (
	webAddr     = listen.NewFlag("web", ":8080", "Web port")
	smtpAddr    = listen.NewFlag("smtp", ":2500", "SMTP port")
	storageRoot = flag.String("root", "", "Root for local disk storage")
)

type server struct {
	httpServer http.Server
	smtpServer *smtpd.Server
	lookup     webfist.Lookup
	storage    webfist.Storage
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

	if *storageRoot == "" {
		varDir := "var"
		if runtime.GOOS == "darwin" {
			varDir = "Library"
		}
		*storageRoot = filepath.Join(os.Getenv("HOME"), varDir, "webfistd")
		if err := os.MkdirAll(*storageRoot, 0700); err != nil {
			log.Fatal(err)
		}
	}

	storage := NewDiskStorage(*storageRoot)
	srv := &server{
		storage: storage,
		lookup:  NewLookup(storage),
	}
	srv.initSMTPServer()
	log.Printf("Server up. web %s, smtp %s", webAddr, smtpAddr)
	go srv.runSMTP(smtpln)

	http.HandleFunc("/.well-known/webfinger", srv.Lookup)
	http.HandleFunc("/add", srv.WebFormAdd)

	log.Fatal(srv.httpServer.Serve(webln))
}
