package main

import (
	"fmt"
	"net/http"
)

func init() {
	http.HandleFunc("/", front)
}

func front(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, `<html><body>
<h1>WebFist!</h1>
<p>This is a WebFist server. See <a href="https://github.com/bradfitz/webfist">github.com/bradfitz/webfist</a>.</p>
<p>You can send email to <a href="mailto:fist@webfist.org">fist@webfist.org</a> to set your WebFinger keys. The email must be DKIM-signed. You will not receive a response email.</p>
<p>Or, send yourself an email and paste the full email headers and body here:</p>
<form method='POST' action='/add'>
<textarea name='email' rows=20 cols=80></textarea><br/>
<input type='submit' value='Submit' />
</form>

<center>
<p><img width=561 height=798 src='http://upload.wikimedia.org/wikipedia/commons/1/17/Fist.png'></p>
</center>
</body></html>`)
}
