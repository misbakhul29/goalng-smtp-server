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
	subject := sanitize(req.Subject)
	body := buildBody(req.Sender, req.Email)

	headers := buildHeaders(s.cfg.SMTPEmail, req.SendTo, req.Sender, subject)
	raw := headers + "\r\n" + body

	addr := fmt.Sprintf("%s:%d", s.cfg.SMTPHost, s.cfg.SMTPPort)
	auth := smtp.PlainAuth("", s.cfg.SMTPEmail, s.cfg.SMTPPassword, s.cfg.SMTPHost)

	done := make(chan error, 1)
	go func() {
		done <- smtp.SendMail(addr, auth, s.cfg.SMTPEmail, []string{req.SendTo}, []byte(raw))
	}()

	select {
	case err := <-done:
		if err != nil {
			return fmt.Errorf("smtp send failed: %w", err)
		}
		log.Printf("email sent successfully | from=%s to=%s subject=%q", req.Sender, req.SendTo, req.Subject)
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

func buildBody(sender, email string) string {
	return fmt.Sprintf(
		"From: %s\r\n\r\nMessage:\r\n%s",
		sender,
		email,
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
