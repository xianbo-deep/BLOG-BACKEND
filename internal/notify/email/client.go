package email

import (
	"Blog-Backend/consts"
	"fmt"
	"log"
	"strings"
	"time"

	"gopkg.in/gomail.v2"
)

type EmailClient struct {
	cfg EmailConfig
}

func NewEmailClient(cfg EmailConfig) *EmailClient {
	return &EmailClient{cfg: cfg}
}

func (e *EmailClient) SendHTML(to []string, subject string, content string) (err error) {
	if len(to) == 0 {
		return fmt.Errorf("no email addresses to send")
	}
	if strings.TrimSpace(subject) == "" {
		return fmt.Errorf("no email subject")
	}
	if strings.TrimSpace(content) == "" {
		return fmt.Errorf("no email html")
	}

	// 设置发件人
	d := gomail.NewDialer(e.cfg.Host, e.cfg.Port, e.cfg.User, e.cfg.Pass)

	// 循环发送
	for _, addr := range to {
		m := gomail.NewMessage()
		m.SetHeader("From", e.cfg.From)
		m.SetHeader("To", addr)
		m.SetHeader("Message-ID", fmt.Sprintf("<%d@%s>", time.Now().UnixNano(), consts.BlogDomain))
		m.SetHeader("Subject", subject)
		m.SetBody("text/html", content)

		if err := d.DialAndSend(m); err != nil {
			log.Printf("error sending email: %s", err)
		}
	}

	return nil
}

func (e *EmailClient) SendPlainText(to []string, subject string, content string) (err error) {
	if len(to) == 0 {
		return fmt.Errorf("no email addresses to send")
	}
	if strings.TrimSpace(subject) == "" {
		return fmt.Errorf("no email subject")
	}

	d := gomail.NewDialer(e.cfg.Host, e.cfg.Port, e.cfg.User, e.cfg.Pass)

	// 循环发送
	for _, addr := range to {
		m := gomail.NewMessage()
		m.SetHeader("From", e.cfg.From)
		m.SetHeader("To", addr)
		m.SetHeader("Message-ID", fmt.Sprintf("<%d@%s>", time.Now().UnixNano(), consts.BlogDomain))
		m.SetHeader("Subject", subject)
		m.SetBody("text/html", content)

		if err := d.DialAndSend(m); err != nil {
			log.Printf("error sending email: %s", err)
		}
	}
	return nil

}
