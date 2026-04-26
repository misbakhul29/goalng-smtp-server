package handler

import (
	"fmt"
	"net/http"
)

const apiDocsHTML = `
<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>API Documentation - golangsmtp</title>
    <link rel="icon" href="data:image/svg+xml,<svg xmlns='http://www.w3.org/2000/svg' viewBox='0 0 100 100'><text y='.9em' font-size='90'>📧</text></svg>">
    <style>
        body { font-family: -apple-system, BlinkMacSystemFont, "Segoe UI", Roboto, Helvetica, Arial, sans-serif; line-height: 1.6; color: #333; max-width: 800px; margin: 0 auto; padding: 2rem; }
        h1, h2, h3 { color: #2c3e50; }
        pre { background: #f8f9fa; padding: 1rem; border-radius: 4px; overflow-x: auto; border: 1px solid #e9ecef; }
        code { font-family: ui-monospace, SFMono-Regular, Consolas, "Liberation Mono", Menlo, monospace; font-size: 0.9em; background: #f8f9fa; padding: 0.2em 0.4em; border-radius: 3px; }
        .endpoint { background: #e8f4f8; padding: 1rem; border-radius: 4px; border-left: 4px solid #3498db; margin-bottom: 2rem; }
        .method { font-weight: bold; color: #fff; background: #2ecc71; padding: 0.3rem 0.6rem; border-radius: 4px; margin-right: 0.5rem; }
        table { width: 100%%; border-collapse: collapse; margin: 1rem 0; }
        th, td { text-align: left; padding: 0.75rem; border-bottom: 1px solid #dee2e6; }
        th { background: #f8f9fa; }
    </style>
</head>
<body>
    <h1>API Documentation</h1>
    <p>Welcome to the <code>golangsmtp</code> API documentation.</p>

    <div class="endpoint">
        <h2><span class="method">POST</span> /api/send-email</h2>
        <p>Sends an email to the configured portfolio owner's email address.</p>
    </div>

    <h3>Authentication</h3>
    <p>All requests must include an API key in the headers. To request an API key, please email <a href="mailto:misbakhul2904@gmail.com">misbakhul2904@gmail.com</a>.</p>
    <table>
        <tr><th>Header</th><th>Description</th></tr>
        <tr><td><code>X-API-Key</code></td><td>Your secret API key (required)</td></tr>
    </table>

    <h3>Request Format</h3>
    <p>The request body must be JSON and include the following fields:</p>
    <table>
        <tr><th>Field</th><th>Type</th><th>Description</th></tr>
        <tr><td><code>sender_email</code></td><td>string</td><td>Valid email address of the sender</td></tr>
        <tr><td><code>subject</code></td><td>string</td><td>Email subject (cannot be empty)</td></tr>
        <tr><td><code>message</code></td><td>string</td><td>Email body content (cannot be empty)</td></tr>
    </table>

    <h4>Example Request</h4>
    <pre><code>curl -X POST https://%s/api/send-email \
  -H "Content-Type: application/json" \
  -H "X-API-Key: your_api_key_here" \
  -d '{
    "sender_email": "test@example.com",
    "subject": "Hello",
    "message": "I want to work with you"
  }'</code></pre>

    <h3>Responses</h3>
    <table>
        <tr><th>Status</th><th>Description</th><th>Body</th></tr>
        <tr><td><code>200 OK</code></td><td>Email sent successfully</td><td><code>{"message": "Email sent successfully"}</code></td></tr>
        <tr><td><code>400 Bad Request</code></td><td>Validation failed</td><td><code>{"error": "sender_email is not a valid email address"}</code></td></tr>
        <tr><td><code>401 Unauthorized</code></td><td>Missing API key</td><td><code>{"error": "missing X-API-Key header"}</code></td></tr>
        <tr><td><code>403 Forbidden</code></td><td>Invalid API key</td><td><code>{"error": "invalid API key"}</code></td></tr>
        <tr><td><code>429 Too Many Requests</code></td><td>Rate limit exceeded (5 req/min)</td><td><code>{"error": "too many requests, please slow down"}</code></td></tr>
        <tr><td><code>500 Internal Server Error</code></td><td>SMTP/Server error</td><td><code>{"error": "failed to send email, please try again later"}</code></td></tr>
    </table>

    <hr>
    
    <div class="endpoint">
        <h2><span style="background: #3498db" class="method">GET</span> /health</h2>
        <p>Basic health check endpoint to verify the service is running.</p>
    </div>
    
    <h4>Example Request</h4>
    <pre><code>curl https://%s/health</code></pre>

    <h4>Example Response</h4>
    <pre><code>{"status":"ok"}</code></pre>

</body>
</html>
`

func APIDocs(domain string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(fmt.Sprintf(apiDocsHTML, domain, domain)))
	}
}
