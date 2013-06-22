package webfist

import (
	"testing"
)

func TestDKIMVerifyAvailable(t *testing.T) {
	path, err := findDKIMVerify()
	if err != nil {
		t.Errorf(dkimFailMessage)
	}
	t.Logf("dkimverify at %s", path)
}
