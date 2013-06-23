package webfist

import (
	"io/ioutil"
	"path/filepath"
	"testing"
)

func TestDKIMVerifyAvailable(t *testing.T) {
	path, err := findDKIMVerify()
	if err != nil {
		t.Errorf(dkimFailMessage)
	}
	t.Logf("dkimverify at %s", path)
}

func TestEmailVerify(t *testing.T) {
	files := []string{
		"gmail_dkim.txt",
		"facebook_dkim.txt",
		"twitter_dkim.txt",
	}
	for _, file := range files {
		full := filepath.Join("testdata", file)
		all, err := ioutil.ReadFile(full)
		if err != nil {
			t.Fatalf("Error opening %v: %v", full, err)
		}
		e, err := NewEmail(all)
		if err != nil {
			t.Errorf("NewEmail(%s) = %v", file, err)
			continue
		}
		if !e.Verify() {
			t.Errorf("%s didn't verify", file)
		}
		ea, err := e.From()
		if err != nil {
			t.Errorf("%s From error = %v", file, err)
		} else {
			t.Logf("%s From = %v", file, ea.Canonical())
		}
	}
}
