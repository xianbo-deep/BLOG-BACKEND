package email

import (
	"bytes"
	"embed"
	"errors"
	"html/template"
)

type Renderer struct {
	tpls map[string]*template.Template
}

//go:embed template/*.html
var tplFS embed.FS

func NewRenderer() *Renderer {
	tpls := map[string]*template.Template{
		MailDeadlinkReport:   template.Must(template.ParseFS(tplFS, DeadLinkFile)),
		MailDiscussionNotify: template.Must(template.ParseFS(tplFS, DiscussionNotifyFile)),
		MailDiscussionDigest: template.Must(template.ParseFS(tplFS, DiscussionReportFile)),
		MailSubscribeNotify:  template.Must(template.ParseFS(tplFS, SubscribeNotifyFile)),
		MailSubscribe:        template.Must(template.ParseFS(tplFS, SubscribeFile)),
		MailUnSubscribe:      template.Must(template.ParseFS(tplFS, UnSubscribeFile)),
		MailSubscribeVerify:  template.Must(template.ParseFS(tplFS, SubscribeVCFile)),
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

func (r *Renderer) RenderPlaintext(t string, data any) (string, error) {
	tpl, ok := r.tpls[t]
	if !ok {
		return "", errors.New("template not found")
	}
	var buf bytes.Buffer

	if err := tpl.Execute(&buf, data); err != nil {
		return "", err
	}
	return buf.String(), nil
}
