package main

import (
	"fmt"
	"net/http"

	"github.com/bradfitz/webfist"
)

func (s *server) WebFormAdd(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "Invalid method.", 400)
		return
	}
	all := r.FormValue("email")
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

	am, err := em.Assignments()
	if err != nil {
		http.Error(w, "Email didn't contain WebFist commands: " + err.Error(), 400)
		return
	}

	err = s.storage.PutEmail(from, em)
	if err != nil {
		http.Error(w, "Storage error: " + err.Error(), 500)
		return
	}
	fmt.Fprintf(w, "Saved. Extracted email = %#v", am)
}
