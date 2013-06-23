package main

import (
	"crypto/sha1"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"regexp"

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

func (s *diskStorage) hexPath(hex string) string {
	return filepath.Join(s.root, hex[:3], hex[3:6], hex)
}

func (s *diskStorage) getEmailRoot(addr *webfist.EmailAddr) string {
	return s.hexPath(addr.HexKey())
}

var (
	rxSHA1            = regexp.MustCompile(`^[0-9a-f]{40,40}$`)
	errInvalidBlobref = errors.New("Invalid sha1")
)

func (s *diskStorage) StatEncryptedEmail(sha1 string) (size int, err error) {
	if !rxSHA1.MatchString(sha1) {
		return 0, errInvalidBlobref
	}
	fi, err := os.Stat(s.hexPath(sha1))
	if err != nil {
		return
	}
	return int(fi.Size()), nil
}

func (s *diskStorage) EncryptedEmail(sha1 string) ([]byte, error) {
	if !rxSHA1.MatchString(sha1) {
		return nil, errInvalidBlobref
	}
	return ioutil.ReadFile(s.hexPath(sha1))
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

func (s *diskStorage) Emails(addr *webfist.EmailAddr) ([]*webfist.Email, error) {
	emailRoot := s.getEmailRoot(addr)
	file, err := os.Open(emailRoot)
	if err != nil {
		return nil, err
	}
	infoList, err := file.Readdir(-1)
	if err != nil {
		return nil, err
	}
	result := make([]*webfist.Email, len(infoList))
	for i, info := range infoList {
		emailPath := filepath.Join(emailRoot, info.Name())
		file, err := os.Open(emailPath)
		if err != nil {
			return nil, err
		}
		all, err := ioutil.ReadAll(addr.Decrypter(file))
		if err != nil {
			return nil, err
		}
		email, err := webfist.NewEmail(all)
		if err != nil {
			return nil, err
		}
		result[i] = email
	}
	return result, nil
}
