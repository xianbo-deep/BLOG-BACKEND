package email

import (
	"fmt"
	"strings"

	"gopkg.in/gomail.v2"
)

type EmailClient struct {
	cfg EmailConfig
}

func NewEmailClient(cfg EmailConfig) *EmailClient {
	return &EmailClient{cfg: cfg}
}

func (e *EmailClient) SendHTML(to []string, subject string, html string) (err error) {
	if len(to) == 0 {
		return fmt.Errorf("no email addresses to send")
	}
	if strings.TrimSpace(subject) == "" {
		return fmt.Errorf("no email subject")
	}
	if strings.TrimSpace(html) == "" {
		return fmt.Errorf("no email html")
	}

	m := gomail.NewMessage()

	// 设置信息
	m.SetHeader("From", e.cfg.From)
	m.SetHeader("To", e.cfg.From)
	m.SetHeader("Bcc", to...)
	m.SetHeader("Subject", subject)
	m.SetBody("text/html", html)

	// 设置发件人
	d := gomail.NewDialer(e.cfg.Host, e.cfg.Port, e.cfg.User, e.cfg.Pass)

	// 发送邮件
	if err := d.DialAndSend(m); err != nil {
		return fmt.Errorf("email send failed: %v", err)
	}
	return nil
}
