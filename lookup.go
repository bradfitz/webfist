// Package webfist implements WebFist verification.
package webfist

import (
  // "fmt"
  "log"
  "net/http"
)

type Result struct {
  JSON map[string]interface{}
}

type Lookup interface {
  Get(email string) Result
}

type Server struct {
  Lookup
}

func (l Server) HandleLookup(w http.ResponseWriter, r *http.Request) {
  if r.ParseForm() != nil {
    http.Error(w, "Bad request", http.StatusBadRequest)
    return
  }
  resource := r.Form.Get("resource")
  if resource == "" {
    http.Error(w, "Bad request", http.StatusBadRequest)
    return
  }
  log.Printf("Lookup request for %v", resource)
}

type Dummy struct {}

func (l Dummy) Get(email string) (r Result) {
  r.JSON["hi"] = "meep";
  return r;
}
