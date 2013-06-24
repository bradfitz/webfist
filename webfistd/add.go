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
	"fmt"
	"net/http"
	"strings"

	"github.com/bradfitz/webfist"
)

func (s *server) WebFormAdd(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "Invalid method.", 400)
		return
	}
	all := r.FormValue("email")
	all = strings.TrimLeft(all, " \t\n\r")
   	em, err := webfist.NewEmail([]byte(all))
	if err != nil {
		http.Error(w, "Bogus email: " + err.Error(), 400)
		return
	}

	from, err := em.From()
	if err != nil {
		http.Error(w, "No From", 400)
		return
	}

	if !em.Verify() {
		http.Error(w, "Email didn't verify. No DKIM.", 400)
		return
	}

	webfist, err := em.WebFist()
	if err != nil {
		http.Error(w, "Email didn't contain WebFist command: " + err.Error(), 400)
		return
	}

	err = s.storage.PutEmail(from, em)
	if err != nil {
		http.Error(w, "Storage error: " + err.Error(), 500)
		return
	}
	fmt.Fprintf(w, "Saved. Extracted email = %#v", webfist)
}
