// Package webfist implements WebFist verification.
package webfist

import (
  "encoding/json"
  "log"
  "net/http"
)

type Result struct {
  JSON map[string]interface{}
}

type Lookup interface {
  Get(email string) *Result
}

type Server struct {
  Lookup
}

func (s Server) HandleLookup(w http.ResponseWriter, r *http.Request) {
  if r.ParseForm() != nil {
    http.Error(w, "Bad request", http.StatusBadRequest)
    return
  }
  resource := r.Form.Get("resource")
  if resource == "" {
    http.Error(w, "Bad request", http.StatusBadRequest)
    return
  }

  foundData := s.Get(resource)
  if foundData == nil {
    log.Printf("Not found: %v", resource)
    http.NotFound(w, r)
    return
  }
  log.Printf("Found: %v", resource)
  b, err := json.Marshal(foundData.JSON)
  if err != nil {
    log.Printf("Bad data for resource: %v -- %v", resource, err)
    http.Error(w, "Bad data", http.StatusInternalServerError)
    return
  }

  w.Write(b)
}

type Dummy struct {}

func (l Dummy) Get(email string) (r *Result) {
  r = &Result{JSON: make(map[string]interface{})}
  r.JSON["hi"] = "meep";
  return r;
}
