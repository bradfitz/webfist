package main

import (
	"bytes"
	"fmt"
	"net/http"
	"time"
)

func (s *server) ServeRecent(w http.ResponseWriter, r *http.Request) {
	rm, err := s.storage.RecentMeta()
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	var max time.Time
	var buf bytes.Buffer
	for _, m := range rm {
		if m.AddTime.After(max) {
			max = m.AddTime
		}
		fmt.Fprintf(&buf, "%s %s-%s\n", m.AddTime.Format(time.RFC3339), m.AddrHexKey, m.EncSHA1)
	}
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	http.ServeContent(w, r, "", max, bytes.NewReader(buf.Bytes()))
}
