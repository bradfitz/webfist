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
			margin: 1em 0;
		}
		.container {
			width: 600px;
			margin: 0 auto;
		}
		code {
			background-color: #eee;
			border: 2px solid #ccc;
			padding: 8px;
			border-radius: 4px;
		}
	</style>
</head>
<body>
<div class="container">
<h1>WebFist!</h1>
<p>This is a WebFist server. See the source at <a
href="https://github.com/bradfitz/webfist">github.com/bradfitz/webfist</a>. You
can run your own WebFist server and join the fist-bump network. Use WebFist
servers to do WebFinger look-ups when the canonical reply (from the domain that
owns an email address) is invalid or missing.</p>
<p>Send email to <a href="mailto:fist@webfist.org">fist@webfist.org</a>
to set your WebFinger delegation. The email must be DKIM-signed. You will not
receive a response email. The contents of the email should be like this. The URL
you point to should be a <a href="http://tools.ietf.org/html/draft-ietf-appsawg-webfinger">JRD document</a>
defined by the WebFinger spec.</p>
<br>
<code>webfist = http://example.com/path/to/your-profile</code>
<br>
<br>
<br>
<p>
Lookup your email address using WebFist:
<form method="GET" action="/.well-known/webfinger">
<input type="text" name="resource" placeholder="acct:foo@example.com" size="40">
<input type="submit" value="Lookup">
</form>
</p>
<br>
<p>Or, send yourself an email and paste the full email headers and body here:
<form method='POST' action='/add'>
<textarea name="email" rows=20 cols=80></textarea><br>
<input type="submit" value="Submit">
</form>
[<a href="/webfist/bump">Recent changes</a>]
<center>
<p><img width="561" height="798" src="http://upload.wikimedia.org/wikipedia/commons/1/17/Fist.png"></p>
</center>
</div>
</body>
</html>
`)
}
