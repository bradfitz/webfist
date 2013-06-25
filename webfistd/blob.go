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
	"bytes"
	"io/ioutil"
	"log"
	"net/http"
	"path"
	"strings"

	"github.com/bradfitz/webfist"
)

func (s *server) ServeBlob(w http.ResponseWriter, r *http.Request) {
	if r.ParseForm() != nil {
		log.Printf("Could not parse form: %+v", r)
		http.Error(w, "Bad request", http.StatusBadRequest)
		return
	}

	base := path.Base(r.URL.Path)
	parts := strings.SplitN(base, "-", 2)
	if len(parts) != 2 {
		log.Printf("Invalid blob path: %q", r.URL.Path)
		http.Error(w, "Bad request", http.StatusBadRequest)
		return
	}

	hexKey, encSHA1 := parts[0], parts[1]

	all, err := s.storage.EncryptedEmail(hexKey, encSHA1)
	if err != nil {
		http.NotFound(w, r)
		return
	}

	decryptKey := r.Form.Get("decrypt")
	if decryptKey == "" {
		w.Header().Add("Content-Type", "application/octet-stream")
		w.Write(all)
		return
	}

	addr := webfist.NewEmailAddr(decryptKey)
	decrypted, err := ioutil.ReadAll(addr.Decrypter(bytes.NewReader(all)))
	if err != nil {
		http.Error(w, "Could not decrypt blob", http.StatusInternalServerError)
		return
	}

	w.Header().Add("Content-Type", "text/plain; charset=utf-8")
	w.Header().Add("Access-Control-Allow-Origin", "*")
	w.Write(decrypted)
}
