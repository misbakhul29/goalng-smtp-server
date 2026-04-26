package service

import (
	"fmt"
	"log"
	"net/smtp"
	"strings"
	"time"

	"github.com/pc-06/golangsmtp/internal/config"
	"github.com/pc-06/golangsmtp/internal/model"
)

type MailService interface {
	Send(req *model.MailRequest) error
}

type mailService struct {
	cfg *config.Config
}

func NewMailService(cfg *config.Config) MailService {
	return &mailService{cfg: cfg}
}

func (s *mailService) Send(req *model.MailRequest) error {
	subject := fmt.Sprintf("[Portfolio Contact] %s", sanitize(req.Subject))
	body := buildBody(req.SenderEmail, req.Message)

	headers := buildHeaders(s.cfg.SMTPEmail, s.cfg.SMTPEmail, req.SenderEmail, subject)
	raw := headers + "\r\n" + body

	addr := fmt.Sprintf("%s:%d", s.cfg.SMTPHost, s.cfg.SMTPPort)
	auth := smtp.PlainAuth("", s.cfg.SMTPEmail, s.cfg.SMTPPassword, s.cfg.SMTPHost)

	done := make(chan error, 1)
	go func() {
		done <- smtp.SendMail(addr, auth, s.cfg.SMTPEmail, []string{s.cfg.SMTPEmail}, []byte(raw))
	}()

	select {
	case err := <-done:
		if err != nil {
			return fmt.Errorf("smtp send failed: %w", err)
		}
		log.Printf("email sent successfully | from=%s subject=%q", req.SenderEmail, req.Subject)
		return nil
	case <-time.After(15 * time.Second):
		return fmt.Errorf("smtp send timed out after 15 seconds")
	}
}

func buildHeaders(from, to, replyTo, subject string) string {
	return strings.Join([]string{
		fmt.Sprintf("From: %s", from),
		fmt.Sprintf("To: %s", to),
		fmt.Sprintf("Reply-To: %s", replyTo),
		fmt.Sprintf("Subject: %s", subject),
		"MIME-Version: 1.0",
		"Content-Type: text/plain; charset=UTF-8",
	}, "\r\n")
}

func buildBody(senderEmail, message string) string {
	return fmt.Sprintf(
		"You received a new message from your portfolio:\r\n\r\nFrom: %s\r\n\r\nMessage:\r\n%s",
		senderEmail,
		message,
	)
}

func sanitize(s string) string {
	replacer := strings.NewReplacer(
		"\r", "",
		"\n", "",
		"\t", " ",
	)
	return strings.TrimSpace(replacer.Replace(s))
}
