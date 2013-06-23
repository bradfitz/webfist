package main

import (
	"fmt"
	"net/http"
)

func init() {
	http.HandleFunc("/", front)
}

func front(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, `
<html>
<head>
	<meta charset="utf-8">
	<title>WebFist</title>
	<style>
		html, body {
			width: 100%;
			height: 100%;
		}
		body {
			font-family: monospace;
			font-size: 12px;
		}
		p {
			width: 580px;
		}
		.container {
			width: 600px;
			margin: 0 auto;
		}
	</style>
</head>
<body>
<div class="container">
<h1>WebFist!</h1>
<p>This is a WebFist server. See the source at <a
href="https://github.com/bradfitz/webfist">github.com/bradfitz/webfist</a>. You
can run your own WebFist fall-back server.</p>
<p>You can send email to <a href="mailto:fist@webfist.org">fist@webfist.org</a>
to set your WebFinger delegation. The email must be DKIM-signed. You will not receive
a response email. The contents of the email should be:</p>
<p>
<code>webfist = http://example.com/path/to/your-profile</code>
</p>
<p>Or, send yourself an email and paste the full email headers and body
here:</p>
<form method='POST' action='/add'>
<textarea name='email' rows=20 cols=80></textarea><br/>
<input type='submit' value='Submit' />
</form>
[<a href="/webfist/bump">Recent changes</a>]
<center>
<p><img width=561 height=798 src='http://upload.wikimedia.org/wikipedia/commons/1/17/Fist.png'></p>
</center>
</div>
</body>
</html>
`)
}
