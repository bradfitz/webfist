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
