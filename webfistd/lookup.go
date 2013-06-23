package main

import (
  "encoding/json"
  "log"
  "net/http"
  "strings"

  "github.com/bradfitz/webfist"
)

func init() {
  Foo()
}

func Foo() {
  log.Printf("Hi")
}

func (s *server) HandleLookup(w http.ResponseWriter, r *http.Request) {
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
  if (emailLikeId == resource) {
    http.Error(w, "'resource' must start with 'acct:'", http.StatusBadRequest)
    return
  }

  emailAddr := NewEmailAddr(resource)
  foundData := s.Get(emailAddr)
  if foundData == nil {
    log.Printf("Not found: %s", emailAddr.email)
    http.NotFound(w, r)
    return
  }
  b, err := json.Marshal(foundData.JSON)
  if err != nil {
    log.Printf("Bad data for resource: %s -- %v", emailAddr.email, err)
    http.Error(w, "Bad data", http.StatusInternalServerError)
    return
  }

  log.Printf("Found user %s -- %v", emailAddr.email, foundData.JSON)
  w.Write(b)
}
