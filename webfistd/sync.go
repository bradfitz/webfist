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
	"bufio"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
	"time"
	"flag"
)

var syncInterval = flag.Duration("sync_interval", 10 * time.Second, "Sync poll interval.")

func (s *server) syncFromPeers() {
	for _, peer := range s.peers {
		go s.syncFromPeer(peer)
	}
}

func (srv *server) syncFromPeer(host string) {
	log.Printf("Starting sync from %v", host)
	var ims time.Time
	var sleep time.Duration
	for {
		time.Sleep(sleep)
		url := "http://"+host+"/webfist/bump"
		log.Printf("Doing request to %v", url)
		req, err := http.NewRequest("GET", url, nil)
		if err != nil {
			panic("Bogus host: " + host + ": " + err.Error())
		}
		sleep = *syncInterval
		res, err := http.DefaultClient.Do(req)
		if err != nil {
			log.Printf("Fetch error from %s: %v", host, err)
			continue
		}
		sc := bufio.NewScanner(res.Body)
		for sc.Scan() {
			s := strings.Split(sc.Text(), " ")
			if len(s) != 2 {
				continue
			}
			modTime, err := time.Parse(time.RFC3339, s[0])
			if err != nil {
				continue
			}
			if modTime.After(ims) {
				ims = modTime
			}
			s = strings.Split(s[1], "-")
			if len(s) != 2 {
				continue
			}
			addrHexKey, encSHA1 := s[0], s[1]
			_, err = srv.storage.StatEncryptedEmail(addrHexKey, encSHA1)
			if err != nil {
				log.Printf("Need to fetch %s-%s", addrHexKey, encSHA1)
				res, err := http.Get("http://" + host + "/webfist/proof/" + addrHexKey + "-" + encSHA1)
				if err != nil {
					log.Printf("Error fetching %s-%s: %v", addrHexKey, encSHA1, err)
					continue
				}
				slurp, err := ioutil.ReadAll(io.LimitReader(res.Body, 100<<10))
				if err != nil {
					log.Printf("Error fetching %s-%s: %v", addrHexKey, encSHA1, err)
					continue
				}
				err = srv.storage.PutEncryptedEmail(addrHexKey, encSHA1, slurp)
				if err != nil {
					log.Printf("Error storing fetched %s-%s: %v", addrHexKey, encSHA1, err)
					continue
				}
				log.Printf("Synced %s-%s from %s", addrHexKey, encSHA1, host)
			}
		}
		if err := sc.Err(); err != nil {
			log.Printf("Scan error: %v", err)
		}
		res.Body.Close()
	}
}
