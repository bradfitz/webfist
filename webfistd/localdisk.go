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

func (s *diskStorage) emailRootFromHex(addrHexKey string) string {
	x := addrHexKey
	if len(x) < 7 {
		panic("bogus emailRootFromHex")
	}
	return filepath.Join(s.root, x[:3], x[3:6], x)
}

func (s *diskStorage) emailRoot(addr *webfist.EmailAddr) string {
	return s.emailRootFromHex(addr.HexKey())
}

var (
	rxAddrKey         = regexp.MustCompile(`^[0-9a-f]{7,}$`)
	rxSHA1            = regexp.MustCompile(`^[0-9a-f]{40,40}$`)
	errInvalidBlobref = errors.New("Invalid sha1")
)

func (s *diskStorage) encFile(addrKey, encSHA1 string) (string, error) {
	if !rxAddrKey.MatchString(addrKey) || !rxSHA1.MatchString(encSHA1) {
		return "", errInvalidBlobref
	}
	return filepath.Join(s.emailRootFromHex(addrKey), encSHA1), nil
}

func (s *diskStorage) StatEncryptedEmail(addrKey, encSHA1 string) (size int, err error) {
	path, err := s.encFile(addrKey, encSHA1)
	if err != nil {
		return
	}
	fi, err := os.Stat(path)
	if err != nil {
		return
	}
	return int(fi.Size()), nil
}

func (s *diskStorage) EncryptedEmail(addrKey, encSHA1 string) ([]byte, error) {
	path, err := s.encFile(addrKey, encSHA1)
	if err != nil {
		return nil, err
	}
	return ioutil.ReadFile(path)
}

func (s *diskStorage) PutEmail(addr *webfist.EmailAddr, email *webfist.Email) error {
	emailRoot := s.emailRoot(addr)
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
	emailRoot := s.emailRoot(addr)
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
