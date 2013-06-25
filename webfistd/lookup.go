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
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"sort"
	"strings"

	"github.com/bradfitz/webfist"
)

func (s *server) Lookup(w http.ResponseWriter, r *http.Request) {
	if r.ParseForm() != nil {
		http.Error(w, "Bad request", http.StatusBadRequest)
		return
	}
	resource := r.Form.Get("resource")
	if resource == "" {
		http.Error(w, "'resource' missing", http.StatusBadRequest)
		return
	}
	emailLikeId := strings.TrimPrefix(resource, "acct:")
	if emailLikeId == resource {
		http.Error(w, "'resource' must start with 'acct:'", http.StatusBadRequest)
		return
	}

	foundData, err := s.lookup.WebFinger(emailLikeId)
	if err != nil {
		log.Printf("Error looking up resource: %s -- %v", emailLikeId, err)
		http.Error(w, "Error doing lookup", http.StatusInternalServerError)
		return
	}
	if foundData == nil {
		log.Printf("Not found: %s", emailLikeId)
		http.NotFound(w, r)
		return
	}
	b, err := json.Marshal(foundData)
	if err != nil {
		log.Printf("Bad data for resource: %s -- %v", emailLikeId, err)
		http.Error(w, "Bad data from lookup", http.StatusInternalServerError)
		return
	}

	log.Printf("Found user %s -- %+v", emailLikeId, foundData)
	w.Header().Add("Content-Type", "application/json; charset=utf-8")
	w.Header().Add("Access-Control-Allow-Origin", "*")
	w.Write(b)
}

type emailLookup struct {
	storage webfist.Storage
}

type byEmailDate []*webfist.Email


func (s byEmailDate) Len() int {
	return len(s)
}

func (s byEmailDate) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}

func (s byEmailDate) Less(i, j int) bool {
	d1, err := s[i].Date()
	if err != nil {
		return false
	}
	d2, err := s[j].Date()
	if err != nil {
		return false
	}
	return d1.Before(d2)
}

func (l *emailLookup) WebFinger(addr string) (*webfist.WebFingerResponse, error) {
	emailAddr := webfist.NewEmailAddr(addr)
	emailList, err := l.storage.Emails(emailAddr)
	if err != nil {
		return nil, err
	}
	if len(emailList) == 0 {
		return nil, nil
	}
	sort.Sort(byEmailDate(emailList))
	lastEmail := emailList[len(emailList) - 1]
	// TODO: Garbage collect old emails

	delegationURL, err := lastEmail.WebFist()
	if err != nil {
		return nil, err
	}
	encSHA1, err := lastEmail.EncSHA1()
	if err != nil {
		return nil, err
	}
	proofURL := fmt.Sprintf("%s/webfist/proof/%s-%s?decrypt=%s", *baseURL, emailAddr.HexKey(), encSHA1, url.QueryEscape(emailAddr.Canonical()))

	resp := &webfist.WebFingerResponse {
		Subject: emailAddr.Canonical(),
		Links: []webfist.Link {
			{
				Rel: "http://webfist.org/spec/rel",
				Href: delegationURL,
				Properties: map[string]string {
					"http://webfist.org/spec/proof": proofURL,
				},
			},
		},
	}
	return resp, nil
}

func NewLookup(storage webfist.Storage) webfist.Lookup {
	return &emailLookup{
		storage: storage,
	}
}
