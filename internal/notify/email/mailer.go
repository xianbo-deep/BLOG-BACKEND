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
	case MailSubscribe:
		return SubscribeSubject
	case MailUnSubscribe:
		return UnSubscribeSubject
	case MailSubscribeVerify:
		return SubscribeVCSubject
	default:
		return ""
	}
}

func (m *Mailer) SendTemplate(to []string, emailType string, data any, isHTML bool) error {
	// 选择模板并填充
	var content string
	var err error
	if isHTML {
		content, err = m.renderer.Render(emailType, data)
	} else {
		content, err = m.renderer.RenderPlaintext(emailType, data)
	}

	if err != nil {
		return err
	}
	// 选择邮件标题
	subject := m.selectSubject(emailType)
	// 发送邮件
	if isHTML {
		return m.client.SendHTML(to, subject, content)
	} else {
		return m.client.SendPlainText(to, subject, content)
	}

}
