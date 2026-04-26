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
        body { font-family: -apple-system, BlinkMacSystemFont, "Segoe UI", Roboto, Helvetica, Arial, sans-serif; line-height: 1.6; color: #333; margin: 0; padding: 0; background: #f8f9fa; }
        .hero { background: linear-gradient(135deg, #2c3e50 0%, #3498db 100%); color: white; text-align: center; padding: 8rem 2rem; display: flex; flex-direction: column; align-items: center; justify-content: center; min-height: 50vh; }
        .hero h1 { font-size: 3.5rem; margin-bottom: 1rem; font-weight: 800; letter-spacing: -1px; }
        .hero p { font-size: 1.25rem; max-width: 600px; margin: 0 auto 2.5rem auto; opacity: 0.9; }
        .button { display: inline-block; background: white; color: #2c3e50; padding: 1rem 2.5rem; text-decoration: none; border-radius: 50px; font-weight: bold; font-size: 1.1rem; transition: transform 0.2s, box-shadow 0.2s; box-shadow: 0 4px 6px rgba(0,0,0,0.1); }
        .button:hover { transform: translateY(-2px); box-shadow: 0 6px 12px rgba(0,0,0,0.15); }
        .features { display: flex; justify-content: center; gap: 2rem; padding: 4rem 2rem; max-width: 1000px; margin: 0 auto; flex-wrap: wrap; }
        .feature { background: white; padding: 2rem; border-radius: 12px; box-shadow: 0 4px 6px rgba(0,0,0,0.05); flex: 1; min-width: 250px; text-align: center; }
        .feature h3 { color: #2c3e50; margin-top: 0; font-size: 1.25rem; }
        .feature p { color: #666; font-size: 0.95rem; margin-bottom: 0; }
        footer { text-align: center; padding: 2rem; color: #888; font-size: 0.9rem; }
    </style>
</head>
<body>
    <div class="hero">
        <h1>golangsmtp</h1>
        <p>A fast, lightweight, production-ready email sender API service built with Go. Designed for portfolio contact forms and simple notification systems.</p>
        <a href="/docs/api" class="button">View API Documentation</a>
    </div>

    <div class="features">
        <div class="feature">
            <h3>⚡ Fast & Lightweight</h3>
            <p>Built with Go's standard library. No heavy frameworks, minimal memory footprint, and blazing fast execution.</p>
        </div>
        <div class="feature">
            <h3>🔒 Secure by Default</h3>
            <p>API Key authentication, header injection sanitization, and built-in rate limiting out of the box.</p>
        </div>
        <div class="feature">
            <h3>📖 Easy to Integrate</h3>
            <p>Simple JSON API that can be consumed by any frontend framework or backend service in minutes.</p>
        </div>
    </div>

    <footer>
        &copy; golangsmtp - SMTP API Service
    </footer>
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
