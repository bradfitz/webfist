package main

import (
	"crypto/sha1"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/bradfitz/webfist"
)

type diskStorage struct {
	root string
}

func NewDiskStorage(root string) webfist.Storage {
	return &diskStorage{
		root: root,
	}
}

func (s *diskStorage) getEmailRoot(addr *webfist.EmailAddr) string {
	x := addr.HexKey()
	return filepath.Join(s.root, x[:3], x[3:6], x)
}

func (s *diskStorage) PutEmail(addr *webfist.EmailAddr, email *webfist.Email) error {
	emailRoot := s.getEmailRoot(addr)
	err := os.MkdirAll(emailRoot, 0755)
	if err != nil {
		return err
	}

	r, err := email.Encrypted()
	if err != nil {
		return err
	}
	enc, err := ioutil.ReadAll(r)
	if err != nil {
		return err
	}

	s1 := sha1.New()
	s1.Write(enc)
	emailPath := filepath.Join(emailRoot, fmt.Sprintf("%x", s1.Sum(nil)))
	return ioutil.WriteFile(emailPath, enc, 0644)
}

func (s *diskStorage) Emails(*webfist.EmailAddr) ([]*webfist.Email, error) {
	panic("TODO")
}
