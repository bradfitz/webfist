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
	"strings"
	"crypto/sha1"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"regexp"
	"sort"

	"github.com/bradfitz/webfist"
)

const recentCount = 1000

type diskStorage struct {
	root      string
	recentDir string
}

func NewDiskStorage(root string) (webfist.Storage, error) {
	ds := &diskStorage{
		root: root,
	}
	ds.recentDir = filepath.Join(root, "recent")
	if err := os.MkdirAll(ds.recentDir, 0700); err != nil {
		return nil, err
	}
	go ds.cleanRecent()
	return ds, nil
}

func (s *diskStorage) RecentMeta() (rm []*webfist.RecentMeta, err error) {
	f, err := os.Open(s.recentDir)
	if err != nil {
		return
	}
	defer f.Close()
	fis, err := f.Readdir(-1)
	if err != nil {
		return
	}
	sort.Sort(byModTime(fis))
	for _, fi := range fis {
		name := fi.Name()
		if !strings.HasSuffix(name, ".recent") {
			continue
		}
		name = strings.TrimSuffix(name, ".recent")
		s := strings.SplitN(name, "-", 2)
		if len(s) != 2 {
			continue
		}
		rm = append(rm, &webfist.RecentMeta{
			AddrHexKey: s[0],
			EncSHA1:    s[1],
			AddTime:    fi.ModTime(),
		})
	}
	return rm, nil
}

func (s *diskStorage) cleanRecent() error {
	f, err := os.Open(s.recentDir)
	if err != nil {
		return err
	}
	defer f.Close()
	fis, err := f.Readdir(-1)
	if err != nil {
		return err
	}
	if len(fis) <= recentCount {
		return nil
	}
	sort.Sort(byModTime(fis))
	toDelete := fis[:len(fis)-recentCount]
	for _, fi := range toDelete {
		path := filepath.Join(s.recentDir, filepath.Base(fi.Name()))
		os.Remove(path)
	}
	return nil
}

type byModTime []os.FileInfo

func (s byModTime) Len() int           { return len(s) }
func (s byModTime) Swap(i, j int)      { s[i], s[j] = s[j], s[i] }
func (s byModTime) Less(i, j int) bool { return s[i].ModTime().Before(s[j].ModTime()) }

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

func (s *diskStorage) touchRecent(addrKey, encSHA1 string) error {
	// Touch the recent file.
	err := ioutil.WriteFile(filepath.Join(s.recentDir, addrKey+"-"+encSHA1+".recent"), nil, 0600)
	if err == nil {
		go s.cleanRecent()
	}
	return err
}

func (s *diskStorage) PutEncryptedEmail(addrKey, encSHA1 string, data []byte) error {
	s1 := sha1.New()
	s1.Write(data)
	if fmt.Sprintf("%x", s1.Sum(nil)) != encSHA1 {
		return errInvalidBlobref
	}
	emailRoot := s.emailRootFromHex(addrKey)
	err := os.MkdirAll(emailRoot, 0755)
	if err != nil {
		return err
	}
	emailPath := filepath.Join(emailRoot, encSHA1)
	if err := ioutil.WriteFile(emailPath, data, 0644); err != nil {
		return err
	}
	return s.touchRecent(addrKey, encSHA1)
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

	addrKey := addr.HexKey()
	encSHA1 := fmt.Sprintf("%x", s1.Sum(nil))
	email.SetEncSHA1(encSHA1)

	emailPath := filepath.Join(emailRoot, encSHA1)
	if err := ioutil.WriteFile(emailPath, enc, 0644); err != nil {
		return err
	}
	return s.touchRecent(addrKey, encSHA1)
}

func (s *diskStorage) Emails(addr *webfist.EmailAddr) ([]*webfist.Email, error) {
	emailRoot := s.emailRoot(addr)
	file, err := os.Open(emailRoot)
	if os.IsNotExist(err) {
		return nil, nil
	}
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
