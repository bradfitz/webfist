package webfist

import (
	"testing"
)

func TestEmailKey(t *testing.T) {
	key := EmailKey("brad@danga.com")
	if want := "f888951d2ddcad78ffebce4a2c3158ecd1a60db0811a924ae7f41204828937c3"; key != want {
		t.Errorf("key = %q; want %q", key, want)
	}
}
