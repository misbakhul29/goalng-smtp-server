package model

type MailRequest struct {
	Sender  string `json:"sender"`
	SendTo  string `json:"send_to"`
	Subject string `json:"subject"`
	Email   string `json:"email"`
}

type MailResponse struct {
	Message string `json:"message"`
}

type ErrorResponse struct {
	Error string `json:"error"`
}
