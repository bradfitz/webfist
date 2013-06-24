/*
Copyright 2013 WebFist AUTHORS

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

     http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

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

func (DummyStorage) Emails(*webfist.EmailAddr) ([]*webfist.Email, error) {
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

func (DummyStorage) StatEncryptedEmail(addrKey, encSHA1 string) (size int, err error) {
	panic("Not implemented")
}

func (DummyStorage) EncryptedEmail(addrKey, sha1 string) ([]byte, error) {
	panic("Not implemented")
}

func (DummyStorage) PutEncryptedEmail(addrKey, encSHA1 string, data []byte) error {
	panic("Not implemented")
}

func (DummyStorage) RecentMeta() ([]*webfist.RecentMeta, error) {
	panic("Not implemented")
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
	wants := `{"subject":"myname@example.com","links":[{"rel":"http://webfist.org/spec/rel","href":"http://www.example.com/foo/bar/baz.json","properties":{"http://webfist.org/spec/proof":"http://webfist.org/webfist/proof/9239956c3d0668d7d0009ef14228bfbbc43dfd10-3a3202736e2f25cae0f5acfb011b6436eb28e27d?decrypt=myname%40example.com"}}]}`
	if body != wants {
		t.Fatalf("Body = %q; want %q", body, wants)
	}
}
