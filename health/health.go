package health

import (
	"net/http"
	"strings"
)

func Health(w http.ResponseWriter, r *http.Request) {
	userAgent := r.Header.Get("User-Agent")
	if strings.Contains(strings.ToLower(userAgent), "curl") { // Check, is it curl?
		w.WriteHeader(http.StatusOK) // Send status 200
		return
	}

	// if request from browser return html page with text OK
	w.Header().Set("Content-Type", "text/html")
	html := `<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <title>Health Check</title>
</head>
<body>
    <h1>I'm OK!</h1>
</body>
</html>`
	w.Write([]byte(html))
}
