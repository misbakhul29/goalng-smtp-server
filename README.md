# golangsmtp

A lightweight, production-ready email sender service built with Go — designed for portfolio website contact forms.

## Features

- `POST /api/send-email` — sends email via Gmail SMTP
- Input validation (email format, non-empty fields)
- Header injection sanitization
- Rate limiting (5 requests/min per IP)
- 15-second SMTP timeout
- Environment-based configuration
- PM2 process management support

---

## Project Structure

```
golangsmtp/
├── cmd/
│   └── server/
│       └── main.go
├── internal/
│   ├── config/
│   │   └── config.go
│   ├── handler/
│   │   └── mail_handler.go
│   ├── middleware/
│   │   └── rate_limiter.go
│   ├── model/
│   │   └── mail.go
│   └── service/
│       └── mail_service.go
├── ecosystem.config.js
├── .env.example
├── .gitignore
├── go.mod
└── go.sum
```

---

## Requirements

- Go 1.21+
- Gmail account with [App Password](https://myaccount.google.com/apppasswords) enabled (requires 2FA)
- PM2 (optional, for production deployment)

---

## Setup

**1. Clone the repository**

```bash
git clone https://github.com/your-username/golangsmtp.git
cd golangsmtp
```

**2. Install dependencies**

```bash
go mod tidy
```

**3. Configure environment**

```bash
cp .env.example .env
```

Edit `.env`:

```env
SMTP_HOST=smtp.gmail.com
SMTP_PORT=587
SMTP_EMAIL=your_email@gmail.com
SMTP_PASSWORD=your_app_password_here
SERVER_PORT=8000
```

> **Important:** Use a Gmail App Password, not your account password. Generate one at [myaccount.google.com/apppasswords](https://myaccount.google.com/apppasswords).

---

## Running

**Development**

```bash
go run ./cmd/server/main.go
```

**Production (binary)**

```bash
go build -o golangsmtp.exe ./cmd/server/main.go
./golangsmtp.exe
```

**Production (PM2)**

```bash
go build -o golangsmtp.exe ./cmd/server/main.go
pm2 start ecosystem.config.js
pm2 save
pm2 startup
```

---

## API Reference

### `POST /api/send-email`

**Request**

```json
{
  "sender_email": "user@example.com",
  "subject": "Hello",
  "message": "I want to work with you"
}
```

**Responses**

| Status | Body |
|--------|------|
| `200 OK` | `{ "message": "Email sent successfully" }` |
| `400 Bad Request` | `{ "error": "sender_email is not a valid email address" }` |
| `429 Too Many Requests` | `{ "error": "too many requests, please slow down" }` |
| `500 Internal Server Error` | `{ "error": "failed to send email, please try again later" }` |

### `GET /health`

Returns `{ "status": "ok" }` — useful for uptime monitoring.

---

**Test with curl:**

```bash
curl -X POST http://localhost:8000/api/send-email \
  -H "Content-Type: application/json" \
  -d '{"sender_email":"test@example.com","subject":"Hello","message":"I want to work with you"}'
```

---

## PM2 Commands

```bash
pm2 logs golangsmtp      # tail logs
pm2 status               # process status
pm2 restart golangsmtp   # restart
pm2 stop golangsmtp      # stop
```

---

## Security Notes

- `.env` is gitignored — never commit credentials
- SMTP credentials are loaded from environment variables only
- Subject line is sanitized against header injection attacks
- Rate limiter prevents abuse (5 requests/min per IP)

---

## License

MIT
