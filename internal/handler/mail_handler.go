package handler

import (
	"encoding/json"
	"log"
	"net/http"
	"net/mail"
	"strings"

	"github.com/pc-06/golangsmtp/internal/model"
	"github.com/pc-06/golangsmtp/internal/service"
)

type MailHandler struct {
	svc service.MailService
}

func NewMailHandler(svc service.MailService) *MailHandler {
	return &MailHandler{svc: svc}
}

func (h *MailHandler) SendEmail(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		writeError(w, http.StatusMethodNotAllowed, "method not allowed")
		return
	}

	var req model.MailRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	if validationErr := validate(&req); validationErr != "" {
		writeError(w, http.StatusBadRequest, validationErr)
		return
	}

	if err := h.svc.Send(&req); err != nil {
		log.Printf("mail send error: %v", err)
		writeError(w, http.StatusInternalServerError, "failed to send email, please try again later")
		return
	}

	writeJSON(w, http.StatusOK, model.MailResponse{Message: "Email sent successfully"})
}

func validate(req *model.MailRequest) string {
	if _, err := mail.ParseAddress(req.Sender); err != nil {
		return "sender is not a valid email address"
	}
	if _, err := mail.ParseAddress(req.SendTo); err != nil {
		return "send_to is not a valid email address"
	}
	if strings.TrimSpace(req.Subject) == "" {
		return "subject must not be empty"
	}
	if strings.TrimSpace(req.Email) == "" {
		return "email must not be empty"
	}
	return ""
}

func writeJSON(w http.ResponseWriter, status int, payload any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(payload)
}

func writeError(w http.ResponseWriter, status int, msg string) {
	writeJSON(w, status, model.ErrorResponse{Error: msg})
}
