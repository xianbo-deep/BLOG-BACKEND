package email

import (
	"Blog-Backend/consts"
	"os"
	"strconv"
)

func RegisterEmail() {
	host := os.Getenv("EMAIL_HOST")
	portStr := os.Getenv("EMAIL_PORT")
	port, _ := strconv.Atoi(portStr)
	user := os.Getenv("EMAIL_USER")
	smtp := os.Getenv("EMAIL_SMTP")
	from := os.Getenv("EMAIL_FROM")
	cfg := EmailConfig{
		Host: host,
		Port: port,
		User: user,
		Pass: smtp,
		From: from,
	}

	emailClient := NewEmailClient(cfg)
	renderer := NewRenderer()

	mailer := NewMailer(emailClient, renderer)

}
