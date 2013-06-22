package main

import (
	"fmt"
	"net/http"
)

func init() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, `<html><body>
<center>
<h1>WebFist!</h1>
<p><img width=561 height=798 src='http://upload.wikimedia.org/wikipedia/commons/1/17/Fist.png'></p>
</center>
</body></html>`)
	})
}
