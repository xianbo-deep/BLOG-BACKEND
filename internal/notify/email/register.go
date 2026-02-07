package email

import (
	"os"
	"strconv"
)

func RegisterEmail() *Mailer {
	host := os.Getenv(EnvEmailHost)
	portStr := os.Getenv(EnvEmailPort)
	port, _ := strconv.Atoi(portStr)
	user := os.Getenv(EnvEmailUser)
	smtp := os.Getenv(EnvEmailSMTP)
	from := os.Getenv(EnvEmailFrom)
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

	return mailer
}
