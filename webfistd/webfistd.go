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

package main

import (
	"flag"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/bradfitz/go-smtpd/smtpd"
	"github.com/bradfitz/runsit/listen"
	"github.com/bradfitz/webfist"
)

var (
	webAddr     = listen.NewFlag("web", ":8080", "Web port")
	smtpAddr    = listen.NewFlag("smtp", ":2500", "SMTP port")
	storageRoot = flag.String("root", "", "Root for local disk storage")
	baseURL     = flag.String("base", "http://webfist.org", "Base URL without trailing slash for all server-side generated URLs.")
	peers       = flag.String("peers", "", "Comma-separated list of hosts to replicate from.")
)

type server struct {
	httpServer http.Server
	smtpServer *smtpd.Server
	lookup     webfist.Lookup
	storage    webfist.Storage
	peers      []string // hosts
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

	storage, err := NewDiskStorage(*storageRoot)
	if err != nil {
		log.Fatalf("Disk storage of %s: %v", *storageRoot, err)
	}
	srv := &server{
		storage: storage,
		lookup:  NewLookup(storage),
	}
	if *peers != "" {
		srv.peers = strings.Split(*peers, ",")
	}
	go srv.syncFromPeers()
	srv.initSMTPServer()
	log.Printf("Server up. web %s, smtp %s", webAddr, smtpAddr)
	go srv.runSMTP(smtpln)

	http.HandleFunc("/.well-known/webfinger", srv.Lookup)
	http.HandleFunc("/webfist/bump", srv.ServeRecent)
	http.HandleFunc("/webfist/proof/", srv.ServeBlob)
	http.HandleFunc("/add", srv.WebFormAdd)

	log.Fatal(srv.httpServer.Serve(webln))
}
