package main

import (
	"net/http"

	"github.com/bradfitz/webfist"
)


func add(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "Invalid method.", 400)
		return
	}
	all := r.FormValue("email")
   	em, err := webfist.NewEmail([]byte(all))
	if err != nil {
		http.Error(w, "Bogus email: " + err.Error(), 500)
		return
	}
	_ = em
}
