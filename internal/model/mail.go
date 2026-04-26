package model

type MailRequest struct {
	SenderEmail string `json:"sender_email"`
	Subject     string `json:"subject"`
	Message     string `json:"message"`
}

type MailResponse struct {
	Message string `json:"message"`
}

type ErrorResponse struct {
	Error string `json:"error"`
}
