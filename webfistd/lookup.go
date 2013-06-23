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
  foundData := s.storage.WebFinger(emailAddr)
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

// TODO: Move this somewhere else
type Dummy struct {}

func (l *Dummy) PutEmail(*webfist.EmailAddr, *webfist.Email) error {
  return nil
}

func (l *Dummy) Emails(*webfist.EmailAddr) ([]*webfist.Email, error) {
  return nil, nil
}

func (l *Dummy) WebFinger(emailAddr *webfist.EmailAddr) (r *webfist.WebFingerResponse) {
  r = &webfist.WebFingerResponse{JSON: make(map[string]interface{})}
  r.JSON["hi"] = "meep";
  return r;
}
