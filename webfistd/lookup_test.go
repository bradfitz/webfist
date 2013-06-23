package main

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/bradfitz/webfist"
)

type DummyStorage struct{}

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
	testServer *server
)

func init() {
	storage := &DummyStorage{}
	testServer = &server{
		storage: storage,
		lookup:  NewLookup(storage),
	}
}

func TestEmailLookup(t *testing.T) {
	req, _ := http.NewRequest("GET", "http://example.com/foo?resource=acct:myname@example.com", nil)
	resp := httptest.NewRecorder()
	testServer.Lookup(resp, req)
	body := resp.Body.String()
	wants := `{"subject":"foo@bar.com"}`
	if body != wants {
		t.Fatalf("Body = %q; want %q", body, wants)
	}
}
