package main

import (
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"path/filepath"
	"testing"

	"github.com/bradfitz/webfist"
)

type DummyStorage struct{}

func (DummyStorage) PutEmail(*webfist.EmailAddr, *webfist.Email) error {
	return nil
}

func (l DummyStorage) Emails(*webfist.EmailAddr) ([]*webfist.Email, error) {
	files := []string{
		"gmail_dkim.txt",
	}
	var res []*webfist.Email
	for _, file := range files {
		full := filepath.Join("../testdata", file)
		all, err := ioutil.ReadFile(full)
		if err != nil {
			return nil, err
		}
		e, err := webfist.NewEmail(all)
		if err != nil {
			return nil, err
		}
		res = append(res, e)
	}
	return res, nil
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
	wants := `{"subject":"myname@example.com","links":[{"rel":"webfist","href":"http://www.example.com/foo/bar/baz.json"}]}`
	if body != wants {
		t.Fatalf("Body = %q; want %q", body, wants)
	}
}
