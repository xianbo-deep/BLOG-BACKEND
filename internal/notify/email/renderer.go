package email

import (
	"bytes"
	"errors"
	"html/template"
)

type Renderer struct {
	tpls map[string]*template.Template
}

func NewRenderer() *Renderer {
	tpls := map[string]*template.Template{
		MailDeadlinkReport:   template.Must(template.ParseFiles(DeadLinkFile)),
		MailDiscussionNotify: template.Must(template.ParseFiles(DiscussionNotifyFile)),
		MailDiscussionDigest: template.Must(template.ParseFiles(DiscussionReportFile)),
		MailSubscribeNotify:  template.Must(template.ParseFiles(SubscribeNotifyFile)),
		MailSubscribe:        template.Must(template.ParseFiles(SubscribeFile)),
		MailUnSubscribe:      template.Must(template.ParseFiles(UnSubscribeFile)),
	}
	return &Renderer{tpls}
}

func (r *Renderer) Render(t string, data any) (string, error) {
	// 拿到template对象
	tpl, ok := r.tpls[t]
	if !ok {
		return "", errors.New("template not found")
	}
	var buf bytes.Buffer
	// 执行渲染
	if err := tpl.Execute(&buf, data); err != nil {
		return "", err
	}
	return buf.String(), nil
}
