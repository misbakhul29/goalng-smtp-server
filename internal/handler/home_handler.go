package handler

import (
	"net/http"
)

const homeHTML = `
<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>golangsmtp - SMTP API Service</title>
    <style>
        body { font-family: -apple-system, BlinkMacSystemFont, "Segoe UI", Roboto, Helvetica, Arial, sans-serif; line-height: 1.6; color: #333; max-width: 800px; margin: 0 auto; padding: 2rem; text-align: center; margin-top: 10vh; }
        h1 { color: #2c3e50; font-size: 2.5rem; margin-bottom: 0.5rem; }
        p { color: #666; font-size: 1.2rem; margin-bottom: 2rem; }
        .button { display: inline-block; background: #3498db; color: white; padding: 0.8rem 1.5rem; text-decoration: none; border-radius: 4px; font-weight: bold; transition: background 0.2s; }
        .button:hover { background: #2980b9; }
    </style>
</head>
<body>
    <h1>golangsmtp</h1>
    <p>A fast, lightweight, production-ready email sender API service.</p>
    <a href="/docs/api" class="button">View API Documentation</a>
</body>
</html>
`

func Home() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/" {
			http.NotFound(w, r)
			return
		}
		
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(homeHTML))
	}
}
