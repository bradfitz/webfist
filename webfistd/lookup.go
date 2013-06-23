package main

import (
  "encoding/json"
  "log"
  "net/http"
  "strings"

  "github.com/bradfitz/webfist"
)

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

  emailAddr := webfist.NewEmailAddr(resource)
  foundData := s.lookup.WebFinger(emailAddr)
  if foundData == nil {
    log.Printf("Not found: %s", emailAddr.Canonical())
    http.NotFound(w, r)
    return
  }
  b, err := json.Marshal(foundData.JSON)
  if err != nil {
    log.Printf("Bad data for resource: %s -- %v", emailAddr.Canonical(), err)
    http.Error(w, "Bad data", http.StatusInternalServerError)
    return
  }

  log.Printf("Found user %s -- %v", emailAddr.Canonical(), foundData.JSON)
  w.Write(b)
}

type emailLookup struct {
  storage webfist.Storage
}

func (l *emailLookup) WebFinger(emailAddr *webfist.EmailAddr) *webfist.WebFingerResponse {
  resp := &webfist.WebFingerResponse{JSON: make(map[string]interface{})}
  resp.JSON["hi"] = "meep";
  return resp;
}

func NewLookup(storage webfist.Storage) webfist.Lookup {
  return &emailLookup{
    storage: storage,
  }
}
