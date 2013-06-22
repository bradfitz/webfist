package webfist

import (
	"bytes"
	"io"
	"testing"
)

func TestEmailKey(t *testing.T) {
	key := EmailKeyString("brad@danga.com")
	if want := "f888951d2ddcad78ffebce4a2c3158ecd1a60db0811a924ae7f41204828937c3"; key != want {
		t.Errorf("key = %q; want %q", key, want)
	}
}

func TestEncrypt(t *testing.T) {
	const msg = "From: foo\r\nTo: bar\r\n"
	const email = "brad@danga.com"

	var encBuf bytes.Buffer
	enc := NewEncrypter(email, &encBuf)
	enc.Write([]byte(msg))

	if encBuf.String() != "\xcd\xe2\x136n\xbe\xd4c\xf0\xefy4\xc5T\xe6\xda5o\x865" {
		t.Errorf("Encrypted value doesn't match what's expected.")
	}

	var decBuf bytes.Buffer
	io.Copy(&decBuf, NewDecrypter(email, &encBuf))
	if decBuf.String() != msg {
		t.Errorf("Decrypted = %q; want %q", decBuf.String(), msg)
	}
}
