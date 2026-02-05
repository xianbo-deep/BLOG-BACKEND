package email

type Mailer struct {
	client   *EmailClient
	renderer *Renderer
}

func NewMailer(client *EmailClient, renderer *Renderer) *Mailer {
	return &Mailer{client: client, renderer: renderer}
}

func (m *Mailer) selectSubject(subject string) string {
	switch subject {
	case MailDeadlinkReport:
		return DeadLinkSubject
	case MailDiscussionNotify:
		return DiscussionNotifySubject
	case MailDiscussionDigest:
		return DiscussionDigestSubject
	case MailSubscribeNotify:
		return SubscribeNotifySubject
	default:
		return ""
	}
}

func (m *Mailer) SendTemplate(to []string, emailType string, data any) error {
	// 选择模板并填充
	html, err := m.renderer.Render(emailType, data)
	if err != nil {
		return err
	}
	// 选择邮件标题
	subject := m.selectSubject(emailType)
	// 发送邮件
	return m.client.SendHTML(to, subject, html)
}
