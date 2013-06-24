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

package webfist

import (
	"bytes"
	"io"
	"testing"
)

func TestEmailKey(t *testing.T) {
	key := NewEmailAddr("brad@danga.com").HexKey()
	if want := "f888951d2ddcad78ffebce4a2c3158ecd1a60db0811a924ae7f41204828937c3"; key != want {
		t.Errorf("key = %q; want %q", key, want)
	}
}

func TestEncrypt(t *testing.T) {
	const msg = "From: foo\r\nTo: bar\r\n"
	email := NewEmailAddr("brad@danga.com")

	var encBuf bytes.Buffer
	enc := email.Encrypter(&encBuf)
	enc.Write([]byte(msg))

	if encBuf.String() != "\xcd\xe2\x136n\xbe\xd4c\xf0\xefy4\xc5T\xe6\xda5o\x865" {
		t.Errorf("Encrypted value doesn't match what's expected.")
	}

	var decBuf bytes.Buffer
	io.Copy(&decBuf, email.Decrypter(&encBuf))
	if decBuf.String() != msg {
		t.Errorf("Decrypted = %q; want %q", decBuf.String(), msg)
	}
}
