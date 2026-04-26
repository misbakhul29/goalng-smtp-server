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
    
    <!-- Primary Meta Tags -->
    <title>Golang SMTP API Service - Fast & Secure Email Sender</title>
    <meta name="title" content="Golang SMTP API Service - Fast & Secure Email Sender">
    <meta name="description" content="A fast, lightweight, and production-ready Golang SMTP API service for developers. Secure by default, easy to integrate for portfolio contact forms and notifications.">
    <meta name="keywords" content="Golang SMTP API Service, Email Sender API, Go email sender, SMTP API, developer tools, golangsmtp">
    <meta name="robots" content="index, follow">
    <meta name="language" content="English">
    <meta name="author" content="Misbakhul Munir">
    <link rel="canonical" href="https://smtp.misbakhul.my.id/">
    <link rel="icon" href="data:image/svg+xml,<svg xmlns='http://www.w3.org/2000/svg' viewBox='0 0 100 100'><text y='.9em' font-size='90'>📧</text></svg>">

    <!-- Open Graph / Facebook -->
    <meta property="og:type" content="website">
    <meta property="og:url" content="https://smtp.misbakhul.my.id/">
    <meta property="og:title" content="Golang SMTP API Service - Fast & Secure Email Sender">
    <meta property="og:description" content="A fast, lightweight, and production-ready Golang SMTP API service for developers. Secure by default, easy to integrate for portfolio contact forms and notifications.">

    <!-- Twitter -->
    <meta property="twitter:card" content="summary_large_image">
    <meta property="twitter:url" content="https://smtp.misbakhul.my.id/">
    <meta property="twitter:title" content="Golang SMTP API Service - Fast & Secure Email Sender">
    <meta property="twitter:description" content="A fast, lightweight, and production-ready Golang SMTP API service for developers. Secure by default, easy to integrate for portfolio contact forms and notifications.">

    <!-- JSON-LD Structured Data -->
    <script type="application/ld+json">
    {
      "@context": "https://schema.org",
      "@type": "SoftwareApplication",
      "name": "golangsmtp",
      "operatingSystem": "Any",
      "applicationCategory": "DeveloperApplication",
      "description": "A fast, lightweight, and production-ready Golang SMTP API service for developers. Secure by default, easy to integrate for portfolio contact forms and notifications.",
      "url": "https://smtp.misbakhul.my.id/",
      "offers": {
        "@type": "Offer",
        "price": "0.00",
        "priceCurrency": "USD"
      }
    }
    </script>

    <style>
        body { font-family: -apple-system, BlinkMacSystemFont, "Segoe UI", Roboto, Helvetica, Arial, sans-serif; line-height: 1.6; color: #333; margin: 0; padding: 0; background: #f8f9fa; }
        .hero { background: linear-gradient(135deg, #2c3e50 0%, #3498db 100%); color: white; text-align: center; padding: 8rem 2rem; display: flex; flex-direction: column; align-items: center; justify-content: center; min-height: 50vh; }
        .hero h1 { font-size: 3.5rem; margin-bottom: 1rem; font-weight: 800; letter-spacing: -1px; }
        .hero p { font-size: 1.25rem; max-width: 650px; margin: 0 auto 2.5rem auto; opacity: 0.9; }
        .button { display: inline-block; background: white; color: #2c3e50; padding: 1rem 2.5rem; text-decoration: none; border-radius: 50px; font-weight: bold; font-size: 1.1rem; transition: transform 0.2s, box-shadow 0.2s; box-shadow: 0 4px 6px rgba(0,0,0,0.1); }
        .button:hover { transform: translateY(-2px); box-shadow: 0 6px 12px rgba(0,0,0,0.15); }
        .features { display: flex; justify-content: center; gap: 2rem; padding: 4rem 2rem; max-width: 1000px; margin: 0 auto; flex-wrap: wrap; }
        .feature { background: white; padding: 2rem; border-radius: 12px; box-shadow: 0 4px 6px rgba(0,0,0,0.05); flex: 1; min-width: 250px; text-align: center; }
        .feature h2 { color: #2c3e50; margin-top: 0; font-size: 1.25rem; }
        .feature p { color: #666; font-size: 0.95rem; margin-bottom: 0; }
        footer { text-align: center; padding: 2rem; color: #888; font-size: 0.9rem; }
    </style>
</head>
<body>
    <main>
        <section class="hero">
            <h1>Golang SMTP API Service</h1>
            <p>A fast, lightweight, production-ready email sender API service built with Go. Designed to seamlessly integrate with portfolio contact forms and simple notification systems.</p>
            <a href="/docs/api" class="button" title="View Golang SMTP API Documentation">View API Documentation</a>
        </section>

        <section class="features" aria-label="Key Features">
            <div class="feature">
                <h2>⚡ Fast & Lightweight</h2>
                <p>Built with Go's standard library. No heavy frameworks, minimal memory footprint, and blazing fast execution for this SMTP API service.</p>
            </div>
            <div class="feature">
                <h2>🔒 Secure by Default</h2>
                <p>API Key authentication, header injection sanitization, and built-in rate limiting out of the box to protect your email sender.</p>
            </div>
            <div class="feature">
                <h2>📖 Easy to Integrate</h2>
                <p>Simple JSON API that can be consumed by any frontend framework or backend service in minutes.</p>
            </div>
        </section>
    </main>

    <footer>
        &copy; 2026 Golang SMTP API Service - Fast & Secure Email Sender
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
