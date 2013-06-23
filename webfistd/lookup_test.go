package main

import (
  "net/http"
  "net/http/httptest"
  "testing"

  "github.com/bradfitz/webfist"
)

type DummyStorage struct {}

func (DummyStorage) PutEmail(*webfist.EmailAddr, *webfist.Email) error {
  return nil
}

func (l DummyStorage) Emails(*webfist.EmailAddr) (result []*webfist.Email, err error) {
  result = make([]*webfist.Email, 1)
  var myVar []byte = []byte("foo")
  result[0], err = webfist.NewEmail(myVar)
  return
}

var (
  testLookup webfist.Lookup
  testHandler http.Handler
)

func init() {
  testLookup = NewLookup(DummyStorage{})
  testHandler = &lookupHandler{
    lookup: testLookup,
  }
}

func TestEmailLookup(t *testing.T) {
  req, _ := http.NewRequest("GET", "http://example.com/foo?resource=acct:myname@example.com", nil)
  resp := httptest.NewRecorder()
  testHandler.ServeHTTP(resp, req)
  body := resp.Body.String()
  wants := `{"hi":"meep"}`
  if body != wants {
    t.Fatalf("Body = %q; want %q", body, wants)
  }
}
